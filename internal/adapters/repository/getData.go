package repository

import (
	"database/sql"
	"fmt"
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

func (repo *PostgresDatabase) GetExtremePriceByExchange(order, exchange, symbol string) (domain.Data, error) {
	var data domain.Data

	query := fmt.Sprintf(`
		SELECT Exchange, Pair_name, Price, StoredTime
		FROM LatestData
		WHERE Exchange = $1 AND Pair_name = $2
		ORDER BY Price %s
		LIMIT 1;
	`, order)

	rows, err := repo.Db.Query(query, exchange, symbol)
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

func (repo *PostgresDatabase) GetLatestDataByAllExchanges(symbol string) (domain.Data, error) {
	var data domain.Data

	rows, err := repo.Db.Query(`
		SELECT Exchange, Pair_name, Price, StoredTime
		FROM LatestData
		WHERE Pair_name = $1
		ORDER BY StoredTime DESC
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

func (repo *PostgresDatabase) GetExtremePriceByAllExchanges(order, symbol string) (domain.Data, error) {
	var data domain.Data

	// DESC для max, ASC для min

	query := fmt.Sprintf(`
		SELECT Exchange, Pair_name, Price, StoredTime
		FROM LatestData
		WHERE Pair_name = $1
		ORDER BY Price %s
		LIMIT 1;
	`, order)

	rows, err := repo.Db.Query(query, symbol)
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

func (repo *PostgresDatabase) GetExtremePriceByDuration(order, exchange, symbol string, startTime time.Time, period time.Duration) (domain.Data, error) {
	var data domain.Data
	endTime := startTime.Add(-period)

	var (
		rows  *sql.Rows
		err   error
		query string
	)

	if exchange == "All" {
		query = fmt.Sprintf(`
			SELECT Exchange, Pair_name, Price, StoredTime
			FROM LatestData
			WHERE Pair_name = $1 AND StoredTime BETWEEN $2 AND $3
			ORDER BY Price %s
			LIMIT 1;
		`, order)

		rows, err = repo.Db.Query(query, symbol, endTime, startTime)
	} else {
		query = fmt.Sprintf(`
			SELECT Exchange, Pair_name, Price, StoredTime
			FROM LatestData
			WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime BETWEEN $3 AND $4
			ORDER BY Price %s
			LIMIT 1;
		`, order)

		rows, err = repo.Db.Query(query, exchange, symbol, endTime, startTime)
	}

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

// Min by all exchange and all time
func (repo *PostgresDatabase) GetMinPriceByAllExchanges(symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: "All",
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MIN(Min_price), 0) FROM AggregatedData
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

// Min by one exchange and all time
func (repo *PostgresDatabase) GetMinPriceByExchange(exchange, symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MIN(Min_price), 0) FROM AggregatedData
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

// Min by one exchange on period
func (repo *PostgresDatabase) GetMinPriceWithDuration(exchange, symbol string, startTime time.Time, duration time.Duration) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MIN(Min_price), 0) FROM AggregatedData
		WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime BETWEEN $3 AND $4
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

// Max by all exchange all time
func (repo *PostgresDatabase) GetMaxPriceByAllExchanges(symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: "All",
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MAX(Max_price), 0) FROM AggregatedData
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


// Max by one exchange on all time
func (repo *PostgresDatabase) GetMaxPriceByExchange(exchange, symbol string) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MAX(Max_price), 0) FROM AggregatedData
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

// Max by one exchange on period
func (repo *PostgresDatabase) GetMaxPriceWithDuration(exchange, symbol string, startTime time.Time, duration time.Duration) (domain.Data, error) {
	data := domain.Data{
		ExchangeName: exchange,
		Symbol:       symbol,
	}

	rows, err := repo.Db.Query(`
		SELECT COALESCE(MAX(Max_price), 0) FROM AggregatedData
		WHERE Exchange = $1 AND Pair_name = $2 AND StoredTime BETWEEN $3 AND $4
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



