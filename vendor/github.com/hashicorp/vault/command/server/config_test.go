package server

import (
	"reflect"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestLoadConfigFile(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfigFile("./test-fixtures/config.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
		},

		Storage: &Storage{
			Type:         "consul",
			RedirectAddr: "foo",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type:         "consul",
			RedirectAddr: "snafu",
			Config: map[string]string{
				"bar": "baz",
			},
			DisableClustering: true,
		},

		Telemetry: &Telemetry{
			StatsdAddr:              "bar",
			StatsiteAddr:            "foo",
			DisableHostname:         false,
			DogStatsDAddr:           "127.0.0.1:7254",
			DogStatsDTags:           []string{"tag_1:val_1", "tag_2:val_2"},
			PrometheusRetentionTime: prometheusDefaultRetentionTime,
		},

		DisableCache:             true,
		DisableCacheRaw:          true,
		DisableMlock:             true,
		DisableMlockRaw:          true,
		DisablePrintableCheckRaw: true,
		DisablePrintableCheck:    true,
		EnableUI:                 true,
		EnableUIRaw:              true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
		ClusterName:        "testcluster",

		PidFile: "./pidfile",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestLoadConfigFile_topLevel(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfigFile("./test-fixtures/config2.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
		},

		Storage: &Storage{
			Type:         "consul",
			RedirectAddr: "top_level_api_addr",
			ClusterAddr:  "top_level_cluster_addr",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type:         "consul",
			RedirectAddr: "top_level_api_addr",
			ClusterAddr:  "top_level_cluster_addr",
			Config: map[string]string{
				"bar": "baz",
			},
			DisableClustering: true,
		},

		Telemetry: &Telemetry{
			StatsdAddr:                 "bar",
			StatsiteAddr:               "foo",
			DisableHostname:            false,
			DogStatsDAddr:              "127.0.0.1:7254",
			DogStatsDTags:              []string{"tag_1:val_1", "tag_2:val_2"},
			PrometheusRetentionTime:    30 * time.Second,
			PrometheusRetentionTimeRaw: "30s",
		},

		DisableCache:    true,
		DisableCacheRaw: true,
		DisableMlock:    true,
		DisableMlockRaw: true,
		EnableUI:        true,
		EnableUIRaw:     true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
		ClusterName:        "testcluster",

		PidFile: "./pidfile",

		APIAddr:     "top_level_api_addr",
		ClusterAddr: "top_level_cluster_addr",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestLoadConfigFile_json(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfigFile("./test-fixtures/config.hcl.json", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
		},

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
			DisableClustering: true,
		},

		ClusterCipherSuites: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",

		Telemetry: &Telemetry{
			StatsiteAddr:                       "baz",
			StatsdAddr:                         "",
			DisableHostname:                    false,
			CirconusAPIToken:                   "",
			CirconusAPIApp:                     "",
			CirconusAPIURL:                     "",
			CirconusSubmissionInterval:         "",
			CirconusCheckSubmissionURL:         "",
			CirconusCheckID:                    "",
			CirconusCheckForceMetricActivation: "",
			CirconusCheckInstanceID:            "",
			CirconusCheckSearchTag:             "",
			CirconusCheckDisplayName:           "",
			CirconusCheckTags:                  "",
			CirconusBrokerID:                   "",
			CirconusBrokerSelectTag:            "",
			PrometheusRetentionTime:            prometheusDefaultRetentionTime,
		},

		MaxLeaseTTL:          10 * time.Hour,
		MaxLeaseTTLRaw:       "10h",
		DefaultLeaseTTL:      10 * time.Hour,
		DefaultLeaseTTLRaw:   "10h",
		ClusterName:          "testcluster",
		DisableCacheRaw:      interface{}(nil),
		DisableMlockRaw:      interface{}(nil),
		EnableUI:             true,
		EnableUIRaw:          true,
		PidFile:              "./pidfile",
		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,
		DisableSealWrap:      true,
		DisableSealWrapRaw:   true,
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestLoadConfigFile_json2(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfigFile("./test-fixtures/config2.hcl.json", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:444",
				},
			},
		},

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
			DisableClustering: true,
		},

		HAStorage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"bar": "baz",
			},
		},

		CacheSize: 45678,

		EnableUI: true,

		EnableRawEndpoint: true,

		DisableSealWrap: true,

		Telemetry: &Telemetry{
			StatsiteAddr:                       "foo",
			StatsdAddr:                         "bar",
			DisableHostname:                    true,
			CirconusAPIToken:                   "0",
			CirconusAPIApp:                     "vault",
			CirconusAPIURL:                     "http://api.circonus.com/v2",
			CirconusSubmissionInterval:         "10s",
			CirconusCheckSubmissionURL:         "https://someplace.com/metrics",
			CirconusCheckID:                    "0",
			CirconusCheckForceMetricActivation: "true",
			CirconusCheckInstanceID:            "node1:vault",
			CirconusCheckSearchTag:             "service:vault",
			CirconusCheckDisplayName:           "node1:vault",
			CirconusCheckTags:                  "cat1:tag1,cat2:tag2",
			CirconusBrokerID:                   "0",
			CirconusBrokerSelectTag:            "dc:sfo",
			PrometheusRetentionTime:            30 * time.Second,
			PrometheusRetentionTimeRaw:         "30s",
		},
	}
	if !reflect.DeepEqual(config, expected) {
	}
}

