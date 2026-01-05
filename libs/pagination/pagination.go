package pagination

import (
	"strconv"

	"gorm.io/gorm"
)

func Paginate(offset, limit int, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit).Order("created_at DESC")
	}
}

type Pagination struct {
	Offset  int
	PerPage int
	Page    int
}

func PaginationBuilder(perPage, page string) *Pagination {
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		perPageInt = 10
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	if pageInt < 1 {
		pageInt = 1
	}

	offset := (pageInt - 1) * perPageInt

	paginator := Pagination{
		Offset:  offset,
		PerPage: perPageInt,
		Page:    pageInt,
	}
	return &paginator
}

func TotalPage(totalRows, perPage int) int {
	totalPage := totalRows / perPage
	if totalRows%perPage > 0 {
		totalPage++
	}
	return totalPage
}

type PaginationResponse struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}
