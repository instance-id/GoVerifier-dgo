package models

import (
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"github.com/sarulabs/di/v2"
)

type VerifiedUserDataAccessObject struct{}

// --- User data container object -----------------------------------------------------------------
type VerifiedUser struct {
	Id       int64
	Username string `xorm:"'username' not null index VARCHAR(50)"`
	Email    string `xorm:"'email' VARCHAR(75)"`
}

// --- User data access object from outside functions ---------------------------------------------
var VerifiedUserDAO *VerifiedUserDataAccessObject

// --- Specify table name in database -------------------------------------------------------------
func (d *VerifiedUserDataAccessObject) TableName() string {
	return "verified_users"
}

// --- Create new user object --------------------------------------------------------------------
func NewVerifiedUser(username string, email string) *VerifiedUser {
	return &VerifiedUser{
		Username: username,
		Email:    email,
	}
}

// --- Add new user to database -------------------------------------------------------------------
func (d *VerifiedUserDataAccessObject) AddUser(user *VerifiedUser, di di.Container) {
	db := DatabaseAccessContainer(di)
	_, err := db.Table(VerifiedUserDAO.TableName()).InsertOne(user)
	LogFatalf("Unable to insert user", err)
	log := LogAccessContainer(di)
	log.Infof("Data from insert: %v", user.Id)
}
