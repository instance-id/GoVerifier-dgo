package models

import "time"

type UserPackages struct {
	Id         int64     `xorm:"'id' index(fk_verified_users_id_idx) INT(10)"`
	Username   string    `xorm:"'username' not null index(par_ind) VARCHAR(50)"`
	Invoice    string    `xorm:"'invoice' unique VARCHAR(15)"`
	Packages   string    `xorm:"'packages' VARCHAR(50)"`
	Purdate    time.Time `xorm:"'purdate'"`
	Verifydate time.Time `xorm:"'verifydate' created"`
}

func (p *UserPackages) TableName() string {
	return "user_packages"
}

//func (p *UserPackages) FindInvoice(invoice string) (string, error) {
//
//}
