// +build go1.9

package pq

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestPing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	db := openTestConn(t)
	defer db.Close()

	if _, ok := reflect.TypeOf(db).MethodByName("Conn"); !ok {
		t.Skipf("Conn method undefined on type %T, skipping test (requires at least go1.9)", db)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatal("expected Ping to succeed")
	}
	defer cancel()

	// grab a connection
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// start a transaction and read backend pid of our connection
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	rows, err := tx.Query("SELECT pg_backend_pid()")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	// read the pid from result
	var pid int
	for rows.Next() {
		if err := rows.Scan(&pid); err != nil {
			t.Fatal(err)
		}
	}
	if rows.Err() != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	// kill the process which handles our connection and test if the ping fails
	if _, err := db.Exec("SELECT pg_terminate_backend($1)", pid); err != nil {
		t.Fatal(err)
	}
	if err := conn.PingContext(ctx); err != driver.ErrBadConn {
		t.Fatalf("expected error %s, instead got %s", driver.ErrBadConn, err)
	}
}
