package audit

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os/user"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/logging"
	"gopkg.in/yaml.v2"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("AUDIT")
	confBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.WithField("error", err).Debugf("Error reading config")
		return
	}
	var conf config
	err = yaml.Unmarshal(confBytes, &conf)
	if err != nil {
		log.WithField("error", err).Panicf("Error unmarshalling config")
	}
	// If we are running test, then don't create auditors.
	// This is pretty hacky, but doing it The Right Way would require
	// restructuring SOPS to use dependency injection instead of just using
	// globals everywhere.
	if flag.Lookup("test.v") != nil {
		return
	}
	for _, pgConf := range conf.Backends.Postgres {
		auditors = append(auditors, NewPostgresAuditor(pgConf.ConnStr))
	}
}

// TODO: Make platform agnostic
const configFile = "/etc/sops/audit.yaml"

type config struct {
	Backends struct {
		Postgres []struct {
			ConnStr string `yaml:"connection_string"`
		} `yaml:"postgres"`
	} `yaml:"backends"`
}

var auditors []Auditor

func SubmitEvent(event interface{}) {
	for _, auditor := range auditors {
		auditor.Handle(event)
	}
}

func Register(auditor Auditor) {
	auditors = append(auditors, auditor)
}

type Auditor interface {
	Handle(event interface{})
}

type DecryptEvent struct {
	File string
}

type EncryptEvent struct {
	File string
}

type RotateEvent struct {
	File string
}

type PostgresAuditor struct {
	DB *sql.DB
}

func NewPostgresAuditor(connStr string) *PostgresAuditor {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatalf("Pinging audit database failed: %s", err)
	} else if result != 1 {
		log.Fatalf("Database malfunction: SELECT 1 should return 1, but returned %d", result)
	}
	return &PostgresAuditor{
		DB: db,
	}
}

func (p *PostgresAuditor) Handle(event interface{}) {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user for auditing: %s", err)
	}
	switch event := event.(type) {
	case DecryptEvent:
		// Save the event to the database
		log.WithField("file", event.File).
			Debug("Saving decrypt event to database")
		_, err = p.DB.Exec("INSERT INTO audit_event (action, username, file) VALUES ($1, $2, $3)", "decrypt", u.Username, event.File)
		if err != nil {
			log.Fatalf("Failed to insert audit record: %s", err)
		}
	case EncryptEvent:
		// Save the event to the database
		log.WithField("file", event.File).
			Debug("Saving encrypt event to database")
		_, err = p.DB.Exec("INSERT INTO audit_event (action, username, file) VALUES ($1, $2, $3)", "encrypt", u.Username, event.File)
		if err != nil {
			log.Fatalf("Failed to insert audit record: %s", err)
		}
	case RotateEvent:
		// Save the event to the database
		log.WithField("file", event.File).
			Debug("Saving rotate event to database")
		_, err = p.DB.Exec("INSERT INTO audit_event (action, username, file) VALUES ($1, $2, $3)", "rotate", u.Username, event.File)
		if err != nil {
			log.Fatalf("Failed to insert audit record: %s", err)
		}
	default:
		log.WithField("type", fmt.Sprintf("%T", event)).
			Info("Received unknown event")
	}
}
