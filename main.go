package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
      godotenv.Load()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Failed to load environmental variables(port)")
	}
	DBURL := os.Getenv("DB_URL")
	if DBURL == "" {
		log.Fatal("Failed to load environmental variables(database)")
	}
	_, err := sql.Open("mysql", DBURL)
	if err != nil {
		log.Fatal("error opening database:", err)
	}

	srv := &http.Server{
		Addr: ":" + PORT,
	}

	log.Printf("server running on port %v", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("error ruunning the serve:", err)
	}

}
