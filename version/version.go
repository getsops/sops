package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/blang/semver"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/urfave/cli"
)

// Version represents the value of the current semantic version.
var Version = "3.10.1"

// PrintVersion prints the current version of sops. If the flag
// `--disable-version-check` is set or if the environment variable
// SOPS_DISABLE_VERSION_CHECK is set to a value that is considered
// true by https://pkg.go.dev/strconv#ParseBool, the function will
// not attempt to retrieve the latest version from the GitHub API.
//
// If the flag is not set, the function will attempt to retrieve
// the latest version from the GitHub API and compare it to the
// current version. If the latest version is newer, the function
// will print a message to stdout.
func PrintVersion(c *cli.Context) {
	out := strings.Builder{}

	out.WriteString(fmt.Sprintf("%s %s", c.App.Name, c.App.Version))

	if c.Bool("disable-version-check") && !c.Bool("check-for-updates") {
		out.WriteString("\n")
	} else {
		upstreamVersion, upstreamURL, err := RetrieveLatestReleaseVersion()
		if err != nil {
			out.WriteString(fmt.Sprintf("\n[warning] failed to retrieve latest version from upstream: %v\n", err))
		} else {
			outdated, err := AIsNewerThanB(upstreamVersion, Version)
			if err != nil {
				out.WriteString(fmt.Sprintf("\n[warning] failed to compare current version with latest: %v\n", err))
			} else {
				if outdated {
					out.WriteString(fmt.Sprintf("\n[info] a new version of sops (%s) is available, you can update by visiting: %s\n", upstreamVersion, upstreamURL))
				} else {
					out.WriteString(" (latest)\n")
				}
			}
		}
		if !c.Bool("check-for-updates") {
			out.WriteString(
				"\n[warning] Note that in a future version, sops will no longer check whether the current version is the latest when asking for sops' version." +
					" If you want to explicitly check for the latest version, add the `--check-for-updates` option to `sops --version`." +
					" This will hide this deprecation warning and will always check, even if the default behavior changes in the future.\n")
		}
	}
	fmt.Fprintf(c.App.Writer, "%s", out.String())
}

// AIsNewerThanB compares two semantic versions and returns true if A is newer
// than B. The function will return an error if either version is not a valid
// semantic version.
func AIsNewerThanB(A, B string) (bool, error) {
	if strings.HasPrefix(B, "1.") {
		// sops 1.0 doesn't use the semver format, which will
		// fail the call to `make` below. Since we now we're
		// more recent than 1.X anyway, return true right away
		return true, nil
	}

	// Trim the leading "v" from the version strings, if present.
	A, B = strings.TrimPrefix(A, "v"), strings.TrimPrefix(B, "v")

	vA, err := semver.Make(A)
	if err != nil {
		return false, err
	}
	vB, err := semver.Make(B)
	if err != nil {
		return false, err
	}
	if vA.Compare(vB) > 0 {
		// vA is newer than vB
		return true, nil
	}
	return false, nil
}

// RetrieveLatestVersionFromUpstream retrieves the most recent release version
// from GitHub. The function returns the latest version as a string, or an
// error if the request fails or the response cannot be parsed.
//
// Deprecated: This function is deprecated in favor of
// RetrieveLatestReleaseVersion, which also provides the URL of the latest
// release.
func RetrieveLatestVersionFromUpstream() (string, error) {
	tag, _, err := RetrieveLatestReleaseVersion()
	return strings.TrimPrefix(tag, "v"), err
}

// RetrieveLatestReleaseVersion fetches the latest release version from GitHub.
// Returns the latest version as a string and the release URL, or an error if
// the request failed or the response could not be parsed.
//
// The function first attempts redirection-based retrieval (HTTP 301). It's
// preferred over GitHub API due to no rate limiting, but may break on
// redirect changes. If the first attempt fails, it falls back to the GitHub
// API.
//
// Unlike RetrieveLatestVersionFromUpstream, it returns the tag (e.g. "v3.7.3").
func RetrieveLatestReleaseVersion() (tag, url string, err error) {
	const repository = "getsops/sops"
	return newReleaseFetcher().LatestRelease(repository)
}

