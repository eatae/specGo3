package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	Id     string `json:"Id"`
	Title  string `json:"Title"`
	Amount string `json:"Amount"`
	Price  string `json:"Price"`
}

type ErrorMessage struct {
	Message string `json:"Message"`
}

//products - local DataBase
var Products []Product

//GET request for /products
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hint: getAllProducts woked.....")
	json.NewEncoder(w).Encode(Products) //ResponseWriter - место , куда пишем. Products - кого пишем
}

//GET request for product with ID
func GetProductWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	find := false
	for _, product := range Products {
		if product.Id == vars["id"] {
			find = true
			json.NewEncoder(w).Encode(product)
		}
	}
	if !find {
		var erM = ErrorMessage{Message: "Not found product with that ID"}
		json.NewEncoder(w).Encode(erM)
	}
}

//PostNewProduct func for create new Product
func PostNewProduct(w http.ResponseWriter, r *http.Request) {
	// {
	// 	"Id" : "3",
	// 	"Title" : "Title from json POST method",
	// 	"Amount" : "150",
	// 	"Price" : "12.4"
	// }
	reqBody, _ := ioutil.ReadAll(r.Body)
	var product Product
	json.Unmarshal(reqBody, &product) // Считываем все из тела зпроса в подготовленный пустой объект Product

	Products = append(Products, product)
	json.NewEncoder(w).Encode(product) //После добавления новой статьи возвращает добавленную
}

// Delete
func DeleteProductWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, product := range Products {
		if product.Id == id {
			Products = append(Products[:index], Products[index+1:]...)
		}
	}
}

// PutExistsProduct ....
func PutExistsProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idKey := vars["id"] // СТРОКА
	finded := false

	for index, product := range Products {
		if product.Id == idKey {
			finded = true
			reqBody, _ := ioutil.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)        // Изменяем статус код на 202
			json.Unmarshal(reqBody, &Products[index]) // перезаписываем всю информацию для статьи с Id
		}
	}

	if !finded {
		w.WriteHeader(http.StatusNotFound) // Изменяем статус код на 404
		var erM = ErrorMessage{Message: "Not found product with that ID. Try use POST first"}
		json.NewEncoder(w).Encode(erM)
	}

}

func main() {
	//Добавляю 2 статьи в свою базу
	Products = []Product{
		Product{Id: "1", Title: "First title", Amount: "10", Price: "140.2"},
		Product{Id: "2", Title: "Second title", Amount: "120", Price: "14.2"},
	}
	fmt.Println("REST API V2.0 worked....")
	//СОздаю свой маршрутизатор на основе либы mux
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/products", GetAllProducts).Methods("GET")
	myRouter.HandleFunc("/product/{id}", GetProductWithId).Methods("GET")
	//Создадим запрос на добавление новой статьи
	myRouter.HandleFunc("/product", PostNewProduct).Methods("POST")

	//Создадим запрос на удаление статьи (гарантировано существует)
	myRouter.HandleFunc("/product/{id}", DeleteProductWithId).Methods("DELETE")

	myRouter.HandleFunc("/product/{id}", PutExistsProduct).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
