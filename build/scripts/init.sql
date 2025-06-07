CREATE TABLE AggregatedData(
    Data_id SERIAL PRIMARY KEY,
    Pair_name VARCHAR NOT NULL,
    Exchange VARCHAR(100) NOT NULL,
    StoredTime TimestampTZ DEFAULT NOW(),
    Average_price FLOAT NOT NULL, 
    Min_price FLOAT NOT NULL,
    Max_price FLOAT NOT NULL
);

CREATE TABLE LatestData(
    Exchange VARCHAR(100) NOT NULL,
    Pair_name VARCHAR NOT NULL,
    Price FLOAT NOT NULL,
    StoredTime BIGINT NOT NULL,
    CONSTRAINT unique_exchange_pair UNIQUE (Exchange, Pair_name)
);

-- Automatically deletes rows older than 7 weeks from expire_table after each insert
CREATE FUNCTION expire_table_delete_old_rows() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  DELETE FROM AggregatedData WHERE StoredTime < NOW() - INTERVAL '7 weeks'; -- or more,  for example '4 years'
  RETURN NEW;
END;
$$;

CREATE TRIGGER expire_table_delete_old_rows_trigger
    AFTER INSERT ON AggregatedData
    EXECUTE PROCEDURE expire_table_delete_old_rows();
