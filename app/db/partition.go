package db

import (
	"fmt"
	"log"
	"time"

	"github.com/wlgd/xutils/orm"
)

const (
	timeFormat = "2006-01-02 15:04:05"
	partFormat = "20060102 150405"
)

type schemaPartition struct {
	Name       string `gorm:"column:PARTITION_NAME;"`
	Desc       string `gorm:"column:PARTITION_DESCRIPTION;"`
	Expression string `gorm:"column:PARTITION_EXPRESSION;"`
	Rows       uint   `gorm:"column:TABLE_ROWS;"`
}

func partNameExpress(t time.Time) (string, string) {
	name := "p" + t.Format(partFormat)[2:8]
	where := t.AddDate(0, 0, 1).Format(timeFormat)[:10]
	return name, where
}

type dbPartition struct {
	Table string
}

func (o *dbPartition) query() []schemaPartition {
	var data []schemaPartition
	orm.DB().Raw(fmt.Sprintf("SELECT PARTITION_NAME, PARTITION_DESCRIPTION, PARTITION_EXPRESSION, TABLE_ROWS FROM information_schema.partitions WHERE table_name = '%s';", o.Table)).Scan(&data)
	return data
}

var gPartInitSQL string = `
ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))
(
 PARTITION %s VALUES LESS THAN (TO_DAYS('%s'))
)
`
func (o *dbPartition) init(t time.Time, primaryKey string) string {
	p0, w0 := partNameExpress(t)
	orm.DB().Exec(fmt.Sprintf(gPartInitSQL, o.Table, primaryKey, p0, w0))
	return p0
}

func (o *dbPartition) AlterRange(primaryKey string, days int) {
	log.Printf("%s partiton\n", o.Table)
	parts := o.query()
	t := time.Now().AddDate(0, 0, -1*days+1)
	if parts[0].Name == "" {
		p0 := o.init(time.Now(), primaryKey)
		parts[0] = schemaPartition{Name: p0, Rows: 0}
	}
	pdel, _ := partNameExpress(t)
	for _, v := range parts {
		if v.Name < pdel {
			orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s;", o.Table, v.Name)) // 删除过时分区
		}
	}
	isExistPart := func(p string) bool {
		for _, v := range parts {
			if v.Name == p {
				return true
			}
		}
		return false
	}
	// 创建新分区
	for i := 0; i < days+1; i++ {
		part, where := partNameExpress(t.AddDate(0, 0, i))
		if part < parts[0].Name || isExistPart(part) {
			continue
		}
		// 创建新分区
		orm.DB().Exec(fmt.Sprintf("ALTER TABLE %s ADD PARTITION( PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));", o.Table, part, where))
	}
}

func PartTable(table string) *dbPartition {
	return &dbPartition{Table: table}
}
