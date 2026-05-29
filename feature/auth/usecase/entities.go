package usecase

type User_Entity struct {
	User_name string `validate:"required,min=3,max=50"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"required,len=10"`
	Password  string `validate:"required,min=8"`
	Role      string `validate:"required,oneof=guest client artisant admin"`
	Sex       string `validate:"required,oneof=male female"`
	Picture   string `validate:"omitempty,url"`
}

type Query struct {
	User User_Entity
}

type Result struct {
	User    User_Entity
	Message string
}
