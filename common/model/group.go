package model

import (
	"github.com/fanghongbo/ops-hbs/common/g"
	"github.com/jinzhu/gorm"
)

type HostGroup struct {
	ID     int64 `gorm:"primary_key" json:"id"`
	HostId int64 `gorm:"column:host_id;" json:"host_id" comment:"机器id"`
	GrpId  int64 `gorm:"column:grp_id;" json:"grp_id" comment:"机器组id"`
}

func (HostGroup) TableName() string {
	return "grp_host"
}

func (u HostGroup) GetAll() ([]HostGroup, error) {
	var (
		db     *gorm.DB
		result []HostGroup
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
