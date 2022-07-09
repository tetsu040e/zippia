package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		zip := r.URL.Query().Get("zip")

		stmt, err := db.Prepare("SELECT zip, pref, city, town, pref_kana, city_kana, town_kana FROM row WHERE zip = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		rows, err := stmt.Query(zip)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		list := []Row{}
		for rows.Next() {
			var row Row
			err = rows.Scan(&row.Zip, &row.Pref, &row.City, &row.Town, &row.PrefKana, &row.CityKana, &row.TownKana)
			if err != nil {
				break
			}
			list = append(list, row)
		}

		body, err := json.Marshal(list)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	})

	fmt.Println("http server started on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

//go:embed kenall.json
var kenall []byte

type Row struct {
	Id       int    `json:"-"`
	Zip      string `json:"zip"`
	Pref     string `json:"pref"`
	City     string `json:"city"`
	Town     string `json:"town"`
	PrefKana string `json:"pref_kana"`
	CityKana string `json:"city_kana"`
	TownKana string `json:"town_kana"`
}

func initialize() (*sql.DB, error) {
	fmt.Println("initializing...")
	db, err := sql.Open("sqlite3", "file:zip.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}

	initialTableSQL := `
    DROP INDEX IF EXISTS zipcode;
    DROP TABLE IF EXISTS smaple;
    CREATE TABLE IF NOT EXISTS row (
        id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        zip      TEXT,
        pref      TEXT,
        city      TEXT,
        town      TEXT,
        pref_kana TEXT,
        city_kana TEXT,
        town_kana TEXT
    );
    CREATE INDEX zipcode ON row (zip);
    `
	_, err = db.Exec(initialTableSQL)
	if err != nil {
		return nil, err
	}

	var zipList []Row
	if err := json.Unmarshal(kenall, &zipList); err != nil {
		return nil, err
	}

	buf := []Row{}
	bufSize := 50
	for _, row := range zipList {
		buf = append(buf, row)
		if len(buf) >= bufSize {
			err = flush(db, &buf)
			if err != nil {
				return nil, err
			}
		}
	}
	if len(buf) > 0 {
		err = flush(db, &buf)
		if err != nil {
			return nil, err
		}
	}

	db.SetConnMaxLifetime(-1)
	fmt.Println("initialized is completed!!")
	return db, nil
}

func flush(db *sql.DB, buf *[]Row) error {
	insertSQL, args := createBulkInsertQuery(*buf)
	*buf = []Row{}

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func createBulkInsertQuery(list []Row) (string, []interface{}) {
	placeholder := []string{}
	args := []interface{}{}

	for _, row := range list {
		placeholder = append(placeholder, "(?,?,?,?,?,?,?)")
		args = append(args, row.Zip, row.Pref, row.City, row.Town, row.PrefKana, row.CityKana, row.TownKana)
	}

	query := fmt.Sprintf("INSERT INTO row (zip, pref, city, town, pref_kana, city_kana, town_kana) VALUES %s", strings.Join(placeholder, ","))

	return query, args
}
