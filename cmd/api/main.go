package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ReyviRahman/to-backend/internal/db"
	"github.com/ReyviRahman/to-backend/internal/store"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func getEnvInt(key string) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		log.Fatalf("CRITICAL ERROR: Environment variable '%s' wajib diisi!", key)
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("CRITICAL ERROR: Environment variable '%s' harus berupa ANGKA (integer). Nilai saat ini: %s", key, valStr)
	}
	return valInt
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
		db: dbConfig{
			addr:         os.Getenv("DB_DSN"),
			maxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS"),
			maxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS"),
			maxIdleTime:  os.Getenv("DB_MAX_IDLE_TIME"),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
