package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/emilgibi/inventory-microservices/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (handler *Handler) Connect(host, user, pass, dbName, port string) {
	var err error
	dsn := "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
	handler.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	handler.DB.AutoMigrate(models.Stock{})
	if err != nil {
		panic(err)
	}
}

func (handler *Handler) CheckStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var stock models.Stock
	handler.DB.First(&stock, params["id"])
	if stock.ID == 0 {
		http.Error(w, "Stock not found", http.StatusBadRequest)
		return
	}

	orders := "http://localhost:8084/order/" + params["id"]
	orderResp, err := http.Get(orders)
	if err != nil {
		http.Error(w, "Error getting order information", http.StatusInternalServerError)
		return
	}
	defer orderResp.Body.Close()

	var OrderQuantity int
	err = json.NewDecoder(orderResp.Body).Decode(&OrderQuantity)
	if err != nil {
		http.Error(w, "Error parsing order information", http.StatusInternalServerError)
		return
	}

	var status int
	var message bool
	if stock.ProductQuantity > OrderQuantity {
		status = 1
		message = true
	} else {
		status = 0
		message = false
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status  int  `json:"status"`
		Message bool `json:"message"`
	}{
		Status:  status,
		Message: message,
	})
}
func (handler *Handler) AddStock(w http.ResponseWriter, r *http.Request) {
	var addStock models.Stock
	_ = json.NewDecoder(r.Body).Decode(&addStock)
	handler.DB.Create(&addStock)
	json.NewEncoder(w).Encode(addStock)
}

func (handler *Handler) DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var stocks models.Stock
	handler.DB.First(&stocks, params["id"])
	handler.DB.Delete(&stocks)
	json.NewEncoder(w).Encode("Person successfully deleted")
}
