package entities

type Category struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}
type Movie struct {
	Id	   int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Year   int    `json:"year"`
	Image  string `json:"image"`
}