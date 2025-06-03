package repository

import (
	"database/sql"
	"fmt"
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

func (repo *PostgresDatabase) GetExtremePrice(op, exchange, symbol string, period int) (domain.Data, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if exchange == "All" {
		query = fmt.Sprintf(`
			SELECT Exchange, Pair_name, %s(Price), MAX(StoredTime)
			FROM LatestData
			WHERE Pair_name = $1 AND StoredTime >= NOW() - INTERVAL '%d seconds'
			GROUP BY Exchange, Pair_name
			ORDER BY %s(Price) DESC
			LIMIT 1;
		`, op, period, op)
		rows, err = repo.Db.Query(query, symbol)
	} else {
		query = fmt.Sprintf(`
			SELECT Exchange, Pair_name, %s(Price), MAX(StoredTime)
			FROM LatestData
			WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime >= NOW() - INTERVAL '%d seconds'
			GROUP BY Exchange, Pair_name
			LIMIT 1;
		`, op, period)
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

func (repo *PostgresDatabase) GetAveragePrice(exchange, symbol string, period int) (domain.Data, error) {
	var (
		query string
		row   *sql.Row
	)

	if exchange == "All" {
		query = `
			SELECT 'All', Pair_name, AVG(Price), MAX(StoredTime)		
			FROM LatestData
			WHERE Pair_name = $1 AND StoredTime >= NOW() - INTERVAL '%d seconds'
			GROUP BY Pair_name
			LIMIT 1;
		`
		query = fmt.Sprintf(query, period)
		row = repo.Db.QueryRow(query, symbol)
	} else {
		query = `
			SELECT Exchange, Pair_name, AVG(Price), MAX(StoredTime)
			FROM LatestData
			WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime >= NOW() - INTERVAL '%d seconds'
			GROUP BY Exchange, Pair_name
			LIMIT 1;
		`
		query = fmt.Sprintf(query, period)
		row = repo.Db.QueryRow(query, exchange, symbol)
	}

	var data domain.Data
	err := row.Scan(&data.ExchangeName, &data.Symbol, &data.Price, &data.Timestamp)
	if err != nil {
		return domain.Data{}, err
	}

	return data, nil
}
