package main // import "github.com/thraxil/sofi

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	port, exists := os.LookupEnv("SOFI_PORT")
	if !exists {
		port = "8080"
	}
	templateDir, exists = os.LookupEnv("SOFI_TEMPLATE_DIR")
	if !exists {
		templateDir = "templates"
	}
	var DB_URL string
	if os.Getenv("DATABASE_URL") != "" {
		DB_URL = os.Getenv("DATABASE_URL")
	} else {
		// local dev settings
		DB_URL = "user=postgres password=postgres dbname=melo_dev sslmode=disable"
	}
	db, err := sqlx.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalln(err)
	}

	s := newSite(db)

	mux := http.NewServeMux()
	addRoutes(mux, s)
	go func() {
		log.Printf("starting server")
		if err := http.ListenAndServe(":"+port, mux); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		// if err := httpServer.Shutdown(shutdownCtx); err != nil {
		// 	fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		// }
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
