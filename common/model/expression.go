package model

import (
	"fmt"
	"github.com/fanghongbo/gorm"
	"github.com/fanghongbo/ops-hbs/common/g"
	"strings"
)

type Expression struct {
	ID         int64  `gorm:"primary_key" json:"id"`
	Expression string `gorm:"type:varchar(1024)" json:"expression" comment:"表达式"`
	Func       string `gorm:"type:varchar(16);default:all(#1)" json:"func" comment:"max(#3) all(#3)"`
	Op         string `gorm:"type:varchar(8)" json:"op" comment:"条件"`
	RightValue string `gorm:"column:right_value;type:varchar(16)" json:"right_value" comment:"报警阈值"`
	MaxStep    int64  `gorm:"default:1" json:"max_step" comment:"间隔时间"`
	Priority   int    `gorm:"default:0" json:"priority" comment:"问题级别"`
	ActionId   int64  `gorm:"default:0" json:"action_id" comment:"动作id"`
	CreateUser string `gorm:"column:create_user;type:varchar(64)" json:"create_user" comment:"创建用户"`
	Pause      int    `gorm:"default:0" json:"pause" comment:"停用"`
	Note       string `gorm:"type:longtext" json:"note" comment:"备注"`

	Metric string            `gorm:"-" json:"metric"`
	Tags   map[string]string `gorm:"-" json:"tags"`
}

func (Expression) TableName() string {
	return "expression"
}

func (u Expression) GetAll() ([]Expression, error) {
	var (
		db     *gorm.DB
		result []Expression
		err    error
	)

	db = g.DB()
	if err = db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Expression) ParseExpression() (string, map[string]string, error) {
	var (
		left   int
		right  int
		tagStr string
		tagArr []string
		tags   map[string]string
		metric string
		exist  bool
	)

	left = strings.Index(u.Expression, "(")
	right = strings.Index(u.Expression, ")")
	tagStr = strings.TrimSpace(u.Expression[left+1 : right])
	tags = make(map[string]string)
	tagArr = strings.Fields(tagStr)

	if len(tagArr) < 2 {
		return metric, tags, fmt.Errorf("tag not enough. exp: %s", u.Expression)
	}

	for _, item := range tagArr {
		item = strings.TrimSpace(item)
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			return metric, tags, fmt.Errorf("parse %s fail", u.Expression)
		}
		tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	metric, exist = tags["metric"]
	if !exist {
		return metric, tags, fmt.Errorf("no metric give of %s", u.Expression)
	}

	delete(tags, "metric")
	return metric, tags, nil
}
