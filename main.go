package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/om00/userCourse/handler"
	"github.com/om00/userCourse/mysqldb"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	mysqldb.Dbpath = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening MySQL database connection: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging MySQL database: ", err)
	}

	fmt.Println("Successfully connected to the database")
}
func main() {

	migrateCmd := flag.String("m", "", "Run migration command (up, down, force N, drop)")
	seederName := flag.String("s", "", "Name of the seeder to run")
	flag.Parse()

	if *migrateCmd != "" {
		mysqldb.RunMigrations(*migrateCmd)
		return
	}

	if *seederName != "" {
		mysqldb.CallSeederFunction(db, *seederName)
	}

	app := handler.App{Db: mysqldb.NewDB(db), Validator: validator.New()}

	// Route handlers
	http.HandleFunc("/create_user", app.CreateUser)
	http.HandleFunc("/getAllcourse", app.GetAllCourse)
	http.HandleFunc("/user-course", app.UserCourse)
	http.HandleFunc("/user-details", app.UserDetials)

	fmt.Println("Server started on :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
