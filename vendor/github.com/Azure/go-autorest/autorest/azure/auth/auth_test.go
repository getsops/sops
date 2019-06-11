// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package auth

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	expectedEnvironment = EnvironmentSettings{
		Values: map[string]string{
			SubscriptionID:      "sub-abc-123",
			TenantID:            "tenant-abc-123",
			ClientID:            "client-abc-123",
			ClientSecret:        "client-secret-123",
			CertificatePath:     "~/some/path/cert.pfx",
			CertificatePassword: "certificate-password",
			Username:            "user-name-abc",
			Password:            "user-password-123",
			Resource:            "my-resource",
		},
		Environment: azure.PublicCloud,
	}
	expectedFile = FileSettings{
		Values: map[string]string{
			ClientID:                "client-id-123",
			ClientSecret:            "client-secret-456",
			SubscriptionID:          "sub-id-789",
			TenantID:                "tenant-id-123",
			ActiveDirectoryEndpoint: "https://login.microsoftonline.com",
			ResourceManagerEndpoint: "https://management.azure.com/",
			GraphResourceID:         "https://graph.windows.net/",
			SQLManagementEndpoint:   "https://management.core.windows.net:8443/",
			GalleryEndpoint:         "https://gallery.azure.com/",
			ManagementEndpoint:      "https://management.core.windows.net/",
		},
	}
)

func setDefaultEnv() {
	os.Setenv(SubscriptionID, expectedEnvironment.Values[SubscriptionID])
	os.Setenv(TenantID, expectedEnvironment.Values[TenantID])
	os.Setenv(ClientID, expectedEnvironment.Values[ClientID])
	os.Setenv(ClientSecret, expectedEnvironment.Values[ClientSecret])
	os.Setenv(CertificatePath, expectedEnvironment.Values[CertificatePath])
	os.Setenv(CertificatePassword, expectedEnvironment.Values[CertificatePassword])
	os.Setenv(Username, expectedEnvironment.Values[Username])
	os.Setenv(Password, expectedEnvironment.Values[Password])
	os.Setenv(Resource, expectedEnvironment.Values[Resource])
}

func TestGetSettingsFromEnvironment(t *testing.T) {
	setDefaultEnv()
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		t.Logf("failed to get settings: %v", err)
		t.Fail()
	}
	if !reflect.DeepEqual(expectedEnvironment, settings) {
		t.Logf("expected %v, got %v", expectedEnvironment, settings)
		t.Fail()
	}
	if settings.GetSubscriptionID() != expectedEnvironment.Values[SubscriptionID] {
		t.Log("settings.GetSubscriptionID() return value didn't match")
		t.Fail()
	}
}

func TestGetSettingsFromEnvironmentBadEnvironmentName(t *testing.T) {
	os.Setenv(EnvironmentName, "badenvironment")
	defer func() {
		// must undo this value else other tests will fail
		os.Setenv(EnvironmentName, "")
	}()
	_, err := GetSettingsFromEnvironment()
	if err == nil {
		t.Log("unexpected nil error")
		t.Fail()
	}
}

func TestEnvGetClientCertificate(t *testing.T) {
	setDefaultEnv()
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		t.Logf("failed to get settings: %v", err)
		t.Fail()
	}
	cfg, err := settings.GetClientCertificate()
	if err != nil {
		t.Logf("failed to get config for client cert: %v", err)
		t.Fail()
	}
	if cfg.CertificatePath != expectedEnvironment.Values[CertificatePath] {
		t.Log("bad certificate path")
		t.Fail()
	}
	if cfg.CertificatePassword != expectedEnvironment.Values[CertificatePassword] {
		t.Log("bad certificate password")
		t.Fail()
	}
	// should fail as the certificate doesn't exist
	_, err = cfg.Authorizer()
	if err == nil {
		t.Log("unexpected nil error")
		t.Fail()
	}
}

func TestEnvGetUsernamePassword(t *testing.T) {
	setDefaultEnv()
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		t.Logf("failed to get settings: %v", err)
		t.Fail()
	}
	cfg, err := settings.GetUsernamePassword()
	if err != nil {
		t.Logf("failed to get config for username/password: %v", err)
		t.Fail()
	}
	_, err = cfg.Authorizer()
	if err != nil {
		t.Logf("failed to get authorizer for username/password: %v", err)
		t.Fail()
	}
}

