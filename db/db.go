package db

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"

	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"


	
)

type Database struct {
	db *sql.DB
}

func ConnectDB() (*Database, error) {
	//connect to the database
	db, _ := func() (*sql.DB, error) {
		dbUser := "root"
		dbPass := "password"
		dbName := "ticketing_system"
		return sql.Open("mysql", dbUser+":"+dbPass+"@(127.0.0.1:3306)/"+dbName)
	}()

	//create db stuct
	DbHandler := &Database{
		db: db,
	}
		
	//return error if can't ping database
	err := DbHandler.db.Ping()

	return DbHandler, err
}

func (h *Database) Close() {
	h.db.Close()
}

// validate password for login
func (h *Database) CheckPassword(username, password string) bool {
	//get hashed password for compare
	var passwordHash string
	h.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&passwordHash)

	//return if passwords maches
	return utils.Comparepassword(passwordHash, password)
}

//check if username already exist
func (h *Database) CheckUserExist(username string) bool {
	var exist string
	h.db.QueryRow("SELECT name FROM users WHERE name = ?", username).Scan(&exist)
	return exist != ""
}


//add user: [name, password, role] into database
func (h *Database) AddUser(user *types.Register) error {
	fmt.Println(user)
	_, err := h.db.Exec("INSERT INTO users (id, name, password, role) VALUES (uuid(), ?, ?, ?)", user.Username, user.Password, user.Role)
	return err
}
