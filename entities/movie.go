package entities

type Category struct {
	Id    			int    		`json:"id"`
	Name  			string 		`json:"name"`
	Slug  			string 		`json:"slug"`
}
type Movie struct {
	Name   			string 		`json:"name"`
	Slug   			string 		`json:"slug"`
	Type  			string 		`json:"type"`
	Release_date   	int    		`json:"release_date"`
	Rating			float32 	`json:"rating"`
	Content 		string 		`json:"content,omitempty"`
	Runtime 		string    	`json:"runtime,omitempty"`
	Age 			string 		`json:"age,omitempty"`
	Trailer 		string 		`json:"trailer,omitempty"`
	//ImageRaw  		string   	`db:"image_raw" json:"-"`
	Image 			Image  		`json:"image"`
	//GenresRaw    	string   	`db:"genres_raw" json:"-"`
	Genres 			[]Genre 	`json:"genres,omitempty"`
	Country 		[]Country 	`json:"country,omitempty"`
	Servers 		[]Server	`json:"servers,omitempty"`
}
type Image struct {
	Poster 			string 		`json:"poster,omitempty"`
	Thumb  			string 		`json:"thumb,omitempty"`
}
type Genre struct {
	Id 				int         `json:"id,omitempty"`
	Name  			string 		`json:"name"`
	Slug  			string 		`json:"slug,omitempty"`
	Image 			string 		`json:"image,omitempty"`
}
type Country struct {
	Name 			string 		`json:"name"`
	Slug 			string 		`json:"slug"`
}
type Episode struct {
	Server_id 		int    		`json:"server_id"`
	Episode			string 		`json:"episode"`
	Hls  			string 		`json:"hls"`
}
type Server struct {
	Id    			int    		`json:"id"`
	Name  			string 		`json:"name"`
	Episodes 		[]Episode 	`json:"episodes"`
}
type GenreWithMovies struct {
	Name			string 		`json:"name"`
	Slug			string 		`json:"slug"`
	Image			string 		`json:"image"`
	Total_Movies 	int 		`json:"total_movies"`
}

type PaginatedMovies struct {
	Movie       	[]Movie 	`json:"movies"`
	Page       		int     	`json:"page"`
	PageSize   		int     	`json:"page_size"`
	TotalPages 		int			`json:"total_pages"`
}
type MovieRaw struct {
	Id   			int    `json:"id,omitempty"`
	Name  			string `json:"name"`
	Slug  			string `json:"slug"`
	Type  			string `json:"type"`
	Release_date 	int    `json:"release_date"`
	Rating			float64 `json:"rating"`
	Content 		string `json:"content,omitempty"`
	Runtime 		string `json:"runtime,omitempty"`
	Age 			string `json:"age,omitempty"`
	Trailer 		string `json:"trailer,omitempty"`
	Thumb 			string `json:"thumb"`
	Poster			string `json:"poster"`
	Genre_name 		string
}
type Collection struct {
	Name 			string 	`json:"name"`
	Slug 			string 	`json:"slug"`
	Image 			string 	`json:"image"`
}
type ThemeInfo struct {
	Id 			int		`json:"id"`
	Name 		string	`json:"name"`
	Genre_id 	*int	`json:"genre_id,omitempty"`
	Country_id	*int	`json:"country_id,omitempty"`
	Type 		*string	`json:"type,omitempty"`
	Year 		*int	`json:"year,omitempty"`
	Limit		int		`json:"limit"`
}
type ThemeWithMovies struct {
	ThemeInfo 			ThemeInfo		`json:"theme"`
	PaginatedMovies 	PaginatedMovies `json:"data"`
}
type PagiateTheme struct {
	ThemeWithMovies []ThemeWithMovies `json:"datas"`
	Page       		int     	`json:"page,omitempty"`
	PageSize   		int     	`json:"page_size,omitempty"`
	TotalPages 		int			`json:"total_pages,omitempty"`
}