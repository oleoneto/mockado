package types

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/cleopatrio/mockado/pkg/dbsql"
	"github.com/spf13/cobra"
)

var dbOnce sync.Once

func DatabaseConnectHook(config dbsql.DatabaseConfig) {
	filepath := fmt.Sprintf("%s/%s", config.Path, config.Filename)

	var err error
	var database dbsql.SqlBackend

	dbOnce.Do(func() {
		database, err = dbsql.ConnectDatabase(
			dbsql.DBConnectOptions{
				VerboseLogging: config.IsVerboseEnabled,
				FilenameOrURL:  filepath,
			},
			dbsql.UseSQLite,
		)

		if err != nil {
			panic(err)
		}

		ctx := context.TODO()
		row := database.QueryRowContext(ctx, config.MigrationAssertionQuery)

		// MARK: Migrate initial schema for SQLite database
		if errors.Is(row.Scan(nil), sql.ErrNoRows) {
			if _, err := database.ExecContext(ctx, config.MigrationExecQuery); err != nil {
				if x, _ := regexp.Match("no such table*", []byte(err.Error())); x {
					return
				}
				panic(err)
			}
		}
	})
}

func BeforeHook(f *CLIFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		f.SetFormatter(cmd, args)
		if !f.IsTimingEnabled {
			return
		}
		f.ExecutionStartTime = time.Now()
	}
}

func AfterHook(f *CLIFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if !f.IsTimingEnabled {
			return
		}
		fmt.Fprintln(
			os.Stderr,
			append([]any{"Elapsed time:", time.Since(f.ExecutionStartTime)}, f.ExecutionExitLog...)...,
		)
	}
}
