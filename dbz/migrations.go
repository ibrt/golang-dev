package dbz

import (
	"context"
	"embed"
	"fmt"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/memz"
	"github.com/ibrt/golang-utils/outz"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"

	"github.com/ibrt/golang-dev/consolez"
)

// MigrationsConfig describes the configuration for migrations.
type MigrationsConfig struct {
	PostgresURL string
	TableName   string
	FS          embed.FS
}

// MustApplyMigrations applies all migrations.
func MustApplyMigrations(cfg *MigrationsConfig) {
	consolez.DefaultCLI.Notice("apply-migrations", "loading migrations...")
	m, f := mustGetMigrator(cfg)
	defer f()

	beforeVersion, err := m.GetCurrentVersion(context.Background())
	errorz.MaybeMustWrap(err)
	targetVersion := int32(len(m.Migrations))

	consolez.DefaultCLI.Notice("apply-migrations", "executing migrations...", fmt.Sprintf("%v → %v", beforeVersion, targetVersion))
	errorz.MaybeMustWrap(m.Migrate(context.Background()))

	afterVersion, err := m.GetCurrentVersion(context.Background())
	errorz.MaybeMustWrap(err)
	mustShow(m.Migrations, beforeVersion, afterVersion)
}

// MustRollBackMigrations rolls migrations back.
func MustRollBackMigrations(cfg *MigrationsConfig, targetVersion int32) {
	consolez.DefaultCLI.Notice("roll-back-migrations", "loading migrations...")
	m, f := mustGetMigrator(cfg)
	defer f()

	beforeVersion, err := m.GetCurrentVersion(context.Background())
	errorz.MaybeMustWrap(err)
	errorz.Assertf(targetVersion < beforeVersion, "target version is not lower than current version: %v → %v", beforeVersion, targetVersion)

	consolez.DefaultCLI.Notice("roll-back-migrations", "executing migrations...", fmt.Sprintf("%v → %v", beforeVersion, targetVersion))
	errorz.MaybeMustWrap(m.MigrateTo(context.Background(), targetVersion))

	afterVersion, err := m.GetCurrentVersion(context.Background())
	errorz.MaybeMustWrap(err)
	mustShow(m.Migrations, beforeVersion, afterVersion)
}

// MustShowMigrations shows migrations.
func MustShowMigrations(cfg *MigrationsConfig) {
	consolez.DefaultCLI.Notice("show-migrations", "loading migrations...")
	m, f := mustGetMigrator(cfg)
	defer f()

	currentVersion, err := m.GetCurrentVersion(context.Background())
	errorz.MaybeMustWrap(err)

	mustShow(m.Migrations, currentVersion, currentVersion)
}

func mustShow(ms []*migrate.Migration, beforeVersion, afterVersion int32) {
	fmt.Println()

	consolez.DefaultCLI.NewTable("Version", "Migration", "Status").
		SetRows(memz.TransformSlice(ms, func(_ int, i *migrate.Migration) []string {
			return []string{
				fmt.Sprintf("%v", i.Sequence),
				i.Name,
				func() string {
					if beforeVersion <= afterVersion {
						if i.Sequence <= beforeVersion {
							return outz.DefaultStyles.Secondary().Sprintf("synced, untouched")
						}

						if i.Sequence <= afterVersion {
							return outz.DefaultStyles.Success().Sprintf("applied")
						}

						return outz.DefaultStyles.Secondary().Sprintf("not synced, untouched")
					}

					if i.Sequence <= afterVersion {
						return outz.DefaultStyles.Secondary().Sprintf("synced, untouched")
					}

					if i.Sequence <= beforeVersion {
						return outz.DefaultStyles.Warning().Sprintf("rolled back")
					}

					return outz.DefaultStyles.Secondary().Sprintf("not synced, untouched")
				}(),
			}
		})).
		Print()
	fmt.Println()
}

func mustGetMigrator(cfg *MigrationsConfig) (*migrate.Migrator, func()) {
	pg, err := pgx.Connect(context.Background(), cfg.PostgresURL)
	errorz.MaybeMustWrap(err)

	m, err := migrate.NewMigrator(context.Background(), pg, cfg.TableName)
	errorz.MaybeMustWrap(err)
	errorz.MaybeMustWrap(m.LoadMigrations(cfg.FS))

	return m, func() {
		errorz.MaybeMustWrap(pg.Close(context.Background()))
	}
}
