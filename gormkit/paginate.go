package gormkit

import "gorm.io/gorm"

type PageLabel struct {
	PageNum    int64 `json:"page_num"`
	PageSize   int64 `json:"page_size"`
	TotalPage  int64 `json:"total_page"`
	TotalCount int64 `json:"total_count"`
}

// GetPaginateData 获取分页数据。filter 为分页标签, m 为查询的表的model，result 传递指针获取结果列表
func GetPaginateData(db *gorm.DB, label *PageLabel, m interface{}, result interface{}) (err error) {
	offset := getOffset(label.PageNum, label.PageSize)
	if err = db.Offset(offset).Limit(int(label.PageSize)).Find(result).Error; err != nil {
		return
	}
	
	db = db.Offset(-1).Limit(-1)
	if err = db.Model(m).Count(&label.TotalCount).Error; err != nil {
		return
	}
	if label.PageSize <= 0 {
		return
	}
	label.TotalPage = getTotalPage(label.TotalCount, label.PageSize)
	return
}

func getTotalPage(totalCount, pageSize int64) int64 {
	return (totalCount + pageSize - 1) / pageSize
}

func getOffset(pageNum, pageSize int64) int {
	return int((pageNum - 1) * pageSize)
}
