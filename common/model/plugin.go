package model

import (
	"github.com/fanghongbo/gorm"
	"github.com/fanghongbo/ops-hbs/common/g"
	"time"
)

type Plugin struct {
	ID         int64     `gorm:"primary_key" json:"id"`
	GrpId      int64     `gorm:"column:grp_id;" json:"grp_id" comment:"机器组id"`
	Dir        string    `gorm:"type:varchar(255)" json:"dir" comment:"目录名称"`
	CreateUser string    `gorm:"column:create_user;type:varchar(64)" json:"create_user" comment:"创建用户"`
	CreatedAt  time.Time `gorm:"column:create_at" json:"create_at" comment:"创建时间"`
}

func (Plugin) TableName() string {
	return "plugin_dir"
}

func (u Plugin) GetAll() ([]Plugin, error) {
	var (
		db     *gorm.DB
		result []Plugin
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
