package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Angelosewase/chatbuddiesgo/Handlers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/go-chi/cors"
)


type ApiCfg struct{
    DB *database.Queries
}

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

	connection, err := sql.Open("mysql", DBURL)
	if err != nil {
		log.Fatal("error opening database:", err)
	}
	ApiConfig := ApiCfg{
		DB: database.New(connection),
	}


	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
    userRouter :=chi.NewRouter()
    userRouter.Post("/signUp", Handlers.SignUpHandler(ApiConfig.DB))
	userRouter.Post("/logIn",Handlers.LoginHandler(ApiConfig.DB))
	userRouter.Get("/logout",Handlers.LogoutHandler)
	router.Mount("/user",userRouter)

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Printf("server running on port %v", PORT)
	err= srv.ListenAndServe()
	if err != nil {
		log.Fatal("error ruunning the serve:", err)
	}

}
