package model

import (
	"github.com/fanghongbo/gorm"
	"github.com/fanghongbo/ops-hbs/common/g"
	"time"
)

type GroupTemplate struct {
	ID       int64  `gorm:"primary_key" json:"id"`
	TplId    int64  `gorm:"column:tpl_id;" json:"tpl_id" comment:"模版id"`
	GrpId    int64  `gorm:"column:grp_id;" json:"grp_id" comment:"机器组id"`
	BindUser string `gorm:"column:bind_user;type:varchar(64)" json:"bind_user" comment:"绑定用户"`
}

func (GroupTemplate) TableName() string {
	return "grp_tpl"
}

func (u GroupTemplate) GetAll() ([]GroupTemplate, error) {
	var (
		db     *gorm.DB
		result []GroupTemplate
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

type Template struct {
	ID         int64     `gorm:"primary_key" json:"id"`
	ParentId   int64     `gorm:"column:parent_id;" json:"parent_id" comment:"父模版id"`
	ActionId   int64     `gorm:"column:action_id;" json:"action_id" comment:"动作id"`
	TplName    string    `gorm:"column:tpl_name;type:varchar(255)" json:"tpl_name" comment:"模版名称"`
	CreateUser string    `gorm:"column:create_user;type:varchar(64)" json:"create_user" comment:"创建用户"`
	CreatedAt  time.Time `gorm:"column:create_at" json:"create_at" comment:"创建时间"`
}

func (Template) TableName() string {
	return "tpl"
}

func (u Template) GetAll() ([]Template, error) {
	var (
		db     *gorm.DB
		result []Template
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

type HostTemplate struct {
	HostId int64 `gorm:"column:host_id;" json:"host_id" comment:"主机id"`
	TplId  int64 `gorm:"column:tpl_id;" json:"tpl_id" comment:"模版id"`
}

func (u HostTemplate) GetAll() ([]HostTemplate, error) {
	var (
		db     *gorm.DB
		result []HostTemplate
		err    error
	)

	db = g.DB()
	if err = db.Table("grp_tpl").Select("grp_tpl.tpl_id, grp_host.host_id").Joins("left join grp_host on grp_tpl.grp_id = grp_host.grp_id").Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
