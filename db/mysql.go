package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/KameeKaze/Ticketing-system/types"
	"github.com/KameeKaze/Ticketing-system/utils"
)

type Database struct {
	db *sql.DB
}

// defiine database connection for mysql
var (
	Mysql = Database{
		db: func() *sql.DB {
			dbUser := "root"
			dbPass := os.Getenv("DATABASE_PASSWORD")
			dbName := "ticketing_system"
			dbHost := "database:3306"
			database, _ := sql.Open("mysql", dbUser+":"+dbPass+"@("+dbHost+")/"+dbName+"?parseTime=true")

			database.SetMaxOpenConns(10)
			database.SetMaxIdleConns(10)
			database.SetConnMaxLifetime(5 * time.Minute)
			return database
		}(),
	}
)

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
		ticket.Issuer, time.Now().Format(time.RFC3339), ticket.Title, ticket.Content)
	return true, err
}

// delete a ticket
func (h *Database) DeleteTicket(ticketId string) error {
	_, err := h.db.Exec("DELETE FROM tickets WHERE id = ?", ticketId)
	return err
}

// Get user by userId
func (h *Database) GetUser(userId string) (user types.User, err error) {
	err = h.db.QueryRow("SELECT * FROM users WHERE id = ?", userId).
		Scan(&user.Id, &user.Name, &user.Password, &user.Role)
	return
}

// get userId
func (h *Database) GetUserId(username string) (userId string) {
	h.db.QueryRow("SELECT id FROM users WHERE name = ?", username).
		Scan(&userId)
	return
}

func (h *Database) GetAllTickets(users []string) (tickets []*types.Ticket, err error) {
	var rows *sql.Rows
	//user not spepcified - get all tickets
	if len(users) == 0 {
		rows, err = h.db.Query("SELECT * FROM tickets")
		if err != nil {
			return
		}
		for rows.Next() {
			ticket := &types.Ticket{}
			var issuer string
			rows.Scan(&ticket.Id, &issuer, &ticket.Date, &ticket.Title, &ticket.Status, &ticket.Content)
			ticket.Issuer, _ = h.GetUser(issuer)
			tickets = append(tickets, ticket)
		}
		//iterate over given users
	} else {
		for _, user := range users {
			user = h.GetUserId(user)
			if user == "" {
				continue
			}
			rows, err = h.db.Query("SELECT * FROM tickets WHERE issuer = ?", user)
			for rows.Next() {
				ticket := &types.Ticket{}
				var issuer string
				rows.Scan(&ticket.Id, &issuer, &ticket.Date, &ticket.Title, &ticket.Status, &ticket.Content)
				ticket.Issuer, _ = h.GetUser(issuer)
				tickets = append(tickets, ticket)
			}
		}
	}
	return
}

func (h *Database) GetTicketIssuer(ticketId string) (userId string, err error) {
	err = h.db.QueryRow("SELECT issuer FROM tickets WHERE id = ?", ticketId).Scan(&userId)
	return
}

func (h *Database) GetTicket(ticketId string) (ticket types.Ticket, err error) {
	var issuer string
	err = h.db.QueryRow("SELECT * FROM tickets WHERE id = ?", ticketId).
		Scan(&ticket.Id, &issuer, &ticket.Date, &ticket.Title, &ticket.Status, &ticket.Content)
	ticket.Issuer, _ = h.GetUser(issuer)
	return
}

func (h *Database) UpdateTicket(id string, ticket *types.CreateTicket) error {
	_, err := h.db.Exec("UPDATE tickets SET title = ?, content = ?  WHERE id = ?",
		ticket.Title, ticket.Content, id)
	return err
}

func (h *Database) UpdateTicketStatus(status int, id string) error {
	_, err := h.db.Exec("UPDATE tickets SET status = ?  WHERE id = ?", status, id)
	return err
}
