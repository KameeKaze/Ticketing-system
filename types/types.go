package types

type User struct{
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

type Ticket struct{
	Issuer *User
	Date    string
	Title   string
	Priority   int 
	Type    string
	Content string
	Status  string 
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