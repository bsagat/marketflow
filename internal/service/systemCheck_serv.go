package service

import (
	"log"
	"marketflow/internal/domain"
)

func (serv *DataModeServiceImp) CheckHealth() []domain.ConnMsg {
	data := make([]domain.ConnMsg, 0)

	if err := serv.Datafetcher.CheckHealth(); err != nil {
		log.Println("Cathed error from Datafetcher health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Datafetcher", Status: err.Error()})
	}

	if err := serv.DB.CheckHealth(); err != nil {
		log.Println("Cathed error from Database health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Database", Status: "unhealthy"})
	}

	if err := serv.Cache.CheckHealth(); err != nil {
		log.Println("Cathed error from Cache health: ", err.Error())
		data = append(data, domain.ConnMsg{Connection: "Cache", Status: "unhealthy"})
	}

	if len(data) == 0 {
		data = append(data, domain.ConnMsg{Status: "all connections are healthy"})
	}

	return data
}
