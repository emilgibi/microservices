package main

import (
	"encoding/json"

	"net/http"

	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"order_name"`
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

	db.AutoMigrate(User{})

	router.HandleFunc("/user", getUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", getUserid).Methods("GET")
	router.HandleFunc("/user", updateUser).Methods("POST")
	router.HandleFunc("/products", getProduct).Methods("GET")
	router.HandleFunc("/product/{product_id}", getProductById).Methods("GET")
	router.HandleFunc("/order", getOrder).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8083", router)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func getUserid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var users User
	db.First(&users, params["id"])
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var update User
	_ = json.NewDecoder(r.Body).Decode(&update)

	db.Save(&update)
	json.NewEncoder(w).Encode(update)
}

func getProduct(w http.ResponseWriter, r *http.Request) {

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8081/products")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getProductById(w http.ResponseWriter, r *http.Request) {

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8081/products{product_id}")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getOrder(w http.ResponseWriter, r *http.Request) {

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8083/order")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	client := resty.New()

	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://localhost:8083/orders")

	var productAvailability bool = true
	if productAvailability {
		_, _ = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(newUser).
			Post("https://localhost:8082/inventory")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User added successfully.")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("One or more products are not available.")
	}
}
