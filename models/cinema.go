package models

type Cinema struct {
	ID       int    `json:"id"`
	CityId       int    `json:"cities_id"`
	Name     string `json:"name"`
	Address string `json:"address"`
	Phone string `json:"phone"` 
}

type CinemaInput struct {
	CityId     int `json:"cities_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address"`
	Phone    string `json:"phone"` 
}
