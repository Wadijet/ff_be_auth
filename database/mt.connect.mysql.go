package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// connectToMySQL kết nối đến cơ sở dữ liệu MySQL và trả về handle của cơ sở dữ liệu
func ConnectToMySQL(uri string) (*sql.DB, error) {

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	// Thiết lập các thông số của connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CloseInstance closes the MongoDB client connection.
func DisconnectFromMySql(client *sql.DB) error {

	err := client.Close()
	if err != nil {
		log.Printf("Failed to disconnect MongoDB client: %v", err)
		return err
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil
}
