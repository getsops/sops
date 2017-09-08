package audit

import (
	"log"
	"os/user"

	"database/sql"

	_ "github.com/lib/pq"
)

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

// postgresConnStr should be defined at compile time with the -X ldflag
var postgresConnStr string

type PostgresAuditor struct {
	DB *sql.DB
}

func NewPostgresAuditor() *PostgresAuditor {
	db, err := sql.Open("postgres", postgresConnStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Pinging audit database failed: %s", err)
	}
	return &PostgresAuditor{
		DB: db,
	}
}

func (p *PostgresAuditor) Handle(event interface{}) {
	switch event := event.(type) {
	case DecryptEvent:
		// Save the event to the database
		log.Printf("Saving decrypt event for file %s to database", event.File)
		u, err := user.Current()
		if err != nil {
			log.Fatalf("Error getting current user for auditing: %s", err)
		}
		_, err = p.DB.Exec("INSERT INTO decrypt_event (username, file) VALUES ($1, $2)", u.Username, event.File)
		if err != nil {
			log.Fatalf("Failed to insert audit record: %s", err)
		}
	default:
		log.Printf("Received event of unknown type %T", event)
	}
}
