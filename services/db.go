package services

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LoginData struct {
	PasswordHash string
	SessionToken string
	CSRFToken    string
}

var db *sql.DB

func InitDB() error {
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", "admin:password@tcp(db:3306)/auth_server_db")
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        username VARCHAR(255) PRIMARY KEY,
        password_hash TEXT,
        session_token TEXT,
        csrf_token TEXT
    )`)
	if err != nil {
		return err
	}

	log.Println("Database connection established and table ensured.")
	return nil
}

func LoadUserData(username string) (LoginData, error) {
	var data LoginData
	err := db.QueryRow("SELECT password_hash, session_token, csrf_token FROM users WHERE username = ?", username).Scan(&data.PasswordHash, &data.SessionToken, &data.CSRFToken)
	if err != nil {
		return data, err
	}

	log.Println("User data loaded from database for user:", username)
	return data, nil
}

func SaveUserData(username string, data LoginData) error {
	_, err := db.Exec("REPLACE INTO users (username, password_hash, session_token, csrf_token) VALUES (?, ?, ?, ?)", username, data.PasswordHash, data.SessionToken, data.CSRFToken)
	if err != nil {
		return err
	}

	log.Println("User data saved to database for user:", username)
	return nil
}
