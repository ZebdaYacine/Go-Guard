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

type Login_Entity struct {
	Email    string
	Password string
}

type Result struct {
	User    interface{}
	Success bool
	Error   string
}
