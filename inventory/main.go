package main

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Stock struct {
	ID              int    `json:"id"`
	ProductName     string `json:"product_name"`
	ProductQuantity int    `json:"product_quantity"`
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

	db.AutoMigrate(Stock{})

	router.HandleFunc("/stock/check", checkStock).Methods("GET")
	router.HandleFunc("/stock/add", addStock).Methods("POST")
	router.HandleFunc("/stock/remove", deleteStock).Methods("DELETE")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8082", router)
}

func checkStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var stock Stock
	db.First(&stock, params["id"])
	if stock.ID == 0 {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Product  string `json:"product"`
		Quantity int    `json:"quantity"`
	}{
		Product:  stock.ProductName,
		Quantity: stock.ProductQuantity,
	})
}

func addStock(w http.ResponseWriter, r *http.Request) {
	var addStock Stock
	_ = json.NewDecoder(r.Body).Decode(&addStock)
	db.Create(&addStock)
	json.NewEncoder(w).Encode(addStock)
}

func deleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var stocks Stock
	db.First(&stocks, params["id"])
	db.Delete(&stocks)
	json.NewEncoder(w).Encode("Person successfully deleted")
}
