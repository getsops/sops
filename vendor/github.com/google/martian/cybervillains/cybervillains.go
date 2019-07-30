// Copyright 2016 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cybervillains provides the publically published Selenium project CyberVillains
// certificate and key.  The CyberVillains cert and key allow for a man in the middle,
// and should only be used in testing scenarios. Client installation of the CyberVillains
// certificate is inherently and intentionally insecure.
package cybervillains

// Cert is the x509 CyberVillains public key published by the Selenium project.
const Cert string = `-----BEGIN CERTIFICATE-----
MIIClTCCAf6gAwIBAgIBATANBgkqhkiG9w0BAQUFADBZMRowGAYDVQQKDBFDeWJl
clZpbGxpYW5zLmNvbTEuMCwGA1UECwwlQ3liZXJWaWxsaWFucyBDZXJ0aWZpY2F0
aW9uIEF1dGhvcml0eTELMAkGA1UEBhMCVVMwHhcNMTEwMjEwMDMwMDEwWhcNMzEx
MDIzMDMwMDEwWjBZMRowGAYDVQQKDBFDeWJlclZpbGxpYW5zLmNvbTEuMCwGA1UE
CwwlQ3liZXJWaWxsaWFucyBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTELMAkGA1UE
BhMCVVMwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAIVQhWYazIfMvJUBP5qh
qRyh2tkrYI9wVZ9/Sj1l4tlWY4HOC6Dy5OYBRCmo2T9N8EXrAxXZKKUPzgmb3gIv
AQJ9DP6woiHyyztZJ5/cbhlp8EbHBIvGWK3T0Oph3kEPPS2FWKjiH/+pV6qY0Yt+
lkzcwxrjZIah/3VHQXUDm8X1AgMBAAGjbTBrMB0GA1UdDgQWBBQKvBeVNGu8hxtb
TP31Y4UttI/1bDASBgNVHRMBAf8ECDAGAQH/AgEAMAsGA1UdDwQEAwIBBjApBgNV
HSUEIjAgBggrBgEFBQcDAQYIKwYBBQUHAwkGCmCGSAGG+EUBCAEwDQYJKoZIhvcN
AQEFBQADgYEAD/6m8czx19uRPuaHVYhsEX5QGwJ4Y1NFswAByAuSBQB9KI9P2C7I
muf1aOoslUC4TxnC6g9H5/XmlK1zbZ+2YuABOb08CTXBC2x3ewJnm94DGPBRzj9o
0rXGEC+jsqsziBw+kg69xFn7PH09ZKUCue8POaaN/z5VoQMoM4ZNTP4=
-----END CERTIFICATE-----`

// Key is the CyberVillains private key published by the Selenium project.
const Key string = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCFUIVmGsyHzLyVAT+aoakcodrZK2CPcFWff0o9ZeLZVmOBzgug
8uTmAUQpqNk/TfBF6wMV2SilD84Jm94CLwECfQz+sKIh8ss7WSef3G4ZafBGxwSL
xlit09DqYd5BDz0thVio4h//qVeqmNGLfpZM3MMa42SGof91R0F1A5vF9QIDAQAB
AoGAEVuTkuDIYqIYp7n64xJLZ4v3Z7FLKEHzFApJy0a5y5yA5kTCpNkbTos5qcbv
SlvGfgQEadLVhPBS3lNqC5S9J7iUmmdpveXxV5ZaOsK3Zh+QCURfjLvqLH5Fzn1c
341YTCXpPdlbZElbARh3WKtW7R4c5GNNdf7zrWRqjYsXacECQQD4CVJ0l2AOTfLh
0uOXr1wwblIVscNv5WO9WLERtDWZP2EhDkRFMFsV8gTTvs01LiX0PRkuUjP+C6/e
g1DlBrqxAkEAiZhE4Ui7AHF6CYg+eamQKf4ECn4KgZ/y68Tan9YiULRXOx4HSpsM
3g+uPvwWnp9Pd/0gVSmQlJn3oNi5LQtIhQJBANF6ZgYL1lceY/NuvUJdGrnYYkDq
Ocml7P98CUePb/j2OxzExMm+Vh8JoCQIr5yrVeiZNUwWpsx2qFh/hPF4JnECQCej
/8wryPxStQcEAoPIjykZ7o4bS+mWbETynM3Jwm8f1bXJa+5ZhzZ+rAOnWtjuKtX1
zhfa9rVpOkdTyN2qT4UCQQD35VDm82aDi9mC8Zs1T/SrYKwRJuz25JPM8Yh9xiuK
7iI4qfwwaX99fo09cH0pfUdx+z7QyNba8bMfTWe8qPHm
-----END RSA PRIVATE KEY-----`
