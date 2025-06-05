package config
import (
	"os"
	"fmt"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)
var DB *goqu.Database
func init() {
    err := godotenv.Load()
    if err != nil {
		fmt.Println("Lỗi khi load file .env. Lỗi:", err)
    }
}
func InitDB() {
	
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
    	user, pass, host, port, name,
	)
    sqlDB, err := sql.Open("mysql", dsn)
    if err != nil {
		fmt.Println("Err:", err)
		return
    }
	//defer sqlDB.Close()
	err = sqlDB.Ping()
	if err != nil {
		fmt.Printf("Không thể kết nối đến MySQL: %v", err)
	}
	
    DB = goqu.New("mysql", sqlDB)
	fmt.Println("Connected to database.")
	
}