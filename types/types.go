package types

import(
	"time"
)

type User struct{
	Id       string
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

type Ticket struct{
	Id      string
	Issuer  User
	Date    time.Time
	Title   string
	Status  int 
	Content string
}


type Login struct{
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Register struct{
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"     validate:"required"`
}

type CreateTicket struct{
	Issuer  string `json:"issuer"  validate:"required"`
	Title   string `json:"title"   validate:"required"`
	Content string `json:"content" validate:"required"`
}

type ResponseBody struct {
    Msg interface{} `json:"Message"`
}