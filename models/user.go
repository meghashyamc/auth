package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Email          *string   `json:"email" validate:"required,email"`
	Password       *string   `json:"password" validate:"required,min=16,max=1000" gorm:"-"`
	PasswordDigest string    `json:"-"`
	FirstName      *string   `json:"first_name" validate:"required,alpha,min=1,max=100"`
	LastName       *string   `json:"last_name" validate:"omitempty,alpha,max=100"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