// newReleaseFetcher creates and returns a new instance of the releaseFetcher,
// preconfigured with the necessary endpoint information for redirection-based
// and API-based release retrieval.
func newReleaseFetcher() releaseFetcher {
	return releaseFetcher{
		endpoint:    "https://github.com",
		apiEndpoint: "https://api.github.com",
	}
}

// releaseFetcher is a helper struct used for fetching release information
// from GitHub. It encapsulates the necessary endpoints for redirection-based
// and API-based retrieval methods.
type releaseFetcher struct {
	endpoint    string
	apiEndpoint string
}

// LatestRelease retrieves the most recent release version for a given repository
// by first attempting to fetch it using redirection-based approach. If this
// attempt fails, it then falls back to the versioned GitHub API for retrieval.
//
// It returns the latest version as a string along with its corresponding URL, or
// an error in case both retrieval methods are unsuccessful.
//
// This function combines the advantages of both retrieval strategies: the resilience
// of the redirection-based approach and the reliability of the versioned API usage.
// However, it's worth noting that the API usage can be affected by GitHub's rate limiting.
func (f releaseFetcher) LatestRelease(repository string) (tag, url string, err error) {
	if tag, url, err = f.LatestReleaseUsingRedirect(repository); err == nil {
		return
	}
	return f.LatestReleaseUsingAPI(repository)
}

// LatestReleaseUsingRedirect fetches the most recent version of a release
// from the GitHub API. It returns the latest version as a string, along with
// its corresponding URL, or an error in case of a failed request or if the
// response couldn't be parsed.
//
// This method employs a customized HTTP client capable of following HTTP 301
// redirects, which might occur due to repository renaming. It's important to
// note that it does not follow HTTP 302 redirects, the type GitHub employs
// for redirecting to the latest release.
//
// Compared to LatestReleaseUsingAPI, this approach circumvents potential GitHub
// API rate limiting. However, it's worth considering that changes in GitHub's
// redirect handling could potentially disrupt its functionality.
func (f releaseFetcher) LatestReleaseUsingRedirect(repository string) (tag, url string, err error) {
	client := cleanhttp.DefaultClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Follow HTTP 301 redirects, which may be present due to the
		// repository being renamed. But do not follow HTTP 302 redirects,
		// which is what GitHub uses to redirect to the latest release.
		if req.Response.StatusCode == 302 {
			return http.ErrUseLastResponse
		}
		return nil
	}

	resp, err := client.Head(fmt.Sprintf("%s/%s/releases/latest", f.endpoint, repository))
	if err != nil {
		return "", "", err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 300 || resp.StatusCode > 399 {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", "", fmt.Errorf("missing Location header")
	}

	tagMarker := "releases/tag/"
	if tagIndex := strings.Index(location, tagMarker); tagIndex != -1 {
		return location[tagIndex+len(tagMarker):], location, nil
	}
	return "", "", fmt.Errorf("unexpected Location header: %s", location)
}

// LatestReleaseUsingAPI retrieves the most recent release version from the
// GitHub API. It returns the latest version as a string, along with its
// corresponding URL, or an error in case of request failure or parsing issues
// with the response.
//
// This approach boasts higher reliability compared to
// LatestReleaseUsingRedirect as it leverages the versioned GitHub API.
// However, it can be affected by GitHub API rate limiting.
func (f releaseFetcher) LatestReleaseUsingAPI(repository string) (tag, url string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/releases/latest", f.apiEndpoint, repository), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := cleanhttp.DefaultClient().Do(req)
	if err != nil {
		return "", "", fmt.Errorf("GitHub API request failed: %v", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API request failed with status code: %d", res.StatusCode)
	}

	type release struct {
		URL string `json:"html_url"`
		Tag string `json:"tag_name"`
	}
	var m release
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return "", "", err
	}
	return m.Tag, m.URL, nil
}
