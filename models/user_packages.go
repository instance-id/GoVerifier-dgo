package models

import "time"

type UserPackagesDataAccessObject struct{}

type UserPackages struct {
	Id         int64     `xorm:"'id' index(fk_verified_users_id_idx) INT(10)"`
	Username   string    `xorm:"'username' not null index(par_ind) VARCHAR(50)"`
	Invoice    string    `xorm:"'invoice' unique VARCHAR(15)"`
	Packages   string    `xorm:"'packages' VARCHAR(75)"`
	Purdate    time.Time `xorm:"'purdate'"`
	Verifydate time.Time `xorm:"'verifydate' created"`
}

var UserPackagesDAO *UserPackagesDataAccessObject

func (p *UserPackagesDataAccessObject) TableName() string {
	return "user_packages"
}

func NewUserPackages(user *VerifiedUser, invoice string, packages string, purdate time.Time) *UserPackages {
	return &UserPackages{
		Id:       user.Id,
		Username: user.Username,
		Invoice:  invoice,
		Packages: packages,
		Purdate:  purdate,
	}
}
