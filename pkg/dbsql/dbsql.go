package dbsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var dbonce sync.Once

func ConnectDatabase(options DBConnectOptions, connector func(string) (*sql.DB, error)) (SqlBackend, error) {
	var db *sql.DB
	var err error

	if db, err = connector(options.FilenameOrURL); err != nil {
		return db, err
	}

	if db == nil {
		return db, fmt.Errorf("database connection failed")
	}

	if options.VerboseLogging {
		loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
		db = sqldblogger.OpenDriver(options.FilenameOrURL, db.Driver(), loggerAdapter)

		if db == nil {
			return db, fmt.Errorf("database logger failed")
		}
	}

	return db, nil
}

func UseSQLite(dbname string) (*sql.DB, error) {
	if dbname == "" {
		return nil, fmt.Errorf("no database name provided")
	}

	// NOTE: Extensions should only be registered once
	dbonce.Do(func() {
		sql.Register("sqlite3_ext",
			&sqlite3.SQLiteDriver{
				ConnectHook: func(conn *sqlite3.SQLiteConn) error {
					// if _, err := conn.Exec(`PRAGMA busy_timeout = 5000`, nil); err != nil { return err }

					// if _, err := conn.Exec("PRAGMA journal_mode=WAL", nil); err != nil { return err }

					// RegEx Match
					// 	err = db.QueryRow(`SELECT regexp("foo.*", "seafood")`).Scan(&i)
					var regEx = func(re, s string) (bool, error) { return regexp.MatchString(re, s) }
					if err := conn.RegisterFunc("regexp", regEx, true); err != nil {
						return err
					}

					return nil
				},
			},
		)
	})

	options := []string{"fts5=on", "_cslike=off", "_fk=on", "_ignore_check_constraints=off", "_journal=WAL"}
	d, err := sql.Open("sqlite3_ext", fmt.Sprintf("%s?%s", dbname, strings.Join(options, "&")))
	if err != nil {
		return nil, err
	}

	return d, nil
}

func UseLibSQL(url string) (*sql.DB, error) { return nil, fmt.Errorf("libsql is not supported yet") }

type (
	SqlBackend interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...any) *sql.Row
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	}

	DBConnectOptions struct {
		VerboseLogging bool
		FilenameOrURL  string // i.e "db.sqlite3"
	}

	DatabaseConfig struct {
		Path                    string // Database path (i.e /root/sam/cache)
		Filename                string // Database filename (i.e db.sqlite3)
		MigrationAssertionQuery string // Query to check if database schema is up-to-date
		MigrationExecQuery      string // Query to migrate/upgrade database schema
		IsVerboseEnabled        bool
	}
)
