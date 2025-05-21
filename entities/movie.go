package entities

type Category struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	// Status int   `json:"status"`
	// Position int `json:"position"`
}
type Movie struct {
	Id	   int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Release_date   int    `json:"release_date"`
	//Image  string `json:"image"`
}