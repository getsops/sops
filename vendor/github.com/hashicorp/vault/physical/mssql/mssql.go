package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// Verify MSSQLBackend satisfies the correct interfaces
var _ physical.Backend = (*MSSQLBackend)(nil)

type MSSQLBackend struct {
	dbTable    string
	client     *sql.DB
	statements map[string]*sql.Stmt
	logger     log.Logger
	permitPool *physical.PermitPool
}

func NewMSSQLBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	username, ok := conf["username"]
	if !ok {
		username = ""
	}

	password, ok := conf["password"]
	if !ok {
		password = ""
	}

	server, ok := conf["server"]
	if !ok || server == "" {
		return nil, fmt.Errorf("missing server")
	}

	port, ok := conf["port"]
	if !ok {
		port = ""
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	var err error
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	} else {
		maxParInt = physical.DefaultParallelOperations
	}

	database, ok := conf["database"]
	if !ok {
		database = "Vault"
	}

	table, ok := conf["table"]
	if !ok {
		table = "Vault"
	}

	appname, ok := conf["appname"]
	if !ok {
		appname = "Vault"
	}

	connectionTimeout, ok := conf["connectiontimeout"]
	if !ok {
		connectionTimeout = "30"
	}

	logLevel, ok := conf["loglevel"]
	if !ok {
		logLevel = "0"
	}

	schema, ok := conf["schema"]
	if !ok || schema == "" {
		schema = "dbo"
	}

	connectionString := fmt.Sprintf("server=%s;app name=%s;connection timeout=%s;log=%s", server, appname, connectionTimeout, logLevel)
	if username != "" {
		connectionString += ";user id=" + username
	}

	if password != "" {
		connectionString += ";password=" + password
	}

	if port != "" {
		connectionString += ";port=" + port
	}

	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, errwrap.Wrapf("failed to connect to mssql: {{err}}", err)
	}

	db.SetMaxOpenConns(maxParInt)

	if _, err := db.Exec("IF NOT EXISTS(SELECT * FROM sys.databases WHERE name = '" + database + "') CREATE DATABASE " + database); err != nil {
		return nil, errwrap.Wrapf("failed to create mssql database: {{err}}", err)
	}

	dbTable := database + "." + schema + "." + table
	createQuery := "IF NOT EXISTS(SELECT 1 FROM " + database + ".INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_NAME='" + table + "' AND TABLE_SCHEMA='" + schema +
		"') CREATE TABLE " + dbTable + " (Path VARCHAR(512) PRIMARY KEY, Value VARBINARY(MAX))"

	if schema != "dbo" {

		var num int
		err = db.QueryRow("SELECT 1 FROM " + database + ".sys.schemas WHERE name = '" + schema + "'").Scan(&num)

		switch {
		case err == sql.ErrNoRows:
			if _, err := db.Exec("USE " + database + "; EXEC ('CREATE SCHEMA " + schema + "')"); err != nil {
				return nil, errwrap.Wrapf("failed to create mssql schema: {{err}}", err)
			}

		case err != nil:
			return nil, errwrap.Wrapf("failed to check if mssql schema exists: {{err}}", err)
		}
	}

	if _, err := db.Exec(createQuery); err != nil {
		return nil, errwrap.Wrapf("failed to create mssql table: {{err}}", err)
	}

	m := &MSSQLBackend{
		dbTable:    dbTable,
		client:     db,
		statements: make(map[string]*sql.Stmt),
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}

	statements := map[string]string{
		"put": "IF EXISTS(SELECT 1 FROM " + dbTable + " WHERE Path = ?) UPDATE " + dbTable + " SET Value = ? WHERE Path = ?" +
			" ELSE INSERT INTO " + dbTable + " VALUES(?, ?)",
		"get":    "SELECT Value FROM " + dbTable + " WHERE Path = ?",
		"delete": "DELETE FROM " + dbTable + " WHERE Path = ?",
		"list":   "SELECT Path FROM " + dbTable + " WHERE Path LIKE ?",
	}

	for name, query := range statements {
		if err := m.prepare(name, query); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *MSSQLBackend) prepare(name, query string) error {
	stmt, err := m.client.Prepare(query)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to prepare %q: {{err}}", name), err)
	}

	m.statements[name] = stmt

	return nil
}

func (m *MSSQLBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"mssql", "put"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, err := m.statements["put"].Exec(entry.Key, entry.Value, entry.Key, entry.Key, entry.Value)
	if err != nil {
		return err
	}

	return nil
}

func (m *MSSQLBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"mssql", "get"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	var result []byte
	err := m.statements["get"].QueryRow(key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	ent := &physical.Entry{
		Key:   key,
		Value: result,
	}

	return ent, nil
}

func (m *MSSQLBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"mssql", "delete"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	_, err := m.statements["delete"].Exec(key)
	if err != nil {
		return err
	}

	return nil
}

func (m *MSSQLBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"mssql", "list"}, time.Now())

	m.permitPool.Acquire()
	defer m.permitPool.Release()

	likePrefix := prefix + "%"
	rows, err := m.statements["list"].Query(likePrefix)
	if err != nil {
		return nil, err
	}
	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, errwrap.Wrapf("failed to scan rows: {{err}}", err)
		}

		key = strings.TrimPrefix(key, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			keys = append(keys, key)
		} else if i != -1 {
			keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
		}
	}

	sort.Strings(keys)

	return keys, nil
}
