package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/helper/dbtxn"
	"github.com/lib/pq"
	"github.com/ory/dockertest"
)

var (
	testPostgresImagePull sync.Once
)

func preparePostgresTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("PG_URL") != "" {
		return func() {}, os.Getenv("PG_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"})
	if err != nil {
		t.Fatalf("Could not start local PostgreSQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("postgres://postgres:secret@localhost:%s/database?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("postgres", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to PostgreSQL docker container: %s", err)
	}

	return
}

func TestPostgreSQL_Initialize(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": 5,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test decoding a string value for max_open_connections
	connectionDetails = map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": "5",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestPostgreSQL_CreateUser(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testPostgresRole},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	statements.Creation = []string{testPostgresReadOnlyRole}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep to make sure we haven't expired if granularity is only down to the second
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestPostgreSQL_RenewUser(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testPostgresRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
	statements.Renewal = []string{defaultPostgresRenewSQL}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Sleep longer than the initial expiration time
	time.Sleep(2 * time.Second)

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

}

func TestPostgreSQL_RotateRootCredentials(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connURL = strings.Replace(connURL, "postgres:secret", `{{username}}:{{password}}`, -1)

	connectionDetails := map[string]interface{}{
		"connection_url":       connURL,
		"max_open_connections": 5,
		"username":             "postgres",
		"password":             "secret",
	}

	db := new()

	connProducer := db.SQLConnectionProducer

	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	newConf, err := db.RotateRootCredentials(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if newConf["password"] == "secret" {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostgreSQL_RevokeUser(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testPostgresRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}

	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(2*time.Second))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statements
	statements.Revocation = []string{defaultPostgresRevocationSQL}
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestPostgresSQL_SetCredentials(t *testing.T) {
	cleanup, connURL := preparePostgresTestContainer(t)
	defer cleanup()

	// create the database user
	dbUser := "vaultstatictest"
	createTestPGUser(t, connURL, dbUser, "password", testRoleStaticCreate)

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GenerateCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	usernameConfig := dbplugin.StaticUserConfig{
		Username: dbUser,
		Password: password,
	}

	// Test with no configured Rotation Statement
	username, password, err := db.SetCredentials(context.Background(), dbplugin.Statements{}, usernameConfig)
	if err == nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Rotation: []string{testPostgresStaticRoleRotate},
	}
	// User should not exist, make sure we can create
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// call SetCredentials again, password will change
	newPassword, _ := db.GenerateCredentials(context.Background())
	usernameConfig.Password = newPassword
	username, password, err = db.SetCredentials(context.Background(), statements, usernameConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if password != newPassword {
		t.Fatal("passwords should have changed")
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	t.Helper()
	// Log in with the new creds
	connURL = strings.Replace(connURL, "postgres:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

const testPostgresRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testPostgresReadOnlyRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO "{{name}}";
`

const testPostgresBlockStatementRole = `
DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$

CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
GRANT "foo-role" TO "{{name}}";
ALTER ROLE "{{name}}" SET search_path = foo;
GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";
`

var testPostgresBlockStatementRoleSlice = []string{
	`
DO $$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_roles WHERE rolname='foo-role') THEN
      CREATE ROLE "foo-role";
      CREATE SCHEMA IF NOT EXISTS foo AUTHORIZATION "foo-role";
      ALTER ROLE "foo-role" SET search_path = foo;
      GRANT TEMPORARY ON DATABASE "postgres" TO "foo-role";
      GRANT ALL PRIVILEGES ON SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA foo TO "foo-role";
      GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA foo TO "foo-role";
   END IF;
END
$$
`,
	`CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';`,
	`GRANT "foo-role" TO "{{name}}";`,
	`ALTER ROLE "{{name}}" SET search_path = foo;`,
	`GRANT CONNECT ON DATABASE "postgres" TO "{{name}}";`,
}

const defaultPostgresRevocationSQL = `
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM "{{name}}";
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM "{{name}}";
REVOKE USAGE ON SCHEMA public FROM "{{name}}";

DROP ROLE IF EXISTS "{{name}}";
`

const testPostgresStaticRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

const testRoleStaticCreate = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}';
`

const testPostgresStaticRoleRotate = `
ALTER ROLE "{{name}}" WITH PASSWORD '{{password}}';
`

const testPostgresStaticRoleGrant = `
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`

// This is a copy of a test helper method also found in
// builtin/logical/database/rotation_test.go , and should be moved into a shared
// helper file in the future.
func createTestPGUser(t *testing.T, connURL string, username, password, query string) {
	t.Helper()
	conn, err := pq.ParseURL(connURL)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", conn)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Start a transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	m := map[string]string{
		"name":     username,
		"password": password,
	}
	if err := dbtxn.ExecuteTxQuery(ctx, tx, m, query); err != nil {
		t.Fatal(err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
