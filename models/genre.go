package models

import (
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)

func GetAllGenre() []entities.Genre {
	var listGenre []entities.Genre
	err := config.DB.From("genres").
		Select(
			goqu.I("name"),
			goqu.I("slug"),
			goqu.I("image"),
		).
		Where(
			goqu.Ex{"status": 1,
		}).
		Order(goqu.I("position").Asc()).
		ScanStructs(&listGenre)

	if err != nil {
		return nil
	}

	return listGenre
}