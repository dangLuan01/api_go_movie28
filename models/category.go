package models

import (
	"fmt"
	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/doug-martin/goqu/v9"
)
func GetAllCategory() []entities.Category  {
	var categories []entities.Category
	err := config.DB.
		From("categories").
		Order(goqu.I("position").Asc()).
		ScanStructs(&categories)
	if err != nil {
		fmt.Println("Err:", err)
		return nil
	}

	return categories
}