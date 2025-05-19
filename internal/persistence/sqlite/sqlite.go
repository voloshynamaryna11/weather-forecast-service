package sqlite

import (
	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"log/slog"
)

type Storage struct {
	db *gorm.DB
}

func New(storagePath string, logger *slog.Logger) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{})
	if err != nil {
		logger.Error("failed to open SQLite database",
			slog.String("op", op),
			slog.String("path", storagePath),
			slog.Any("err", err),
		)
		return nil, err
	}

	logger.Info("SQLite storage initialized",
		slog.String("op", op),
		slog.String("path", storagePath),
	)

	m := gormigrate.New(db, gormigrate.DefaultOptions, initMigrations())

	if err := m.Migrate(); err != nil {
		logger.Error("migrate", "err", err)
		return nil, err
	}
	logger.Info("migrations up-to-date")

	return &Storage{db: db}, nil
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}
