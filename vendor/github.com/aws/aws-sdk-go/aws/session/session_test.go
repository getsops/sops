// +build go1.7

package session

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestNewDefaultSession(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	s := New(&aws.Config{Region: aws.String("region")})

	if e, a := "region", *s.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := http.DefaultClient, s.Config.HTTPClient; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if s.Config.Logger == nil {
		t.Errorf("expect not nil")
	}
	if e, a := aws.LogOff, *s.Config.LogLevel; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestNew_WithCustomCreds(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	customCreds := credentials.NewStaticCredentials("AKID", "SECRET", "TOKEN")
	s := New(&aws.Config{Credentials: customCreds})

	if e, a := customCreds, s.Config.Credentials; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

type mockLogger struct {
	*bytes.Buffer
}

func (w mockLogger) Log(args ...interface{}) {
	fmt.Fprintln(w, args...)
}

func TestNew_WithSessionLoadError(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_invalid_source_profile")

	logger := bytes.Buffer{}
	s := New(&aws.Config{Logger: &mockLogger{&logger}})

	if s == nil {
		t.Errorf("expect not nil")
	}

	svc := s3.New(s)
	_, err := svc.ListBuckets(&s3.ListBucketsInput{})

	if err == nil {
		t.Errorf("expect not nil")
	}
	if e, a := "ERROR: failed to create session with AWS_SDK_LOAD_CONFIG enabled", logger.String(); !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
	if e, a := (SharedConfigAssumeRoleError{
		RoleARN: "assume_role_invalid_source_profile_role_arn",
	}).Error(), err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
}

func TestSessionCopy(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_REGION", "orig_region")

	s := Session{
		Config:   defaults.Config(),
		Handlers: defaults.Handlers(),
	}

	newSess := s.Copy(&aws.Config{Region: aws.String("new_region")})

	if e, a := "orig_region", *s.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "new_region", *newSess.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestSessionClientConfig(t *testing.T) {
	s, err := NewSession(&aws.Config{
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("orig_region"),
		EndpointResolver: endpoints.ResolverFunc(
			func(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
				if e, a := "mock-service", service; e != a {
					t.Errorf("expect %q service, got %q", e, a)
				}
				if e, a := "other-region", region; e != a {
					t.Errorf("expect %q region, got %q", e, a)
				}
				return endpoints.ResolvedEndpoint{
					URL:           "https://" + service + "." + region + ".amazonaws.com",
					SigningRegion: region,
				}, nil
			},
		),
	})
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	cfg := s.ClientConfig("mock-service", &aws.Config{Region: aws.String("other-region")})

	if e, a := "https://mock-service.other-region.amazonaws.com", cfg.Endpoint; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "other-region", cfg.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "other-region", *cfg.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestNewSession_NoCredentials(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	s, err := NewSession()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if s.Config.Credentials == nil {
		t.Errorf("expect not nil")
	}
	if e, a := credentials.AnonymousCredentials, s.Config.Credentials; e == a {
		t.Errorf("expect different credentials, %v", e)
	}
}

func TestNewSessionWithOptions_OverrideProfile(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "other_profile")

	s, err := NewSessionWithOptions(Options{
		Profile: "full_profile",
	})
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if e, a := "full_profile_region", *s.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	creds, err := s.Config.Credentials.Get()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}
	if e, a := "full_profile_akid", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "full_profile_secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if e, a := "SharedConfigCredentials", creds.ProviderName; !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
}

func TestNewSessionWithOptions_OverrideSharedConfigEnable(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "0")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "full_profile")

	s, err := NewSessionWithOptions(Options{
		SharedConfigState: SharedConfigEnable,
	})
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if e, a := "full_profile_region", *s.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	creds, err := s.Config.Credentials.Get()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}
	if e, a := "full_profile_akid", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "full_profile_secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if e, a := "SharedConfigCredentials", creds.ProviderName; !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
}

func TestNewSessionWithOptions_OverrideSharedConfigDisable(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "full_profile")

	s, err := NewSessionWithOptions(Options{
		SharedConfigState: SharedConfigDisable,
	})
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if v := *s.Config.Region; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}

	creds, err := s.Config.Credentials.Get()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}
	if e, a := "full_profile_akid", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "full_profile_secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if e, a := "SharedConfigCredentials", creds.ProviderName; !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
}

func TestNewSessionWithOptions_OverrideSharedConfigFiles(t *testing.T) {
	restoreEnvFn := initSessionTestEnv()
	defer restoreEnvFn()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "config_file_load_order")

	s, err := NewSessionWithOptions(Options{
		SharedConfigFiles: []string{testConfigOtherFilename},
	})
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}

	if e, a := "shared_config_other_region", *s.Config.Region; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	creds, err := s.Config.Credentials.Get()
	if err != nil {
		t.Errorf("expect nil, %v", err)
	}
	if e, a := "shared_config_other_akid", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "shared_config_other_secret", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if e, a := "SharedConfigCredentials", creds.ProviderName; !strings.Contains(a, e) {
		t.Errorf("expect %v, to contain %v", e, a)
	}
}

func TestNewSessionWithOptions_Overrides(t *testing.T) {
	cases := []struct {
		InEnvs    map[string]string
		InProfile string
		OutRegion string
		OutCreds  credentials.Value
	}{
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_PROFILE":                 "other_profile",
			},
			InProfile: "full_profile",
			OutRegion: "full_profile_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "full_profile_akid",
				SecretAccessKey: "full_profile_secret",
				ProviderName:    "SharedConfigCredentials",
			},
		},
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_REGION":                  "env_region",
				"AWS_ACCESS_KEY":              "env_akid",
				"AWS_SECRET_ACCESS_KEY":       "env_secret",
				"AWS_PROFILE":                 "other_profile",
			},
			InProfile: "full_profile",
			OutRegion: "env_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "env_akid",
				SecretAccessKey: "env_secret",
				ProviderName:    "EnvConfigCredentials",
			},
		},
		{
			InEnvs: map[string]string{
				"AWS_SDK_LOAD_CONFIG":         "0",
				"AWS_SHARED_CREDENTIALS_FILE": testConfigFilename,
				"AWS_CONFIG_FILE":             testConfigOtherFilename,
				"AWS_PROFILE":                 "shared_profile",
			},
			InProfile: "config_file_load_order",
			OutRegion: "shared_config_region",
			OutCreds: credentials.Value{
				AccessKeyID:     "shared_config_akid",
				SecretAccessKey: "shared_config_secret",
				ProviderName:    "SharedConfigCredentials",
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			restoreEnvFn := initSessionTestEnv()
			defer restoreEnvFn()

			for k, v := range c.InEnvs {
				os.Setenv(k, v)
			}

			s, err := NewSessionWithOptions(Options{
				Profile:           c.InProfile,
				SharedConfigState: SharedConfigEnable,
			})
			if err != nil {
				t.Errorf("expect nil, %v", err)
			}

			creds, err := s.Config.Credentials.Get()
			if err != nil {
				t.Errorf("expect nil, %v", err)
			}
			if e, a := c.OutRegion, *s.Config.Region; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := c.OutCreds.AccessKeyID, creds.AccessKeyID; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := c.OutCreds.SecretAccessKey, creds.SecretAccessKey; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := c.OutCreds.SessionToken, creds.SessionToken; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := c.OutCreds.ProviderName, creds.ProviderName; !strings.Contains(a, e) {
				t.Errorf("expect %v, to contain %v", e, a)
			}
		})
	}
}
