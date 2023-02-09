package main

import (
	"net/http"
	"os"

	"github.com/emilgibi/inventory-microservices/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	var handlerObj handlers.Handler

	godotenv.Load()

	HOST := os.Getenv("DB_HOST")
	USER := os.Getenv("USER_NAME")
	PASS := os.Getenv("PASS")

	handlerObj.Connect(HOST, USER, PASS, "postgres", "5432")

	dbinstance, _ := handlerObj.DB.DB()
	defer dbinstance.Close()

	router := mux.NewRouter()

	router.HandleFunc("/stock/check", handlerObj.CheckStock).Methods("GET")
	router.HandleFunc("/stock/add", handlerObj.AddStock).Methods("POST")
	router.HandleFunc("/stock/remove", handlerObj.DeleteStock).Methods("DELETE")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8082", router)
}
