package db

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"

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
func (h *Database) CheckPassword(username, password string) (bool) {
	//get hashed password for compare
	var passwordHash string
	h.db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&passwordHash)
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
	user.Password = utils.HashPassword(user.Password)
	_, err := h.db.Exec("INSERT INTO users (id, name, password, role) VALUES (UUID(), ?, ?, ?)", user.Username, user.Password, user.Role)
	return err
}


func (h *Database) AddTicket(ticket *types.CreateTicket) (bool, error) {
	// insert new ticket into database
	_, err := h.db.Exec("INSERT INTO tickets (id, issuer, date, title, status, content) VALUES (UUID(), ?, ?, ?, 0, ?)",
											 ticket.Issuer, time.Now().Local().Unix(), ticket.Title, ticket.Content)
	return true, err
}

func (h *Database) GetUserId(username string) (userId string) {
	h.db.QueryRow("SELECT id FROM users WHERE name = ?", username).
		Scan(&userId)
	return
}

func (h *Database) GetUser(userId string) (user types.User) {
	h.db.QueryRow("SELECT id, name, role FROM users WHERE id = ?", userId).
		Scan(&user.Id, &user.Name, &user.Role)
	return
}

func (h *Database) SaveCookie(userId, cookie string, expires *time.Time) error {
	_, err := h.db.Exec("INSERT INTO sessions (userid, cookie, expires) VALUES (?, ?, ?)",
									userId, cookie, expires.Unix())	
	return err
}

func (h *Database) UpdateCookie(userId, cookie string, expires *time.Time) error {
	_, err := h.db.Exec("UPDATE sessions SET cookie = ?, expires = ? WHERE userid = ?",
									cookie, expires.Unix(), userId)	
	return err
}

func (h *Database) SessionExist(cookie string) (bool){
	var expires int64
	h.db.QueryRow("SELECT expires FROM sessions WHERE cookie = ?", cookie).Scan(&expires)
	return expires != 0
}

func (h *Database) UserHasSession(userId string) (bool){
	var expires int64
	h.db.QueryRow("SELECT expires FROM sessions WHERE userid = ?", userId).Scan(&expires)
	return expires != 0
}


func (h *Database) CookieUserId(cookie string) (userId string){
	h.db.QueryRow("SELECT userid FROM sessions WHERE cookie = ?", cookie).Scan(&userId)
	return 
}

func (h *Database) CookieAuthorized(cookie string) (bool){
	if !h.SessionExist(cookie){
		return false
	}
	var user types.User
	h.db.QueryRow("SELECT userid FROM sessions WHERE cookie = ?", cookie).Scan(&user.Id)
	user = h.GetUser(user.Id)
	return user.Role == "admin"
}


func (h *Database) DeleteCookie(cookie string) error {
	_, err := h.db.Exec("DELETE FROM sessions WHERE cookie = ?", cookie)	
 	return err
}

func (h *Database) DeleteTicket(ticketId string) error {
	_, err := h.db.Exec("DELETE FROM tickets WHERE id = ?", ticketId)	
 	return err
}

func (h *Database) GetAllTickets(users []string) (tickets []*types.Ticket, err error){
	var rows *sql.Rows
	//user not spepcified - get all tickets
	if len(users) == 0 {
		rows, err = h.db.Query("SELECT * FROM tickets")
		if err != nil{
			return
		}
		for rows.Next() {
			ticket := &types.Ticket{}
			var date   int64
			var issuer string
			rows.Scan(&ticket.Id, &issuer, &date, &ticket.Title, &ticket.Status, &ticket.Content)
			ticket.Issuer = h.GetUser(issuer)
			ticket.Date = time.Unix(date, 0)
			tickets = append(tickets, ticket)
		}
	//iterate over given users
	}else{
		for _,user :=  range users{
			user = h.GetUserId(user)
			if user == ""{
				continue
			}
			rows, err = h.db.Query("SELECT id, title, status, content, issuer, date FROM tickets WHERE issuer = ?",user)
			for rows.Next() {
				ticket := &types.Ticket{}
				var date   int64
				var issuer string
				rows.Scan(&ticket.Id, &ticket.Title, &ticket.Status, &ticket.Content, &issuer, &date)
				ticket.Issuer = h.GetUser(issuer)
				ticket.Date = time.Unix(date, 0)
				tickets = append(tickets, ticket)
			}
		}
	}
	return
}

func (h *Database) ChangePassword(username, password string) error {
	_, err := h.db.Exec("UPDATE users SET password = ? WHERE name = ?",
								utils.HashPassword(password), username)	
	return err
}

func (h *Database) GetTicketIssuer(ticketId string) (userId string, err error){
	err = h.db.QueryRow("SELECT issuer FROM tickets WHERE id = ?", ticketId).Scan(&userId)
	return
}