package test

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	pg "github.com/ermyar/pg-util/internal/postgres"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestRemove(t *testing.T) {
	for _, tc := range []struct {
		name    string
		conf    pg.Config
		args    []string
		deleted []string
	}{
		{
			name:    "simple",
			conf:    pgCfg,
			args:    []string{"remove_db"},
			deleted: []string{"remove_db"},
		},
		{
			name:    "complex",
			conf:    pgCfg,
			args:    []string{"remove_db1", "remove_db2", "remove_db3"},
			deleted: []string{"remove_db1", "remove_db2", "remove_db3"},
		},
		{
			name:    "not exist",
			conf:    pgCfg,
			args:    []string{"tmp"},
			deleted: []string(nil),
		},
		{
			name: "regex",
			conf: pgCfg,
			args: []string{"remove_regular", "remove_regular_(a|b).+"},
			deleted: []string{
				"remove_regular",
				"remove_regular_a1",
				"remove_regular_a2",
				"remove_regular_a3_wow",
				"remove_regular_b_extra",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := tc.conf.Connect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			before, err := pg.ListDatabases(context.Background(), conn)
			if err != nil {
				t.Fatal(err)
			}

			remove := getCmd("remove", strings.Join(tc.args, ","), tc.conf)
			if err := remove.Run(); err != nil {
				t.Fatal(err)
			}

			after, err := pg.ListDatabases(context.Background(), conn)
			if err != nil {
				t.Fatal(err)
			}

			diff := getDiff(before, after)
			require.Equal(t, tc.deleted, diff)
		})
	}
}

func TestBackup(t *testing.T) {
	for _, tc := range []struct {
		name    string
		conf    pg.Config
		dbname  string
		restore string
		tables  []string
	}{
		{
			name:    "simple",
			conf:    pgCfg,
			dbname:  "test",
			restore: "restored",
			tables:  []string{"users", "numbers"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// backuping.
			backup := getCmd("backup", tc.dbname, tc.conf)
			if err := backup.Run(); err != nil {
				t.Fatal(err)
			}
			backupFilename := "backup_" + tc.dbname + ".sql"
			require.FileExists(t, backupFilename)

			var before [][]string
			for _, table := range tc.tables {
				data, err := getValuesFromTable(context.Background(), tc.conf, tc.dbname, table)
				require.NoError(t, err)
				before = append(before, data)
			}

			// restoring.
			restore := exec.Command("psql", "-f", backupFilename)
			restore.Args = append(restore.Args,
				"-d", pgCfg.ConnString()+"/"+tc.restore,
			)
			if err := restore.Run(); err != nil {
				t.Fatal(err)
			}

			var after [][]string
			for _, table := range tc.tables {
				data, err := getValuesFromTable(context.Background(), tc.conf, tc.restore, table)
				require.NoError(t, err)
				after = append(after, data)
			}

			require.Equal(t, before, after)
			os.Remove(backupFilename)
		})
	}
}
