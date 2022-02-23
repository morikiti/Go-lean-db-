package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Opening struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("mysql", os.Getenv("DB_ROLE")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	parms := mux.Vars(r)
	var name Opening
	err = db.QueryRow("SELECT * FROM sample WHERE id = ?", parms["id"]).Scan(&name.ID, &name.Name)

	if err != nil {
		panic(err.Error())
	}
	//NewEncode　でJSONに変換
	json.NewEncoder(w).Encode(name)
	fmt.Println(name)
}

func getNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", os.Getenv("DB_ROLE")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//rowを取得
	rows, err := db.Query("SELECT * FROM sample")
	if err != nil {
		panic(err.Error())
	}
	//	json.NewEncoder(w).Encode(rows)
	fmt.Println(rows)
	openingArgs := make([]Opening, 0)
	for rows.Next() {
		var opening Opening
		err = rows.Scan(&opening.ID, &opening.Name)
		if err != nil {
			panic(err.Error())
		}
		openingArgs = append(openingArgs, opening)
	}
	json.NewEncoder(w).Encode(openingArgs)
	fmt.Println(openingArgs)

}

func addName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//var name Opening
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var book Opening
	_ = json.NewDecoder(r.Body).Decode(&book)
	json.NewDecoder(r.Body)

	fmt.Println(book)
	fmt.Println(book.ID)
	fmt.Println(book.Name)
	//books = append(books, book)
	db, err := sql.Open("mysql", os.Getenv("DB_ROLE")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO sample(id,name) VALUES(?,?)", book.ID, book.Name)
	if err != nil {
		panic(err.Error())
	}

	return
}

func main() {

	//ルーターの初期化
	router := mux.NewRouter()

	page := "main.go "
	fmt.Println(page + "OK???")
	//ルーティング（エンドポイント）
	router.HandleFunc("/api/sampls/{id}", getName).Methods("GET")
	router.HandleFunc("/api/samples", getNames).Methods("GET")
	router.HandleFunc("/api/samples", addName).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", router))
}
