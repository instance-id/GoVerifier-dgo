package models

import (
	"github.com/go-xorm/xorm"

	"github.com/instance-id/GoVerifier-dgo/components"
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
	db := d.DataAccessContainer(di)
	_, err := db.Table(VerifiedUserDAO.TableName()).InsertOne(user)
	ErrCheckf("Unable to insert user", err)
}

func SearchUser() {

}

// --- Retrieve database connection session from dependency injection container -------------------
func (d *VerifiedUserDataAccessObject) DataAccessContainer(di di.Container) *xorm.Engine {
	db, err := di.SubContainer()
	ErrCheckf("Error accessing DI container within User Model: ", err)

	dba := db.Get("db").(*components.XormDB).Engine
	return dba
}
