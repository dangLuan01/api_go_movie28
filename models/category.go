package models

import (
	"fmt"
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)
func GetAllCategory() []entities.Category  {
	var categories []entities.Category
	err := config.DB.
		From("categories").
		Order(goqu.I("position").Asc()). // optional: sắp xếp theo ID
		ScanStructs(&categories)
	if err != nil {
		fmt.Println("Lỗi khi lấy danh sách thể loại:", err)
		return nil
	}

	return categories
}