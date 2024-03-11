package sqls

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

var o sync.Once

const FileName = "rest-sqlite.db"

type SQLiteClient struct {
	DB *gorm.DB
}

var err error
var sqliteClient *SQLiteClient

func InitGORMSQLiteDB() (*SQLiteClient, error) {
	o.Do(func() {
		if _, err = os.Stat(FileName); err == nil {
			err = os.Remove(FileName)
			if err != nil {
				log.Debugf("unable to remove database file, %v", err)
				os.Exit(1)
			}
		}

		var db *gorm.DB
		db, err = gorm.Open(sqlite.Open(FileName), &gorm.Config{})
		if err != nil {
			log.Debugf("database connection error, %v", err)
			os.Exit(1)
		}
		serviceName := os.Getenv("SERVICE_NAME")
		collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if len(serviceName) > 0 && len(collectorURL) > 0 {
			if err := db.Use(otelgorm.NewPlugin()); err != nil {
				log.Debugf("unable to attach opentel plugin error, %v", err)
				os.Exit(1)
			}
		}

		sqliteClient = &SQLiteClient{
			DB: db,
		}

	})

	return sqliteClient, err
}
