package models

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=16,max=1000"`
}

type SendConfirmMailRequest struct {
	ID string `json:"user_id" validate:"required,uuid4"`
}
