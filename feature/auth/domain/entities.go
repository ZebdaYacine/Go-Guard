package domain

type User_Entity struct {
	User_name string
	Email     string
	Phone     string
	Password  string
	Role      string
	Sex       string
	Picture   string
}

type Query struct {
	User User_Entity
}

type Result struct {
	User    User_Entity
	Success bool
}
