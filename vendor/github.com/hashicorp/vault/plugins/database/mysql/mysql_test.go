package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/credsutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/ory/dockertest"
)

var _ dbplugin.Database = (*MySQL)(nil)

func prepareMySQLTestContainer(t *testing.T, legacy bool) (cleanup func(), retURL string) {
	if os.Getenv("MYSQL_URL") != "" {
		return func() {}, os.Getenv("MYSQL_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	imageVersion := "5.7"
	if legacy {
		imageVersion = "5.6"
	}

	resource, err := pool.Run("mysql", imageVersion, []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		t.Fatalf("Could not start local MySQL docker container: %s", err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	retURL = fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		var db *sql.DB
		db, err = sql.Open("mysql", retURL)
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to MySQL docker container: %s", err)
	}

	return
}

func TestMySQL_Initialize(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
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

func TestMySQL_CreateUser(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-long-displayname",
		RoleName:    "test-long-rolename",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test a second time to make sure usernames don't collide
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test with a manually prepare statement
	statements.Creation = []string{testMySQLRolePreparedStmt}

	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

}

func TestMySQL_CreateUser_Legacy(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, true)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(credsutil.NoneLength, LegacyMetadataLen, LegacyUsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test-long-displayname",
		RoleName:    "test-long-rolename",
	}

	// Test with no configured Creation Statement
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConfig, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("Expected error when no creation statement is provided")
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test a second time to make sure usernames don't collide
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMySQL_RotateRootCredentials(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connURL = strings.Replace(connURL, "root:secret", `{{username}}:{{password}}`, -1)

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
		"username":       "root",
		"password":       "secret",
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
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

func TestMySQL_RevokeUser(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testMySQLRoleWildCard},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
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

	statements.Creation = []string{testMySQLRoleWildCard}
	username, password, err = db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test custom revoke statements
	statements.Revocation = []string{testMySQLRevocationSQL}
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}

func TestMySQL_SetCredentials(t *testing.T) {
	cleanup, connURL := prepareMySQLTestContainer(t, false)
	defer cleanup()

	// create the database user and verify we can access
	dbUser := "vaultstatictest"
	createTestMySQLUser(t, connURL, dbUser, "password", testRoleStaticCreate)
	if err := testCredsExist(t, connURL, dbUser, "password"); err != nil {
		t.Fatalf("Could not connect with credentials: %s", err)
	}

	connectionDetails := map[string]interface{}{
		"connection_url": connURL,
	}

	db := new(MetadataLen, MetadataLen, UsernameLen)
	_, err := db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	newPassword, err := db.GenerateCredentials(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	userConfig := dbplugin.StaticUserConfig{
		Username: dbUser,
		Password: newPassword,
	}

	statements := dbplugin.Statements{
		Rotation: []string{testRoleStaticRotate},
	}

	_, _, err = db.SetCredentials(context.Background(), statements, userConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// verify new password works
	if err := testCredsExist(t, connURL, dbUser, newPassword); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// call SetCredentials again, password will change
	newPassword, _ = db.GenerateCredentials(context.Background())
	userConfig.Password = newPassword
	_, _, err = db.SetCredentials(context.Background(), statements, userConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := testCredsExist(t, connURL, dbUser, newPassword); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func testCredsExist(t testing.TB, connURL, username, password string) error {
	// Log in with the new creds
	connURL = strings.Replace(connURL, "root:secret", fmt.Sprintf("%s:%s", username, password), 1)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

func createTestMySQLUser(t *testing.T, connURL, username, password, query string) {
	t.Helper()
	db, err := sql.Open("mysql", connURL)
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

	// copied from mysql.go
	for _, query := range strutil.ParseArbitraryStringSlice(query, ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}
		query = dbutil.QueryHelper(query, map[string]string{
			"name":     username,
			"password": password,
		})

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			if e, ok := err.(*stdmysql.MySQLError); ok && e.Number == 1295 {
				_, err = tx.ExecContext(ctx, query)
				if err != nil {
					t.Fatal(err)
				}
				stmt.Close()
				continue
			}

			t.Fatal(err)
		}
		if _, err := stmt.ExecContext(ctx); err != nil {
			stmt.Close()
			t.Fatal(err)
		}
		stmt.Close()
	}
}

const testMySQLRolePreparedStmt = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
set @grants=CONCAT("GRANT SELECT ON ", "*", ".* TO '{{name}}'@'%'");
PREPARE grantStmt from @grants;
EXECUTE grantStmt;
DEALLOCATE PREPARE grantStmt;
`
const testMySQLRoleWildCard = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{name}}'@'%';
`
const testMySQLRevocationSQL = `
REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'%'; 
DROP USER '{{name}}'@'%';
`

const testRoleStaticCreate = `
CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
GRANT SELECT ON *.* TO '{{name}}'@'%';
`

const testRoleStaticRotate = `
ALTER USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';
`