func TestEnvGetMSI(t *testing.T) {
	setDefaultEnv()
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		t.Logf("failed to get settings: %v", err)
		t.Fail()
	}
	cfg := settings.GetMSI()
	_, err = cfg.Authorizer()
	if err != nil {
		t.Logf("failed to get authorizer for MSI: %v", err)
		t.Fail()
	}
}

func TestEnvGetDeviceFlow(t *testing.T) {
	setDefaultEnv()
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		t.Logf("failed to get settings: %v", err)
		t.Fail()
	}
	cfg := settings.GetDeviceFlow()
	// TODO mock device flow?
	if cfg.ClientID != expectedEnvironment.Values[ClientID] {
		t.Log("bad client ID")
		t.Fail()
	}
	if cfg.TenantID != expectedEnvironment.Values[TenantID] {
		t.Log("bad tenant ID")
		t.Fail()
	}
}

func TestGetSettingsFromFile(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", "./testdata/credsutf16le.json")
	settings, err := GetSettingsFromFile()
	if err != nil {
		t.Logf("failed to load config file: %v", err)
		t.Fail()
	}
	if !reflect.DeepEqual(expectedFile, settings) {
		t.Logf("expected %v, got %v", expectedFile, settings)
		t.Fail()
	}
	if settings.GetSubscriptionID() != expectedFile.Values[SubscriptionID] {
		t.Log("settings.GetSubscriptionID() return value didn't match")
		t.Fail()
	}
}

func TestNewAuthorizerFromFile(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", "./testdata/credsutf16le.json")
	authorizer, err := NewAuthorizerFromFile("https://management.azure.com")
	if err != nil || authorizer == nil {
		t.Logf("NewAuthorizerFromFile failed, got error %v", err)
		t.Fail()
	}
}

func TestNewAuthorizerFromFileWithResource(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", "./testdata/credsutf16le.json")
	authorizer, err := NewAuthorizerFromFileWithResource("https://my.vault.azure.net")
	if err != nil || authorizer == nil {
		t.Logf("NewAuthorizerFromFileWithResource failed, got error %v", err)
		t.Fail()
	}
}

func TestNewAuthorizerFromEnvironment(t *testing.T) {
	setDefaultEnv()
	authorizer, err := NewAuthorizerFromEnvironment()

	if err != nil || authorizer == nil {
		t.Logf("NewAuthorizerFromEnvironment failed, got error %v", err)
		t.Fail()
	}
}

func TestNewAuthorizerFromEnvironmentWithResource(t *testing.T) {
	setDefaultEnv()
	authorizer, err := NewAuthorizerFromEnvironmentWithResource("https://my.vault.azure.net")

	if err != nil || authorizer == nil {
		t.Logf("NewAuthorizerFromEnvironmentWithResource failed, got error %v", err)
		t.Fail()
	}
}

func TestDecodeAndUnmarshal(t *testing.T) {
	tests := []string{
		"credsutf8.json",
		"credsutf16le.json",
		"credsutf16be.json",
	}
	for _, test := range tests {
		os.Setenv("AZURE_AUTH_LOCATION", filepath.Join("./testdata/", test))
		settings, err := GetSettingsFromFile()
		if err != nil {
			t.Logf("error reading file '%s': %s", test, err)
			t.Fail()
		}
		if !reflect.DeepEqual(expectedFile, settings) {
			t.Logf("unmarshaled map expected %v, got %v", expectedFile, settings)
			t.Fail()
		}
	}
}

func TestFileClientCertificateAuthorizer(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", "./testdata/credsutf8.json")
	settings, err := GetSettingsFromFile()
	if err != nil {
		t.Logf("failed to load file settings: %v", err)
		t.Fail()
	}
	// add certificate settings
	settings.Values[CertificatePath] = "~/fake/path/cert.pfx"
	settings.Values[CertificatePassword] = "fake-password"
	_, err = settings.ClientCertificateAuthorizer("https://management.azure.com")
	if err == nil {
		t.Log("unexpected nil error")
		t.Fail()
	}
}
