package db

import (
	"fmt"
	"time"

	"github.com/wlgd/xutils/orm"
)

const (
	timeFormat = "2006-01-02 15:04:05"
	partFormat = "20060102 150405"
)

type schemapartitions struct {
	PartitionName  string `gorm:"column:PARTITION_NAME;"`
	PartitionDesc  string `gorm:"column:PARTITION_DESCRIPTION;"`
	PartitionExpre string `gorm:"column:PARTITION_EXPRESSION;"`
	Rows           uint   `gorm:"column:TABLE_ROWS;"`
}

type Partition struct {
	RetentionDays int
	Table         string
	PrimaryKey    string
}

func (p *Partition) query() []schemapartitions {
	var data []schemapartitions
	orm.DB().Raw(fmt.Sprintf("SELECT PARTITION_NAME, PARTITION_DESCRIPTION, PARTITION_EXPRESSION, TABLE_ROWS FROM information_schema.partitions WHERE table_name = '%s';", p.Table)).Scan(&data)
	return data
}

func (p *Partition) namExpress(t time.Time) (string, string) {
	name := "p" + t.Format(partFormat)[2:8]
	where := t.AddDate(0, 0, 1).Format(timeFormat)[:10]
	return name, where
}

var gPartInitSQL string = `
ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))
(
 PARTITION %s VALUES LESS THAN (TO_DAYS('%s'))
)
`

func (p *Partition) init(t time.Time) string {
	p0, w0 := p.namExpress(t)
	orm.DB().Exec(fmt.Sprintf(gPartInitSQL, p.Table, p.PrimaryKey, p0, w0))
	return p0
}

func (p *Partition) exec() {
	parts := p.query()
	t := time.Now().AddDate(0, 0, -1*p.RetentionDays+1)
	if parts[0].PartitionName == "" {
		p0 := p.init(time.Now())
		parts[0] = schemapartitions{PartitionName: p0, Rows: 0}
	}
	pdel, _ := p.namExpress(t)
	for _, v := range parts {
		if v.PartitionName < pdel {
			orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s;", p.Table, v.PartitionName)) // 删除过时分区
		}
	}
	isExistPart := func(p string) bool {
		for _, v := range parts {
			if v.PartitionName == p {
				return true
			}
		}
		return false
	}
	// 创建新分区
	for i := 0; i < p.RetentionDays+1; i++ {
		part, where := p.namExpress(t.AddDate(0, 0, i))
		if isExistPart(part) || part < parts[0].PartitionName {
			continue
		}
		// 创建新分区
		orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s ADD PARTITION( PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));", p.Table, part, where))
	}
}

func NewPartition(table, primaryKey string, days int) *Partition {
	return &Partition{RetentionDays: days, Table: table, PrimaryKey: primaryKey}
}
