package models

type City struct {
	ID              int    `json:"id"`
	Name            string `json:"city"`
	FederalDistrict string `json:"federal_district"` // федеральный округ
	Region          string `json:"region_with_type"`
	CityDistrict    string `json:"city_district_with_type"` // район
	Street          string `json:"street_with_type"`
	Country         string `json:"country"`
}
