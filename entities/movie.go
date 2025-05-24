package entities

type Category struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}
type Movie struct {
	Id	   			int    `json:"id"`
	Name   			string `json:"name"`
	Slug   			string `json:"slug"`
	Release_date   	int    `json:"release_date"`
	Rating			float64 `json:"rating"`
	//Genres 			Genre `json:"genres"`
	Image 			Image  `json:"image"`
}
type Genre struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}
type Image struct {
	Poster 	string `json:"poster"`
	Thumb  	string `json:"thumb"`
}
type PaginatedMovies struct {
	Data       []Movie `json:"data"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	Total      uint    `json:"total"`
	TotalPages int     `json:"total_pages"`
}