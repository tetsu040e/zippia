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

var Version = "v0.3.266"

//go:embed var/last-modified.txt
var lastModified []byte

//go:embed var/banner.txt
var banner []byte

//go:embed var/kenall.json
var kenall []byte

//go:embed var/jigyosyo.json
var jigyosyo []byte

func main() {
	host := flag.String("host", "127.0.0.1", "host part of bind address.")
	port := flag.String("port", "5000", "port part of bind address.")
	showVersion := flag.Bool("v", false, "show version.")
	showVersionFull := flag.Bool("vv", false, "show version and last modified date.")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %s\n", appName, strings.TrimSuffix(Version, "\n"))
		return
	}
	if *showVersionFull {
		fmt.Printf(
			"%s %s\nZip code was last modified on \"%s\"\n",
			appName,
			strings.TrimSuffix(Version, "\n"),
			strings.TrimSuffix(string(lastModified), "\n"),
		)
		return
	}

	db, err := initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		zip := r.URL.Query().Get("zip")
		zip = strings.Replace(zip, "-", "", 1)

		stmt, err := db.Prepare("SELECT zip, pref, city, town, pref_kana, city_kana, town_kana, office, office_kana FROM address WHERE zip = ?")
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
			var row Address
			err = rows.Scan(&row.Zip, &row.Pref, &row.City, &row.Town, &row.PrefKana, &row.CityKana, &row.TownKana, &row.Office, &row.OfficeKana)
			if err != nil {
				log.Println(err)
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

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Println(fmt.Sprintf(`start the Japanese zip code search API server.
%s
version %s
Zip code was last modified on "%s"
`, string(banner), string(Version), string(lastModified)))
	log.Printf("http server started on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type Address struct {
	Id         int    `json:"-"`
	Zip        string `json:"zip"`
	Pref       string `json:"pref"`
	City       string `json:"city"`
	Town       string `json:"town"`
	Office     string `json:"office"`
	PrefKana   string `json:"pref_kana"`
	CityKana   string `json:"city_kana"`
	TownKana   string `json:"town_kana"`
	OfficeKana string `json:"office_kana"`
}

func initialize() (*sql.DB, error) {
	log.Println("initializing database...")
	db, err := sql.Open("sqlite3", "file:zip.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}

	initialTableSQL := `
DROP INDEX IF EXISTS zipcode;
DROP TABLE IF EXISTS address;
CREATE TABLE IF NOT EXISTS address (
    id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    zip         TEXT NOT NULL,
    pref        TEXT,
    city        TEXT,
    town        TEXT,
    office      TEXT,
    pref_kana   TEXT,
    city_kana   TEXT,
    town_kana   TEXT,
    office_kana TEXT
);
CREATE INDEX zipcode ON address (zip);
    `
	_, err = db.Exec(initialTableSQL)
	if err != nil {
		return nil, err
	}

	for _, bytes := range [][]byte{kenall, jigyosyo} {
		var list []Address
		if err := json.Unmarshal(bytes, &list); err != nil {
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
	}

	db.SetConnMaxLifetime(-1)
	log.Println("database initialization is completed!")
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

	for _, row := range list {
		placeholder = append(placeholder, "(?,?,?,?,?,?,?,?,?)")
		args = append(args, row.Zip, row.Pref, row.City, row.Town, row.PrefKana, row.CityKana, row.TownKana, row.Office, row.OfficeKana)
	}

	query := fmt.Sprintf("INSERT INTO address (zip, pref, city, town, pref_kana, city_kana, town_kana, office, office_kana) VALUES %s", strings.Join(placeholder, ","))

	return query, args
}
