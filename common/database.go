package common

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// GetDBPath returns the database path from environment or default.
// Exported for use in tests.
func GetDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/gorm.db"
	}
	return dbPath
}

// GetTestDBPath returns the test database path from environment or default.
// Exported for use in tests.
func GetTestDBPath() string {
	testDBPath := os.Getenv("TEST_DB_PATH")
	if testDBPath == "" {
		testDBPath = "./data/gorm_test.db"
	}
	return testDBPath
}

func getPostgresDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("DATABASE_URL")); dsn != "" {
		return dsn
	}
	if dsn := strings.TrimSpace(os.Getenv("DB_DSN")); dsn != "" {
		return dsn
	}

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	if parsed, err := url.Parse(host); err == nil && parsed.Host != "" {
		host = parsed.Hostname()
	} else {
		host = strings.Trim(host, "/[]")
	}
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := os.Getenv("DB_PASSWORD")
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	if host == "" || user == "" || password == "" || name == "" {
		return ""
	}

	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	if port == "" {
		port = "5432"
	}
	sslMode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslMode == "" {
		sslMode = "disable"
	}
	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}
	query := connectionURL.Query()
	query.Set("sslmode", sslMode)
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String()
}

// ensureDir creates the directory for the database file if it doesn't exist
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		return os.MkdirAll(dir, 0750)
	}
	return nil
}

// Init opens the database selected by DB_DRIVER and stores its reference.
// DB_DRIVER accepts "sqlite" (the default), "mysql", "mariadb", or "postgres".
// PostgreSQL also supports DATABASE_URL and DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME.
func Init() (*gorm.DB, error) {
	DB = nil
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("DB_DRIVER")))
	if driver == "" {
		driver = "sqlite"
	}

	var (
		db  *gorm.DB
		err error
	)
	switch driver {
	case "sqlite":
		dbPath := GetDBPath()
		if err := ensureDir(dbPath); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	case "mysql", "mariadb":
		dsn := strings.TrimSpace(os.Getenv("DB_DSN"))
		if dsn == "" {
			return nil, fmt.Errorf("DB_DSN is required when DB_DRIVER is %q", driver)
		}
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres", "postgresql":
		dsn := getPostgresDSN()
		if dsn == "" {
			return nil, fmt.Errorf("DATABASE_URL or DB_DSN, or DB_HOST, DB_USER, DB_PASSWORD, and DB_NAME are required when DB_DRIVER is %q", driver)
		}
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q (expected sqlite, mysql, mariadb, or postgres)", driver)
	}
	if err != nil {
		return nil, fmt.Errorf("open %s database: %w", driver, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection pool: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping %s database: %w", driver, err)
	}
	DB = db
	return DB, nil
}

// This function will create a temporarily database for running testing cases
func TestDBInit() *gorm.DB {
	testDBPath := GetTestDBPath()

	// Ensure the directory exists
	if err := ensureDir(testDBPath); err != nil {
		fmt.Println("db err: (TestDBInit - create dir) ", err)
	}

	test_db, err := gorm.Open(sqlite.Open(testDBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("db err: (TestDBInit) ", err)
	}
	sqlDB, err := test_db.DB()
	if err != nil {
		fmt.Println("db err: (TestDBInit - get sql.DB) ", err)
	} else {
		sqlDB.SetMaxIdleConns(3)
	}
	DB = test_db
	return DB
}

// Delete the database after running testing cases.
func TestDBFree(test_db *gorm.DB) error {
	sqlDB, err := test_db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	testDBPath := GetTestDBPath()
	err = os.Remove(testDBPath)
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
