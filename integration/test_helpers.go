package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"

	pg "github.com/ermyar/pg-util/internal/postgres"
)

var binaryPath string
var pgCfg pg.Config
var pgFinisher func()

var ErrUnableToConnect = errors.New("unable to connect to Postgres")

func setup() {
	buildBinary()
	startPostgres()
}

func teardown() {
	if binaryPath != "" {
		os.Remove(binaryPath)
	}
	if pgFinisher != nil {
		pgFinisher()
	}
}

func buildBinary() {
	cmd := exec.Command("go", "build", "-o", "pg-util-test", "./cmd/...")
	cmd.Dir = filepath.Join("..")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build binary: %s", string(out))
		os.Exit(1)
	}

	binaryPath = filepath.Join("..", "pg-util-test")
}

func startPostgres() {
	cmd := exec.Command("docker", "compose", "up", "-d", "postgres-test")
	cmd.Dir = filepath.Join("..")
	cmd.Run()

	pgCfg = pg.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "test",
		Password: "test",
	}

	waitDB(60 * time.Second)

	pgFinisher = func() {
		cmd := exec.Command("docker", "compose", "down", "postgres-test")
		cmd.Dir = filepath.Join("..")
		cmd.Run()
	}
}

// getDiff returns elements of a that are not in b,
// where b is subset of a.
func getDiff(a, b []string) []string {
	slices.Sort(a)
	slices.Sort(b)
	var ans []string
	for i, j := 0, 0; i < len(a) && j < len(b); {
		if a[i] == b[j] {
			i++
			j++
		} else if a[i] < b[j] {
			ans = append(ans, a[i])
			i++
		}
	}
	return ans
}

func waitDB(timeout time.Duration) error {
	ch := time.After(timeout)
	for {
		select {
		case <-ch:
			return fmt.Errorf("err: %w, timeout: %d", ErrUnableToConnect, timeout)
		default:
			if _, err := pgCfg.Connect(context.Background()); err == nil {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func getCmd(operation, databases string, conf pg.Config) *exec.Cmd {
	cmd := exec.Command(binaryPath, "--databases", databases)
	cmd.Args = append(cmd.Args, "--operation", operation)
	cmd.Args = append(cmd.Args,
		"--host", conf.Host,
		"--port", conf.Port,
		"--username", conf.User,
		"--password", conf.Password,
	)
	return cmd
}

func getValuesFromTable(ctx context.Context, cfg pg.Config, table string) ([]string, error) {
	conn, err := cfg.Connect(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, err
	}
	var ans []string
	for rows.Next() {
		data, err := rows.Values()
		if err != nil {
			return nil, err
		}
		ans = append(ans, fmt.Sprintf("%v", data...))
	}
	return ans, nil
}
