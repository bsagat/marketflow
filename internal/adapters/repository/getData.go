package repository

import (
	"database/sql"
	"marketflow/internal/domain"
)

func (repo *PostgresDatabase) GetLatestData(exchange, symbol string) (domain.Data, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if exchange == "All" {
		query = `
			SELECT Exchange, Pair_name, Price, StoredTime
			FROM LatestData
			WHERE Pair_name = $1
			ORDER BY StoredTime DESC
			LIMIT 1;
		`
		rows, err = repo.Db.Query(query, symbol)
	} else {
		query = `
			SELECT Exchange, Pair_name, Price, StoredTime
			FROM LatestData
			WHERE Exchange = $1 AND Pair_name = $2
			ORDER BY StoredTime DESC
			LIMIT 1;
		`
		rows, err = repo.Db.Query(query, exchange, symbol)
	}

	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	var data domain.Data
	if rows.Next() {
		if err := rows.Scan(&data.ExchangeName, &data.Symbol, &data.Price, &data.Timestamp); err != nil {
			return domain.Data{}, err
		}
		return data, nil
	}

	return domain.Data{}, nil
}
