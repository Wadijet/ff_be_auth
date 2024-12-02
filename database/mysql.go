// database/mysql.go
package database

import (
	"atk-go-server/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// GetMySQLInstance khởi tạo kết nối MySQL
func GetMySQLInstance(cfg *config.Configuration) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.MySQLConnectionURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
