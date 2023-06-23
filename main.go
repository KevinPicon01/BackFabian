package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"kevinPicon/go/rest-ws/Middleware"
	"kevinPicon/go/rest-ws/handlers"
	"kevinPicon/go/rest-ws/server"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")
	fmt.Println(PORT, JWT_SECRET, DATABASE_URL, "env variables")
	serv, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})
	if err != nil {
		log.Fatal(err)
	}
	serv.Start(BindRouters)

}
func BindRouters(s server.Server, r *mux.Router) {

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods("GET")
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods("POST")
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods("GET")
	r.HandleFunc("/users", handlers.UsersHandler(s)).Methods("GET")
	r.HandleFunc("/service", handlers.ServiceHandler(s)).Methods("POST")
	r.HandleFunc("/ecan", handlers.UpdateEcan(s)).Methods("GET")
	r.HandleFunc("/complaints", handlers.InsertComplaints(s)).Methods("POST")
	r.Use(Middleware.CheckAuthMiddleware(s))
}
