package version

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/blang/semver"
	"gopkg.in/urfave/cli.v1"
)

// Version represents the value of the current semantic version
const Version = "3.6.0"

// PrintVersion handles the version command for sops
func PrintVersion(c *cli.Context) {
	out := fmt.Sprintf("%s %s", c.App.Name, c.App.Version)
	upstreamVersion, err := RetrieveLatestVersionFromUpstream()
	if err != nil {
		out += fmt.Sprintf("\n[warning] failed to retrieve latest version from upstream: %v\n", err)
	}
	outdated, err := AIsNewerThanB(upstreamVersion, Version)
	if err != nil {
		out += fmt.Sprintf("\n[warning] failed to compare current version with latest: %v\n", err)
	}
	if outdated {
		out += fmt.Sprintf("\n[info] sops %s is available, update with `go get -u go.mozilla.org/sops/v3/cmd/sops`\n", upstreamVersion)
	} else {
		out += " (latest)\n"
	}
	fmt.Fprintf(c.App.Writer, "%s", out)
}

// AIsNewerThanB takes 2 semver strings are returns true
// is the A is newer than B, false otherwise
func AIsNewerThanB(A, B string) (bool, error) {
	if strings.HasPrefix(B, "1.") {
		// sops 1.0 doesn't use the semver format, which will
		// fail the call to `make` below. Since we now we're
		// more recent than 1.X anyway, return true right away
		return true, nil
	}
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

// RetrieveLatestVersionFromUpstream gets the latest version from the source code at Github
func RetrieveLatestVersionFromUpstream() (string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/mozilla/sops/master/version/version.go")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, `const Version = "`) {
			comps := strings.Split(line, `"`)
			if len(comps) < 2 {
				return "", fmt.Errorf("Failed to parse version from upstream source")
			}
			// try to parse the version as semver
			_, err := semver.Make(comps[1])
			if err != nil {
				return "", fmt.Errorf("Retrieved version %q does not match semver format: %v", comps[1], err)
			}
			return comps[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("Version information not found in upstream file")
}
