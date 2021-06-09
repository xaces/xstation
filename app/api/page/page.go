package page

import "xstation/pkg/orm"

// Page 查询页
type Page struct {
	PageNum  uint64 `form:"pageNum"`  // 当前页码
	PageSize uint64 `form:"pageSize"` // 每页数
}

// Where 初始化
func (s *Page) Where() *orm.DbWhere {
	var where orm.DbWhere
	return &where
}