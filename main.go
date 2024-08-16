package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Angelosewase/chatbuddiesgo/Handlers"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/Angelosewase/chatbuddiesgo/middleware"
	"github.com/Angelosewase/chatbuddiesgo/socket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ApiCfg struct {
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
	userRouter := chi.NewRouter()
	chatRouter := chi.NewRouter()
	chatRouter.Use(middleware.AuthMiddleware)

	userRouter.Post("/signUp", Handlers.SignUpHandler(ApiConfig.DB))
	userRouter.Post("/logIn", Handlers.LoginHandler(ApiConfig.DB))
	userRouter.Get("/logout", Handlers.LogoutHandler)
	router.Mount("/user", userRouter)

	chatRouter.Get("/chats", Handlers.GetChatHandler(ApiConfig.DB))
	chatRouter.Post("/newChat", Handlers.CreateChatHandler(ApiConfig.DB))
	chatRouter.Delete("/deleteChat", Handlers.DeleteChatHandler(ApiConfig.DB))

	socketServer := &socket.Server{}
	if err := socketServer.NewServer(); err != nil {
		log.Fatalf("Failed to create socket server: %v", err)
	}

	chatRouter.Handle("/message", socketServer.SocketServer)

	if err := socketServer.Start(); err != nil {
		log.Fatalf("Failed to start socket server: %v", err)
	}

	router.Mount("/chat", chatRouter)

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Printf("server running on port %v", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("error ruunning the serve:", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := srv.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close the Socket.IO server
	socketServer.Close()
}
