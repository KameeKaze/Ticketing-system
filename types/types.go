package types

import(
	"time"
)

type User struct{
	Id       string `json:"id"`
	Name     string `json:"username" validate:"required"`
	Role     string `json:"role"`
}

type Ticket struct{
	Id      string    `json:"id"`
	Issuer  User	  `json:"issuer"`
	Date    time.Time `json:"date"`
	Title   string    `json:"title"`
	Status  int       `json:"status"`
	Content string    `json:"content"`
}


type Login struct{
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePassword struct{
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newpassword" validate:"required"`
}

type Register struct{
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"     validate:"required"`
}

type CreateTicket struct{
	Issuer  string `json:"issuer"`
	Title   string `json:"title"   validate:"required"`
	Content string `json:"content" validate:"required"`
}

type ResponseBody struct {
    Msg interface{} `json:"Message"`
}