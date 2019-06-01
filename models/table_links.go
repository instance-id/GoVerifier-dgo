package models

type UserPackagesCombine struct {
	VerifiedUser `xorm:"extends"`
	UserPackages `xorm:"extends"`
}

type UserPackagesLink struct {
}
