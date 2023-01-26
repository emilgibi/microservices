package main

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	ID            int    `json:"id"`
	OrderName     string `json:"order_name"`
	OrderQuantity string `json:"order_quantity"`
}

var db *gorm.DB

func main() {

	var err error
	dsn := "host=host.docker.internal user=postgres password=Emilgibi@123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbinstance, _ := db.DB()
	defer dbinstance.Close()

	router := mux.NewRouter()

	db.AutoMigrate(Order{})

	router.HandleFunc("/order", getOrder).Methods("GET")
	router.HandleFunc("/order", addOrder).Methods("POST")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8083", router)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	db.Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	var addOrder Order
	_ = json.NewDecoder(r.Body).Decode(&addOrder)
	db.Create(&addOrder)
	json.NewEncoder(w).Encode(addOrder)
}
