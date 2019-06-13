package models

import "time"

type AssetPackagesDataAccessObject struct{}

type AssetPackages struct {
	Id               int64     `xorm:"'id' pk autoincr notnull"`
	AssetID          string    `xorm:"'asset_id' not null index(par_ind) VARCHAR(50)"`
	AssetApiKey      string    `xorm:"'asset_apikey' unique VARCHAR(30)"`
	AssetName        string    `xorm:"'asset_name' VARCHAR(75)"`
	AssetVersion     string    `xorm:"'asset_version' VARCHAR(75)"`
	AssetReplaced    bool      `xorm:"'asset_replaced' NOT NULL DEFAULT 0"`
	AssetReplaceDate time.Time `xorm:"'asset_replace_date'"`
	Purdate          time.Time `xorm:"'purdate'"`
	Verifydate       time.Time `xorm:"'verifydate' created"`
}

var AssetPackagesDAO *AssetPackagesDataAccessObject

func (p *AssetPackagesDataAccessObject) TableName() string {
	return "asset_packages"
}
