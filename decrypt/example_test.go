package decrypt

import (
	"encoding/json"

	"go.mozilla.org/sops/logging"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("DECRYPT")
}

type configuration struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Age       float64 `json:"age"`
	Address   struct {
		City          string `json:"city"`
		PostalCode    string `json:"postalCode"`
		State         string `json:"state"`
		StreetAddress string `json:"streetAddress"`
	} `json:"address"`
	PhoneNumbers []struct {
		Number string `json:"number"`
		Type   string `json:"type"`
	} `json:"phoneNumbers"`
	AnEmptyValue string `json:"anEmptyValue"`
}

func Example_DecryptFile() {
	var (
		confPath string = "./example.json"
		cfg      configuration
		err      error
	)
	confData, err := File(confPath, "json")
	if err != nil {
		log.Fatalf("cleartext configuration marshalling failed with error: %v", err)
	}
	err = json.Unmarshal(confData, &cfg)
	if err != nil {
		log.Fatalf("cleartext configuration unmarshalling failed with error: %v", err)
	}
	if cfg.FirstName != "John" ||
		cfg.LastName != "Smith" ||
		cfg.Age != 25.4 ||
		cfg.PhoneNumbers[1].Number != "646 555-4567" {
		log.Fatalf("configuration does not contain expected values: %+v", cfg)
	}
	log.Printf("%+v", cfg)
}
