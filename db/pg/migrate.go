package pg

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

func Migrate(fs *embed.FS, path string) error {
	sourceInstance, err := httpfs.New(http.FS(fs), path)
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	defer sourceInstance.Close()

	targetInstance, err := postgres.WithInstance(Conn().DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance, %w", err)
	}
	m, err := migrate.NewWithInstance(
		"httpfs", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
