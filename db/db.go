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

//add user: [name, password, role] into database
func (h *Database) AddUser(user *types.Register) error {
	user.Password = utils.HashPassword(user.Password)
	_, err := h.db.Exec("INSERT INTO users (id, name, password, role) VALUES (UUID(), ?, ?, ?)", user.Username, user.Password, user.Role)
	return err
}

//change password of a given user
func (h *Database) ChangePassword(username, password string) error {
	_, err := h.db.Exec("UPDATE users SET password = ? WHERE name = ?",
								utils.HashPassword(password), username)	
	return err
}

//create a ticket
func (h *Database) AddTicket(ticket *types.CreateTicket) (bool, error) {
	// insert new ticket into database
	_, err := h.db.Exec("INSERT INTO tickets (id, issuer, date, title, status, content) VALUES (UUID(), ?, ?, ?, 0, ?)",
											 ticket.Issuer, time.Now().Local().Unix(), ticket.Title, ticket.Content)
	return true, err
}

// delete a ticket
func (h *Database) DeleteTicket(ticketId string) error {
	_, err := h.db.Exec("DELETE FROM tickets WHERE id = ?", ticketId)	
 	return err
}

//save session cookie
func (h *Database) SaveCookie(userId, cookie string, expires *time.Time) error {
	_, err := h.db.Exec("INSERT INTO sessions (userid, cookie, expires) VALUES (?, ?, ?)",
									userId, cookie, expires.Unix())	
	return err
}

//update session cookie expiration date
func (h *Database) UpdateCookie(userId, cookie string, expires *time.Time) error {
	_, err := h.db.Exec("UPDATE sessions SET cookie = ?, expires = ? WHERE userid = ?",
									cookie, expires.Unix(), userId)	
	return err
}

//delete session cookie
func (h *Database) DeleteCookie(cookie string) error {
	_, err := h.db.Exec("DELETE FROM sessions WHERE cookie = ?", cookie)	
 	return err
}

// Get user by userId
func (h *Database) GetUser(userId string) (user types.User, err error) {
	err = h.db.QueryRow("SELECT * FROM users WHERE id = ?", userId).
		Scan(&user.Id, &user.Name,&user.Password, &user.Role)
	return 
}

// get userId
func (h *Database) GetUserId(username string) (userId string) {
	h.db.QueryRow("SELECT id FROM users WHERE name = ?", username).
		Scan(&userId)
	return
}

// get session cookie
func (h *Database) GetSessionCookie(sessionCookie string) (cookie types.SessionCookie, err error){
	err = h.db.QueryRow("SELECT * FROM sessions WHERE cookie = ?", sessionCookie).
			Scan(&cookie.UserId, &cookie.Cookie, &cookie.Expires)
	return 
}

func (h *Database) UserHasSession(userId string) (bool){
	var expires int64
	h.db.QueryRow("SELECT expires FROM sessions WHERE userid = ?", userId).Scan(&expires)
	return expires != 0
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
			//ticket.Issuer = h.GetUser(issuer)
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
				//ticket.Issuer = h.GetUser(issuer)
				ticket.Date = time.Unix(date, 0)
				tickets = append(tickets, ticket)
			}
		}
	}
	return
}

func (h *Database) GetTicketIssuer(ticketId string) (userId string, err error){
	err = h.db.QueryRow("SELECT issuer FROM tickets WHERE id = ?", ticketId).Scan(&userId)
	return
}