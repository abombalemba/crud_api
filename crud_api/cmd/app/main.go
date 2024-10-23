package main

import (
	"log"
	"net/http"
	"encoding/json"

	"CRUD_API/internal/models"
	pkg_logger "CRUD_API/pkg/logger"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

var (
	logger *log.Logger
	db *gorm.DB
)

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../../database/users.db"), &gorm.Config{})
	if err != nil {
		logger.Fatalln(err)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func createUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("handlerCreateUser successfully started")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Println("handlerCreateUser", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		logger.Println("handlerCreateUser", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)

	logger.Println("handlerCreateUser successfully executed")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	logger.Println("handlerGetUsers successfully started")

	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)

	logger.Println("handlerGetUsers successfully executed")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("handlerUpdateUser successfully started")

	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Model(&user).Where("id = ?", id).Updates(user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)

	logger.Println("handlerUpdateUser successfully executed")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("handlerDeleteUser successfully started")

	vars := mux.Vars(r)
	id := vars["id"]

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	logger.Println("handlerDeleteUser successfully executed")
}

func init() {
	pkg_logger.InitLogger()
	logger = pkg_logger.GetLogger()
}

func main() {
	defer pkg_logger.CloseLogger()

	setupDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")

	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	logger.Panicln(http.ListenAndServe(":5005", r))
}