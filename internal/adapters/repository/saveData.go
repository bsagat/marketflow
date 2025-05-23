package repository

import (
	"marketflow/internal/domain"
)

func (repo *PostgresDatabase) SaveAggregatedData(aggregatedData map[string]domain.ExchangeData) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}

	for _, data := range aggregatedData {
		_, err := repo.Db.Exec(`
		INSERT INTO AggregatedData(Pair_name, Exchange, Average_price, Min_price, Max_price)
		VALUES($1, $2, $3, $4, $5)
		`, data.Pair_name, data.Exchange, data.Average_price, data.Min_price, data.Max_price)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
