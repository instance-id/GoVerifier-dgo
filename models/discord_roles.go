package models

type DiscordRoles struct {
	Id       int64
	RoleID   string `xorm:"'role_id' not null index(par_ind) VARCHAR(50)"`
	RoleName string `xorm:"'role_name' unique VARCHAR(30)"`
}

func (p *DiscordRoles) TableName() string {
	return "discord_roles"
}
