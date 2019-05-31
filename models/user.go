package models

type ValidatedUsers struct {
	Id       int64
	Username string `xorm:"'username' not null index VARCHAR(50)"`
	Email    string `xorm:"'email' VARCHAR(75)"`
}

func (c *ValidatedUsers) TableName() string {
	return "validated_users"
}

func SearchUser() {

}
