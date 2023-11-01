package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	port         = ":3000"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName)

	mux := mux.NewRouter()
	mux.HandleFunc("/home", index(dbx))
	mux.HandleFunc("/post/{postID}", post(dbx))
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/admin", admin(dbx))

	mux.HandleFunc("/api/post", createPost(dbx)).Methods(http.MethodPost)
	//mux.HandleFunc("/api/login", searchUser(dbx)).Methods(http.MethodPost)
	//mux.HandleFunc("/api/logout", deleteUser(dbx)).Methods(http.MethodPost)
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Start server " + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:HEJIOMANTErADbI2021@tcp(localhost:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
