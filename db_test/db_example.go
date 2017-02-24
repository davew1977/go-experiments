package main

import "database/sql"

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	db, _ := sql.Open("mysql", "dev:dev@tcp(localhost:8889)/go_wiki")
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, _ := db.Query("select * from Page")
	for rows.Next() {
		var page Page
		if err := rows.Scan(&page.Title, &page.Body); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", page)
	}
}
