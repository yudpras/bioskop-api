package models

type City struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

type CityInput struct {
	Name     string `json:"name" validate:"required"`
}
