package service

import (
	"errors"
	"xstation/configs"
	"xstation/models"
	"xstation/pkg/orm"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func initTables(db *gorm.DB) {
	db.AutoMigrate(
		&models.XServer{},
		&models.XDevice{},
		&models.XLink{},
		&models.XAlarm{},
		&models.XStatus{},
	)
	if models.KXStatusTabNumber > 1 {
		db.AutoMigrate(&models.XStatus1{})
	}
	if models.KXStatusTabNumber > 2 {
		db.AutoMigrate(&models.XStatus2{})
	}
	if models.KXStatusTabNumber > 3 {
		db.AutoMigrate(&models.XStatus3{})
	}
	if models.KXStatusTabNumber > 4 {
		db.AutoMigrate(&models.XStatus4{})
	}
	orm.SetDB(db)
}

var gconf = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}

func sqlInit() error {
	var err error
	var db *gorm.DB
	switch configs.Default.SQL.Name {
	case "mysql":
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: configs.Default.SQL.Address,
			// DefaultStringSize:         64,    // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}), &gconf)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(configs.Default.SQL.LiteDB), &gconf)
	case "postgresql":
		db, err = gorm.Open(postgres.Open(configs.Default.SQL.Postgre), &gconf)
	default:
	}
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("db invalid")
	}
	sqldb, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqldb.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqldb.SetMaxOpenConns(100)
	initTables(db)
	return nil
}
