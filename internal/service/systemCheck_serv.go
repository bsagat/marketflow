package service

import (
	"log"
	"marketflow/internal/domain"
)

type SystemHealthService struct {
	Datafetcher domain.DataFetcher
	Db          domain.Database
	Cache       domain.CacheMemory
}

func NewSystemHealthService(Datafetcher domain.DataFetcher, Db domain.Database, Cache domain.CacheMemory) *SystemHealthService {
	return &SystemHealthService{Datafetcher: Datafetcher, Db: Db, Cache: Cache}
}

func (serv *SystemHealthService) CheckHealth() []domain.ConnMsg {
	data := make([]domain.ConnMsg, 0)

	if err := serv.Datafetcher.CheckHealth(); err != nil {
		log.Println("Failed to check Datafetcher health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Datafetcher", Status: "unhealthy"})
	}

	if err := serv.Db.CheckHealth(); err != nil {
		log.Println("Failed to check Database health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Database", Status: "unhealthy"})
	}

	if err := serv.Cache.CheckHealth(); err != nil {
		log.Println("Failed to check Cache health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Cache", Status: "unhealthy"})
	}

	if len(data) == 0 {
		data = append(data, domain.ConnMsg{Status: "all connections are healthy"})
	}

	return data
}