func TestLoadConfigDir(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	config, err := LoadConfigDir("./test-fixtures/config-dir", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		DisableCache: true,
		DisableMlock: true,

		DisableClustering:    false,
		DisableClusteringRaw: false,

		APIAddr:     "https://vault.local",
		ClusterAddr: "https://127.0.0.1:444",

		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
		},

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
			RedirectAddr:      "https://vault.local",
			ClusterAddr:       "https://127.0.0.1:444",
			DisableClustering: false,
		},

		EnableUI: true,

		EnableRawEndpoint: true,

		Telemetry: &Telemetry{
			StatsiteAddr:            "qux",
			StatsdAddr:              "baz",
			DisableHostname:         true,
			PrometheusRetentionTime: prometheusDefaultRetentionTime,
		},

		MaxLeaseTTL:     10 * time.Hour,
		DefaultLeaseTTL: 10 * time.Hour,
		ClusterName:     "testcluster",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestParseListeners(t *testing.T) {
	obj, _ := hcl.Parse(strings.TrimSpace(`
listener "tcp" {
	address = "127.0.0.1:443"
	cluster_address = "127.0.0.1:8201"
	tls_disable = false
	tls_cert_file = "./certs/server.crt"
	tls_key_file = "./certs/server.key"
	tls_client_ca_file = "./certs/rootca.crt"
	tls_min_version = "tls12"
	tls_require_and_verify_client_cert = true
	tls_disable_client_certs = true
}`))

	var config Config
	list, _ := obj.Node.(*ast.ObjectList)
	objList := list.Filter("listener")
	parseListeners(&config, objList)
	listeners := config.Listeners
	if len(listeners) == 0 {
		t.Fatalf("expected at least one listener in the config")
	}
	listener := listeners[0]
	if listener.Type != "tcp" {
		t.Fatalf("expected tcp listener in the config")
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":                            "127.0.0.1:443",
					"cluster_address":                    "127.0.0.1:8201",
					"tls_disable":                        false,
					"tls_cert_file":                      "./certs/server.crt",
					"tls_key_file":                       "./certs/server.key",
					"tls_client_ca_file":                 "./certs/rootca.crt",
					"tls_min_version":                    "tls12",
					"tls_require_and_verify_client_cert": true,
					"tls_disable_client_certs":           true,
				},
			},
		},
	}

	if !reflect.DeepEqual(config, *expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, *expected)
	}

}
