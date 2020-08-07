package model

import (
	"github.com/fanghongbo/gorm"
	"github.com/fanghongbo/ops-hbs/common/g"
	"time"
)

type Host struct {
	ID            int64     `gorm:"primary_key" json:"id"`
	Hostname      string    `gorm:"column:hostname;type:varchar(255)" json:"hostname" comment:"主机名"`
	Ip            string    `gorm:"type:varchar(16)" json:"ip" comment:"主机ip"`
	AgentVersion  string    `gorm:"column:agent_version;type:varchar(16)" json:"agent_version" comment:"agent版本"`
	PluginVersion string    `gorm:"column:plugin_version;type:varchar(128)" json:"plugin_version" comment:"插件版本"`
	MaintainBegin int64     `gorm:"column:maintain_begin;default:0" json:"maintain_begin" comment:"维护开始时间"`
	MaintainEnd   int64     `gorm:"column:maintain_end;default:0" json:"maintain_end" comment:"维护结束时间"`
	UpdateAt      time.Time `gorm:"column:update_at" json:"update_at" comment:"更新时间"`
}

func (Host) TableName() string {
	return "host"
}

func (u Host) GetAll() ([]Host, error) {
	var (
		db     *gorm.DB
		result []Host
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u Host) GetHostNotInMaintain() ([]Host, error) {
	var (
		db     *gorm.DB
		result []Host
		err    error
	)
	now := time.Now().Unix()

	db = g.DB()
	if err = db.Where("maintain_begin > ? or maintain_end < ?", now, now).Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func (u *Host) IsExist() bool {
	var (
		db   *gorm.DB
		host Host
	)

	db = g.DB()
	if db.Where("hostname = ?", u.Hostname).First(&host).RecordNotFound() {
		return false
	}

	return true
}

func (u *Host) Create() error {
	var (
		db  *gorm.DB
		err error
	)

	db = g.DB()
	if err = db.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *Host) Update() error {
	var (
		db   *gorm.DB
		host Host
		err  error
	)

	db = g.DB()
	if err = db.Where("hostname = ?", u.Hostname).First(&host).Error; err != nil {
		return err
	}

	u.ID = host.ID
	if err = db.Save(u).Error; err != nil {
		return err
	}

	return nil
}
