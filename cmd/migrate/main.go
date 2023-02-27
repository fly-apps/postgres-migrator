package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

type migrationOpts struct {
	sourceURI string
	targetURI string
	noOwner   bool
	clean     bool
	create    bool
	dataOnly  bool
}

func main() {
	var (
		noOwner  = flag.Bool("no-owner", true, "")
		clean    = flag.Bool("clean", true, "")
		create   = flag.Bool("create", true, "")
		dataOnly = flag.Bool("data-only", false, "")
	)

	flag.Parse()

	log.SetFlags(0)

	ctx := context.Background()

	sourceURI := os.Getenv("SOURCE_DATABASE_URI")
	if sourceURI == "" {
		panic("SOURCE_DATABASE_URI secret is required")
	}

	targetURI := os.Getenv("TARGET_DATABASE_URI")
	if targetURI == "" {
		panic("TARGET_DATABASE_URI secret is required")
	}

	opts := migrationOpts{
		sourceURI: sourceURI,
		targetURI: targetURI,
		noOwner:   *noOwner,
		clean:     *clean,
		create:    *create,
		dataOnly:  *dataOnly,
	}

	if err := runPreChecks(ctx, opts); err != nil {
		panic(err)
	}

	if err := runMigration(ctx, opts); err != nil {
		panic(err)
	}
}

func runPreChecks(ctx context.Context, opts migrationOpts) error {
	// Ensure connectivity
	sourceConn, err := openConnection(ctx, opts.sourceURI)
	if err != nil {
		return fmt.Errorf("failed to connect to source: %s", err)
	}
	defer func() { _ = sourceConn.Close(ctx) }()
	log.Println("Source connnection is healthy")

	targetConn, err := openConnection(ctx, opts.targetURI)
	if err != nil {
		return fmt.Errorf("failed to connect to target: %s", err)
	}
	log.Println("Target connnection is healthy")

	defer func() { _ = targetConn.Close(ctx) }()

	// Verify source version is not greater than target
	var sourceVersion string
	if err := sourceConn.QueryRow(ctx, "SHOW server_version;").Scan(&sourceVersion); err != nil {
		return fmt.Errorf("failed to query source version: %s", err)
	}

	log.Println("Source PG version: " + sourceVersion)

	var targetVersion string
	if err := targetConn.QueryRow(ctx, "SHOW server_version;").Scan(&targetVersion); err != nil {
		return fmt.Errorf("failed to query target version: %s", err)
	}

	log.Println("Target PG version: " + targetVersion)

	sourceSlice := strings.Split(sourceVersion, ".")
	targetSlice := strings.Split(targetVersion, ".")

	if sourceSlice[0] > targetSlice[0] {
		return fmt.Errorf("source is running a more recent version than target. expected >= %s, got %s", targetVersion, sourceVersion)
	}

	log.Println("Migration pre-checks have completed without issues.")

	return nil
}

func runMigration(ctx context.Context, opts migrationOpts) error {
	log.Println("Starting migration")

	dumpStr := fmt.Sprintf("pg_dump -d %s", opts.sourceURI)
	if opts.noOwner {
		dumpStr = dumpStr + " --no-owner"
	}
	if opts.clean {
		dumpStr = dumpStr + " --clean"
	}
	if opts.create {
		dumpStr = dumpStr + " --create"
	}
	if opts.dataOnly {
		dumpStr = dumpStr + " --data-only"
	}

	restoreStr := fmt.Sprintf("psql -d %s", opts.targetURI)

	cmd := fmt.Sprintf("%s | %s", dumpStr, restoreStr)

	log.Printf("Running command: %s \n", cmd)

	if _, err := runCommand(cmd); err != nil {
		return err
	}

	log.Println("Postgres migration has completed")

	return nil
}

func openConnection(parentCtx context.Context, uri string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	conf, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	conf.ConnectTimeout = 5 * time.Second

	return pgx.ConnectConfig(ctx, conf)
}

func runCommand(cmdStr string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	return cmd.Output()
}
