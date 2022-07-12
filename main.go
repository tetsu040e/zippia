package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const appName = "zippia"
const version = "0.0.1"

func main() {
	addr := flag.String("a", "127.0.0.1:8000", "address to bind.")
	path := flag.String("p", "/", "API endpoint path")
	showVersion := flag.Bool("v", false, "show version.")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %s\n", appName, version)
		return
	}

	db, err := initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc(*path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != *path {
			http.NotFound(w, r)
			return
		}

		zip := r.URL.Query().Get("zip")

		stmt, err := db.Prepare("SELECT zip, pref, city, town, pref_kana, city_kana, town_kana FROM address WHERE zip = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		rows, err := stmt.Query(zip)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		list := []Address{}
		for rows.Next() {
			var address Address
			err = rows.Scan(&address.Zip, &address.Pref, &address.City, &address.Town, &address.PrefKana, &address.CityKana, &address.TownKana)
			if err != nil {
				break
			}
			list = append(list, address)
		}
		if len(list) == 0 {
			http.NotFound(w, r)
			return
		}

		body, err := json.Marshal(list)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	})

	log.Printf("http server started on http://%s%s\n", *addr, *path)
	http.ListenAndServe(*addr, nil)
}

//go:embed kenall.json
var kenall []byte

type Address struct {
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
	log.Println("initializing DB...")
	db, err := sql.Open("sqlite3", "file:zip.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}

	initialTableSQL := `
DROP INDEX IF EXISTS zipcode;
DROP TABLE IF EXISTS address;
CREATE TABLE IF NOT EXISTS address (
    id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    zip       TEXT,
    pref      TEXT,
    city      TEXT,
    town      TEXT,
    pref_kana TEXT,
    city_kana TEXT,
    town_kana TEXT
);
CREATE INDEX zipcode ON address (zip);
    `
	_, err = db.Exec(initialTableSQL)
	if err != nil {
		return nil, err
	}

	var list []Address
	if err := json.Unmarshal(kenall, &list); err != nil {
		return nil, err
	}

	const bufSize = 100
	buf := []Address{}
	for _, address := range list {
		buf = append(buf, address)
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
	log.Println("DB initialization is completed!")
	return db, nil
}

func flush(db *sql.DB, buf *[]Address) error {
	insertSQL, args := createBulkInsertQuery(*buf)
	*buf = []Address{}

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

func createBulkInsertQuery(list []Address) (string, []interface{}) {
	placeholder := []string{}
	args := []interface{}{}

	for _, address := range list {
		placeholder = append(placeholder, "(?,?,?,?,?,?,?)")
		args = append(args, address.Zip, address.Pref, address.City, address.Town, address.PrefKana, address.CityKana, address.TownKana)
	}

	query := fmt.Sprintf("INSERT INTO address (zip, pref, city, town, pref_kana, city_kana, town_kana) VALUES %s", strings.Join(placeholder, ","))

	return query, args
}
