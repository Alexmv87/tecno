package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"tecno/database"
	"tecno/routers"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error al cargar archivo de variables de entorno, se utilizarán las variables del sistema")
	}

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}

	if err := database.Migrate(db); err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	ticketsRouter := routers.NewTicketsRouter(db)

	mux.Handle("/tickets/", http.StripPrefix("/tickets", ticketsRouter.RegisterRoutes()))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "Database connection error", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
