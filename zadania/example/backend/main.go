package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func waitForDB(dsn string) {
	for {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Nie można połączyć się z bazą danych, ponawianie próby za 2 sekundy...")
			time.Sleep(2 * time.Second)
			continue
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Println("Nie można połączyć się z bazą danych, ponawianie próby za 2 sekundy...")
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Połączenie z bazą danych nawiązane.")
		break
	}
}

func main() {
	// Pobranie zmiennych środowiskowych
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Połączenie z bazą danych
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	waitForDB(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Nie można połączyć się z bazą danych: %v", err)
	}
	defer db.Close()

	// Utworzenie tabeli (jeśli nie istnieje)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INT AUTO_INCREMENT PRIMARY KEY, message TEXT NOT NULL)")
	if err != nil {
		log.Fatalf("Błąd przy tworzeniu tabeli: %v", err)
	}

	// Endpoint do zapisywania wiadomości
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Tylko metoda POST jest obsługiwana", http.StatusMethodNotAllowed)
			return
		}

		message := r.FormValue("message")
		if message == "" {
			http.Error(w, "Wiadomość nie może być pusta", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO messages (message) VALUES (?)", message)
		if err != nil {
			http.Error(w, "Błąd przy zapisywaniu wiadomości", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Wiadomość zapisana")
	})

	// Endpoint do wyświetlania wiadomości
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, message FROM messages")
		if err != nil {
			http.Error(w, "Błąd przy pobieraniu wiadomości", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var message string
			rows.Scan(&id, &message)
			fmt.Fprintf(w, "ID: %d, Message: %s\n", id, message)
		}
	})

	log.Println("Serwer działa na porcie 8080")
	http.ListenAndServe(":8080", nil)
}
