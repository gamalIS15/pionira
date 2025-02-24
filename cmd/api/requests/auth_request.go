package requests

type RegisterAuthRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FrontEndUrl string `json:"frontend_url" validate:"required,url"`
}

type ResetPasswordRequest struct {
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	Token           string `json:"token" validate:"required,min=5"`
	Meta            string `json:"meta" validate:"required"`
}
