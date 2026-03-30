package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"april-fools-day/internal/db"
)

type PageData struct {
	VisitorNumber int
	IsNewVisitor  bool
}

func newUUID() string {
	var b [16]byte
	rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func runMigrations(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS visitors (
			id             UUID      PRIMARY KEY,
			visitor_number INTEGER   NOT NULL,
			created_at     TIMESTAMP NOT NULL DEFAULT now()
		);
		CREATE SEQUENCE IF NOT EXISTS visitor_number_seq START 1;
	`)
	return err
}

func getOrCreateVisitor(database *sql.DB, w http.ResponseWriter, r *http.Request) (int, bool, error) {
	cookie, err := r.Cookie("visitor_id")
	if err == nil {
		var num int
		err = database.QueryRow(`SELECT visitor_number FROM visitors WHERE id = $1`, cookie.Value).Scan(&num)
		if err == nil {
			return num, false, nil
		}
		if err != sql.ErrNoRows {
			return 0, false, err
		}
	}

	uuid := newUUID()
	var num int
	err = database.QueryRow(
		`INSERT INTO visitors (id, visitor_number) VALUES ($1, nextval('visitor_number_seq')) RETURNING visitor_number`,
		uuid,
	).Scan(&num)
	if err != nil {
		return 0, false, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "visitor_id",
		Value:    uuid,
		MaxAge:   365 * 24 * 3600,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return num, true, nil
}

func homeHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		num, isNew, err := getOrCreateVisitor(database, w, r)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, PageData{VisitorNumber: num, IsNewVisitor: isNew})
	}
}

func main() {
	database, err := db.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := runMigrations(database); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler(database))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
