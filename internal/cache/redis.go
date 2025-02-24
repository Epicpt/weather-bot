package cache

import (
	"github.com/redis/go-redis/v9"
)

type City struct {
	ID              int    `json:"id"`
	Name            string `json:"city"`
	FederalDistrict string `json:"federal_district"` // федеральный округ
	Region          string `json:"region_with_type"`
	CityDistrict    string `json:"city_district_with_type"` // район
	Street          string `json:"street_with_type"`
}

// На удаление
func NewCity(id int, name string, region string) *City {
	return &City{
		ID:     id,
		Name:   name,
		Region: region,
	}
}

type Cache struct {
	c *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{c: c}
}

// Инициализация PostgreSQL
func Init(addr string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})
	return rdb

}
