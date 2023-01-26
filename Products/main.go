package main

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID              int    `json:"id"`
	ProductName     string `json:"product_name"`
	ProductQuantity string `json:"product_quantity"`
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

	db.AutoMigrate(Product{})

	router.HandleFunc("/products", getProduct).Methods("GET")
	router.HandleFunc("/product/{product_id}", getuser).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8081", router)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	var products []Product
	db.Find(&products)
	json.NewEncoder(w).Encode(products)
}

func getuser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var products Product
	db.First(&products, params["id"])
	json.NewEncoder(w).Encode(products)
}
