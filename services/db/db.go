package db

import (
	"fmt"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meghashyamc/auth/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBClient struct {
	db *gorm.DB
}

func PGDSN(dbhost, dbusername, dbpassword, dbname, dbport string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbhost, dbusername, dbpassword, dbname, dbport)
}

func NewClient() (*DBClient, error) {

	dsn := PGDSN("localhost", os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"), os.Getenv("PG_DB"), os.Getenv("PG_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to open DB: %v", err)
		return nil, err
	}

	return &DBClient{db: db}, nil
}

func (dbc *DBClient) CreateUser(user *models.User) (uuid.UUID, error) {
	timeNowUTC := time.Now().UTC()
	user.CreatedAt = timeNowUTC
	user.UpdatedAt = timeNowUTC
	if err := dbc.db.Create(user).Error; err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not create user")

		return uuid.Nil, err
	}

	return user.ID, nil
}

func (dbc *DBClient) GetUserByEmail(email string) (bool, *models.User, error) {

	user := &models.User{}

	result := dbc.db.Where("email=?", email).First(&user)
	if result.Error != nil {
		return true, user, nil
	}
	if result.Error == gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{"err": result.Error, "email": email}).Info("could not find user")
		return false, nil, nil
	}

	return false, nil, result.Error

}

func (dbc *DBClient) GetUserByID(id uuid.UUID) (bool, *models.User, error) {

	user := &models.User{}

	result := dbc.db.Where("id=?", id).First(&user)
	if result.Error != nil {
		return true, user, nil
	}
	if result.Error == gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{"err": result.Error, "id": id.String()}).Info("could not find user")
		return false, nil, nil
	}

	return false, nil, result.Error

}

func (dbc *DBClient) UpdateUserConfirmationToken(id uuid.UUID, confirmationToken string, tokenValidity time.Time) error {
	result := dbc.db.Model(&models.User{}).Updates(models.User{ID: id, ConfirmationToken: confirmationToken, ConfirmationTokenValidity: tokenValidity})
	if result.Error != nil {
		log.WithFields(log.Fields{"err": result.Error, "user_id": id}).Error("error updating confirmation token")
		return result.Error
	}
	return nil
}

func (dbc *DBClient) UpdateUserConfirmationMailSent(id uuid.UUID) error {
	result := dbc.db.Model(&models.User{}).Where("id = ?", id).Update("confirmation_mail_sent", true)
	if result.Error != nil {
		log.WithFields(log.Fields{"err": result.Error, "user_id": id}).Error("error updating confirmation email sent status")
		return result.Error
	}
	return nil
}
