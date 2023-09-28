package gorm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	loggerGorm = logger.Setup("infrastructure.database.gorm")
	errGorm    = errors.New("gorm or database error")
	errGormOp  = errors.New("gorm operation error")
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Gorm is the interface for the gorm database connection.
type Gorm interface {
	WhereFirst(model interface{}, query interface{}, args ...interface{}) error
	Create(model interface{}) error
	WhereFind(model interface{}, query interface{}, args ...interface{}) error
}

// Gorm is the struct that contains the gorm database connection.
type GormImpl struct {
	db *gorm.DB
}

// NewGormDB creates a new gorm database connection and returns it.
// If there is an error creating the connection, the application will be stopped.
func NewGormDB(connStr string, maxIdle, maxOpen int) Gorm {
	dsn := convertConnStrToDSN(connStr)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		desc := "Error creating the gorm database"
		loggerGorm.WithError(err).Fatal(desc)

		return nil
	}

	pgDB, err := db.DB()
	if err != nil {
		desc := "Error getting the gorm database"
		loggerGorm.WithError(err).Fatal(desc)

		return nil
	}

	pgDB.SetMaxIdleConns(maxIdle)
	pgDB.SetMaxOpenConns(maxOpen)

	return &GormImpl{db: db}
}

// convertConnStrToDSN converts the connection string to the data source name.
// input example: "postgres://postgres:postgres@localhost:5432/finanger"
// output example: "host=localhost user=postgres password=postgres dbname=finanger port=5432 sslmode=disable TimeZone=UTC"
func convertConnStrToDSN(connStr string) string {
	parts := strings.Split(connStr, "://")
	if len(parts) != 2 {
		loggerGorm.Error("Invalid connection string - missing protocol")
	}

	parts = strings.Split(parts[1], "/")
	if len(parts) != 2 {
		loggerGorm.Error("Invalid connection string - missing database name")
	}

	dbName := parts[1]

	credsAndHost := strings.Split(parts[0], "@")
	if len(credsAndHost) != 2 {
		loggerGorm.Error("Invalid connection string - missing credentials or host")
	}

	creds := strings.Split(credsAndHost[0], ":")
	if len(creds) != 2 {
		loggerGorm.Error("Invalid connection string - missing credentials")
	}

	user := creds[0]
	password := creds[1]

	hostAndPort := strings.Split(credsAndHost[1], ":")
	if len(hostAndPort) != 2 {
		loggerGorm.Error("Invalid connection string - missing host or port")
	}

	host := hostAndPort[0]
	port := hostAndPort[1]

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		dbName,
		port,
	)

	return dsn
}

// WhereFirst is a wrapper for the gorm Where and First methods.
func (g *GormImpl) WhereFirst(model interface{}, query interface{}, args ...interface{}) error {
	if err := g.db.Where(query, args...).First(model).Error; err != nil {
		return fmt.Errorf(crosscuting.WrapLabel, "Error getting the first element", errGormOp, err)
	}

	return nil
}

// Create is a wrapper for the gorm Create method.
func (g *GormImpl) Create(model interface{}) error {
	if err := g.db.Create(model).Error; err != nil {
		return fmt.Errorf(crosscuting.WrapLabel, "Error creating the element", errGormOp, err)
	}

	return nil
}

// WhereFind is a wrapper for the gorm Where and Find methods.
func (g *GormImpl) WhereFind(model interface{}, query interface{}, args ...interface{}) error {
	if err := g.db.Where(query, args...).Find(model).Error; err != nil {
		return fmt.Errorf(crosscuting.WrapLabel, "Error getting the elements", errGormOp, err)
	}

	return nil
}
