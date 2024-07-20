package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"        // замените на ваше имя пользователя
	password = "workout+5"       // замените на ваш пароль
	dbname   = "message_db_test" // замените на имя вашей базы данных
)

var db *sql.DB

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods("GET")
	http.Handle("/", r)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Добавление пользователя в командной строке windows:   curl -X POST --header "Content-Type: application/json" -d "{\"username\":\"amongus\",\"password\":\"42422\"}" http://localhost:8081/users
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var busy_id int
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Автоматически присваивает пользователям уникальный ID
	sql_statement := "SELECT id from users;"
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		switch err := rows.Scan(&busy_id); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			if user.ID == busy_id {
				for user.ID == busy_id {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					user.ID = r.Intn(10000) + 10

				}
			}
		default:
			checkError(err)
		}
	}
	// Выполнение запроса INSERT в базу данных для сохранения данных пользователя
	sql_statement = "INSERT INTO users (id, username, password) VALUES ($1, $2,$3);"
	_, err = db.Exec(sql_statement, user.ID, user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Успешно добавили пользователя
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	err := db.QueryRow("SELECT id, username,password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Успешно извлекли пользователя
	json.NewEncoder(w).Encode(user)
}
