package storage

import (
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type Page struct {
	ID    string    `json:"id,omitempty"` //url
	Text  string    `json:"page"`
	Links []string  `json:"links"`
	Date  time.Time `json:"date"`
}

type URL struct {
	ID string `json:"id,omitempty"`
	LastCheck time.Time `json:"lastCheck"`
	Names []string `json:"names"`
	Pages []string `json:"pages"`
}

func Init(addr, user, pass, namespace, database string) (db *surrealdb.DB) {
	db, err := surrealdb.New(addr)
	if err != nil {
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": user,
		"pass": pass,
	}); err != nil {
		panic(err)
	}

	if _, err = db.Use(namespace, database); err != nil {
		panic(err)
	}

	return db
}

func AddPage(db *surrealdb.DB, page Page) {

	// Insert page
	data, err := db.Create("pages", page)
	if err != nil {
		panic(err)
	}

	//Make mentions

	// Unmarshal data
	addedPage := make([]Page, 1)
	err = surrealdb.Unmarshal(data, &addedPage)
	if err != nil {
		panic(err)
	}

	//fmt.Println(addedPage[0].ID, addedPage[0].Date.String())


}
