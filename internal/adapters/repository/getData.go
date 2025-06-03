package repository

import (
	"marketflow/internal/domain"
	"time"
)

// Gets the latest price data by exchange for specific symbol
func (repo *PostgresDatabase) GetLatestDataByExchange(exchange, symbol string) (domain.Data, error) {
	var data domain.Data

	rows, err := repo.Db.Query(`
		SELECT Exchange, Pair_name, Price, StoredTime
			FROM LatestData
		WHERE Exchange = $1 AND Pair_name = $2
		ORDER BY StoredTime DESC
		LIMIT 1;
		`, exchange, symbol)

	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&data.ExchangeName, &data.Symbol, &data.Price, &data.Timestamp); err != nil {
			return domain.Data{}, err
		}
		return data, nil
	}

	return domain.Data{}, nil
}

// Gets the latest price data for a specific symbol
func (repo *PostgresDatabase) GetLatestDataByAllExchanges(symbol string) (domain.Data, error) {
	var data domain.Data

	rows, err := repo.Db.Query(`
		SELECT Exchange, Pair_name, Price, StoredTime
		LatestData
		Pair_name = $1
		BY StoredTime DESC
		LIMIT 1;
	`, symbol)

	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&data.ExchangeName, &data.Symbol, &data.Price, &data.Timestamp); err != nil {
			return domain.Data{}, err
		}
		return data, nil
	}

	return domain.Data{}, nil
}

// Gets the average price data by exchange over all period
func (repo *PostgresDatabase) GetAveragePriceByExchange(exchange, symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
	SELECT COALESCE(AVG(Average_price), 0) FROM AggregatedData
	WHERE Exchange = $1 AND Pair_name = $2
	`, exchange, symbol)
	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data.Price); err != nil {
			return domain.Data{}, err
		}
	}

	return data, nil
}

// Gets the average price by exchange over all period
func (repo *PostgresDatabase) GetAveragePriceByAllExchanges(symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: "All",
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
	SELECT COALESCE(AVG(Average_price), 0) from AggregatedData
	WHERE Pair_name = $1
	`, symbol)
	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data.Price); err != nil {
			return domain.Data{}, err
		}
	}

	return data, nil
}

// Gets the average price within the last {duration}
func (repo *PostgresDatabase) GetAveragePriceWithDuration(exchange, symbol string, startTime time.Time, duration time.Duration) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
	SELECT COALESCE(AVG(Average_price), 0) FROM AggregatedData
	WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime BETWEEN $3 and $4
	`, exchange, symbol, startTime.Add(-duration), startTime)
	if err != nil {
		return domain.Data{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data.Price); err != nil {
			return domain.Data{}, err
		}
	}

	return data, nil
}
