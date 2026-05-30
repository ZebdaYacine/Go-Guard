package usecase

type User_Entity struct {
	User_name string `json:"user_name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,len=10"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=guest client artisant admin"`
	Sex       string `json:"sex" validate:"required,oneof=male female"`
	Picture   string `json:"picture,omitempty" validate:"omitempty,url"`
}

type Query struct {
	User User_Entity `json:"user"`
}

type Result struct {
	User    User_Entity `json:"user,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}
