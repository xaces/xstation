package page

// Page 查询页
type Page struct {
	PageNum  uint64 `form:"pageNum"`  // 当前页码
	PageSize uint64 `form:"pageSize"` // 每页数
}