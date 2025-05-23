package repository

func (repo *PostgresDatabase) CheckHealth() error {
	if err := repo.Db.Ping(); err != nil {
		return err
	}
	return nil
}
