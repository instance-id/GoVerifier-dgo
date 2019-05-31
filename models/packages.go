package models

import "time"

type Packages struct {
	Id         int64     `xorm:"'id'index(par_ind) INT(10)"`
	Username   string    `xorm:"'username' not null index(par_ind) VARCHAR(50)"`
	Invoice    string    `xorm:"'invoice' unique VARCHAR(15)"`
	Packages   string    `xorm:"'packages' VARCHAR(50)"`
	Purdate    time.Time `xorm:"'purdate'"`
	Verifydate time.Time `xorm:"'verifydate' created"`
}

func (p *Packages) TableName() string {
	return "packages"
}

//func (p *Packages) FindInvoice(invoice string) (string, error) {
//
//}
