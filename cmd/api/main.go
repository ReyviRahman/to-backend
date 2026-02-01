package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))

}
