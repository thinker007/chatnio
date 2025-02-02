package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

var _ *sql.DB

func ConnectMySQL() *sql.DB {
	// connect to MySQL
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.db"),
	))
	if err != nil {
		log.Fatalln("Failed to connect to MySQL server: ", err)
	} else {
		log.Println("Connected to MySQL server successfully")
	}

	CreateUserTable(db)
	CreateConversationTable(db)
	CreateSubscriptionTable(db)
	CreatePackageTable(db)
	CreatePaymentLogTable(db)
	return db
}

func CreateUserTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS auth (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  bind_id INT UNIQUE,
		  username VARCHAR(24) UNIQUE,
		  token VARCHAR(255) NOT NULL,
		  password VARCHAR(64) NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreatePaymentLogTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS payment_log (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  amount DECIMAL(12,2) DEFAULT 0,
		  description VARCHAR(3600),
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateSubscriptionTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS subscription (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  plan_id INT,
		  expired_at DATETIME,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreatePackageTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS package (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  money DECIMAL(12,2) DEFAULT 0,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateConversationTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS conversation (
		  id INT PRIMARY KEY AUTO_INCREMENT,
		  user_id INT,
		  conversation_id INT,
		  conversation_name VARCHAR(255),
		  data TEXT,
		  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  UNIQUE KEY (user_id, conversation_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
