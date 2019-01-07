package tools

//MySQL数据库分页计算，传入页数和每页数量，得到sql需要的分页参数
func DBPage(number, size int) (limit int, offset int) {
	// 每页的数量，就是限制
	limit = size
	// 页数-1，再乘以，每页的数量，就等于偏移量
	offset = (number - 1) * size
	// 返回
	return limit, offset
}
