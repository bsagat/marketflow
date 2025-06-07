package repository

import (
	"fmt"
	"log/slog"
	"marketflow/internal/domain"
)

func (repo *PostgresDatabase) SaveAggregatedData(aggregatedData map[string]domain.ExchangeData) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO AggregatedData(Pair_name, Exchange, StoredTime, Average_price, Min_price, Max_price)
		VALUES($1, $2, $3, $4, $5, $6)
		`)
	if err != nil {
		tx.Rollback()
		slog.Error("Failed to prepare statement", "error", err.Error())
		return err
	}
	defer stmt.Close()
	fmt.Println(len(aggregatedData))
	for _, data := range aggregatedData {
		_, err := stmt.Exec(data.Pair_name, data.Exchange, data.Timestamp, data.Average_price, data.Min_price, data.Max_price)
		if err != nil {
			tx.Rollback()
			slog.Error("Failed to execute statement", "pair", data.Pair_name, "exchange", data.Exchange, "error", err.Error())
			return err
		}
	}
	slog.Info("Committing transaction", "records", len(aggregatedData))
	return tx.Commit()
}

func (repo *PostgresDatabase) SaveLatestData(latestData map[string]domain.Data) error {
	tx, err := repo.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO LatestData (Exchange, Pair_name, Price, StoredTime)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (Exchange, Pair_name) DO UPDATE
		SET Price = EXCLUDED.Price,
    	StoredTime = EXCLUDED.StoredTime;
		`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, data := range latestData {
		if _, err := stmt.Exec(data.ExchangeName, data.Symbol, data.Price, data.Timestamp); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
