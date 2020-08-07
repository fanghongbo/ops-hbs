package model

import (
	"github.com/fanghongbo/gorm"
	"github.com/fanghongbo/ops-hbs/common/g"
	"time"
)

type Strategy struct {
	ID         int64  `gorm:"primary_key" json:"id"`
	Func       string `gorm:"type:varchar(16);default:all(#1)" json:"func" comment:"max(#3) all(#3)"`
	Op         string `gorm:"type:varchar(8)" json:"op" comment:"条件"`
	RightValue string `gorm:"column:right_value;type:varchar(16)" json:"right_value" comment:"报警阈值"`
	MaxStep    int64  `gorm:"default:1" json:"max_step" comment:"间隔时间"`
	Priority   int    `gorm:"default:0" json:"priority" comment:"问题级别"`
	Note       string `gorm:"type:longtext" json:"note" comment:"备注"`
	Metric     string `gorm:"type:varchar(128)" json:"metric"`
	Tags       string `gorm:"type:varchar(256)" json:"tags"`
	RunBegin   string `gorm:"column:run_begin;type:varchar(16)" json:"run_begin" comment:"运行开始"`
	RunEnd     string `gorm:"column:run_end;type:varchar(16)" json:"run_end" comment:"运行结束"`
	TplId      int64  `gorm:"column:tpl_id;" json:"tpl_id" comment:"模版id"`
}

func (Strategy) TableName() string {
	return "strategy"
}

func (u Strategy) GetAll() ([]Strategy, error) {
	var (
		db     *gorm.DB
		now    string
		result []Strategy
		err    error
	)

	now = time.Now().Format("15:04")
	db = g.DB()
	if err = db.Where("run_begin = ? and run_begin = ?", "", "").
		Or("run_begin <= ? and run_end >= ?", now, now).
		Or("run_begin > run_end and !(run_begin > ? and run_end < ?)", now, now).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u Strategy) GetBuiltinMetrics(templateIds string) ([]BuiltinMetric, error) {
	var (
		db      *gorm.DB
		result  []Strategy
		metrics []string
		data    []BuiltinMetric
		err     error
	)

	metrics = []string{"net.port.listen", "proc.num", "du.bs", "url.check.health"}

	db = g.DB()
	if err = db.Where("tpl_id in (?) and metric in (?)", templateIds, metrics).Find(&result).Error; err != nil {
		return nil, err
	}

	data = []BuiltinMetric{}
	for _, item := range result {
		data = append(data, BuiltinMetric{
			Metric: item.Metric,
			Tags:   item.Tags,
		})
	}
	return data, nil
}
