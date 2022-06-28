package types

type User struct{
	Name     string
	Emain    string
	Password string
	Role     string
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
	Username string `json:"username"`
	Password string `json:"password"`
}