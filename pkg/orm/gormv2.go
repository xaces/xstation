package orm

import "gorm.io/gorm"

var _db *gorm.DB

// whereValues 分页条件
type whereValues struct {
	Where string
	Value []interface{}
}

// ModelWhere 分页条件
type DbWhere struct {
	Wheres []whereValues
	Orders []string
}

// H 多列处理
type H map[string]interface{}

// Append  添加条件
func (p *DbWhere) Append(where string, value ...interface{}) {
	var w whereValues
	w.Where = where
	w.Value = value
	p.Wheres = append(p.Wheres, w)
}

// SetDB gorm对象
func SetDB(db *gorm.DB) {
	_db = db
}

// SetDB gorm对象
func DB() *gorm.DB {
	return _db
}

// DbCount 数目
func DbCount(model, where interface{}) uint64 {
	var count int64
	_db.Model(model).Where(where).Count(&count)
	return uint64(count)
}

// DbCreate 创建
func DbCreate(model interface{}) error {
	return _db.Create(model).Error
}

// DbSave 保存
func DbSave(value interface{}) error {
	return _db.Save(value).Error
}

// DbUpdateModel 更新
func DbUpdateModel(model interface{}) error {
	return _db.Model(model).Updates(model).Error
}

// DbUpdateModelBy 条件更新
func DbUpdateModelBy(model interface{}, where string, args ...interface{}) error {
	return _db.Where(where, args...).Updates(model).Error
}

// DbUpdatesById 更新
func DbUpdateById(model, id interface{}) error {
	return _db.Where("id = ?", id).Updates(model).Error
}

// DbUpdateColById 单列更新
func DbUpdateColById(model, id interface{}, column string, value interface{}) error {
	return _db.Model(model).Where("id = ?", id).Update(column, value).Error
}

// DbUpdateColsById 更新多列
// 用于0不更新
func DbUpdateColsById(model, id interface{}, value H) error {
	return _db.Model(model).Where("id = ?", id).Updates(value).Error
}

// DbUpdateColsBy 更新多列
// 用于0不更新
func DbUpdateColsBy(model interface{}, value H, where string, args ...interface{}) error {
	return _db.Model(model).Where(where, args...).Updates(value).Error
}

// DbUpdateByIds 批量更新
// ids id数组
func DbUpdateByIds(model, ids interface{}, value H) error {
	return _db.Model(model).Where("id in (?)", ids).Updates(value).Error
}

// DbDeletes 批量删除
// ids id数组 []
func DbDeletes(model, ids interface{}) error {
	return _db.Delete(model, ids).Error
}

// DbDeleteBy 删除
func DbDeleteBy(model interface{}, where string, args ...interface{}) (count int64, err error) {
	db := _db.Where(where, args...).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DbFirstBy 指定条件查找
func DbFirstBy(out interface{}, where string, args ...interface{}) (err error) {
	err = _db.Where(where, args...).First(out).Error
	return
}

// DbFirstById 查找
func DbFirstById(id uint64, out interface{}) error {
	return _db.First(out, id).Error
}

// DbFirstWhere 查找
func DbFirstWhere(out, where interface{}) error {
	return _db.Where(where).First(out).Error
}

// DbFind 多个查找
func DbFind(out interface{}, orders ...string) error {
	db := _db
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// DbFindBy 多个条件查找
func DbFindBy(out interface{}, where string, args ...interface{}) (int64, error) {
	db := _db.Where(where, args...).Find(out)
	return db.RowsAffected, db.Error
}

// DbPage 分页
type dbPage struct {
	Db         *gorm.DB
	TotalCount int64
}

// Find 分页
func (o *dbPage) Find(pageIndex, pageSize uint64, out interface{}, conds ...interface{}) (int64, error) {
	if o.TotalCount <= 0 {
		return 0, nil
	}
	if pageSize > 0 {
		return o.TotalCount, o.Db.Offset(int((pageIndex-1)*pageSize)).Limit(int(pageSize)).Find(out, conds...).Error
	}
	return o.TotalCount, o.Db.Find(out, conds...).Error
}

// Find 分页
func (o *dbPage) Scan(pageIndex, pageSize uint64, out interface{}) (int64, error) {
	if o.TotalCount <= 0 {
		return 0, nil
	}
	if pageSize > 0 {
		return o.TotalCount, o.Db.Offset(int((pageIndex - 1) * pageSize)).Limit(int(pageSize)).Scan(out).Error
	}
	return o.TotalCount, o.Db.Scan(out).Error
}

// Preload 关联加载
func (o *dbPage) Preload(preloads ...string) *dbPage {
	if len(preloads) > 0 {
		for _, preload := range preloads {
			o.Db.Preload(preload)
		}
	}
	return o
}

// Joins join
func (o *dbPage) Joins(query string, args ...interface{}) *dbPage {
	o.Db = o.Db.Joins(query, args...)
	return o
}

// DbPage
func DbPage(model interface{}, where *DbWhere) *dbPage {
	db := _db.Model(model)
	if where != nil {
		for _, wo := range where.Wheres {
			if wo.Where != "" {
				db = db.Where(wo.Where, wo.Value...)
			}
		}
	}
	if len(where.Orders) > 0 {
		for _, order := range where.Orders {
			db = db.Order(order)
		}
	}
	var totalCount int64
	if db.Count(&totalCount).Error != nil {
		totalCount = 0
	}
	return &dbPage{Db: db, TotalCount: totalCount}
}
