CREATE TABLE AggregatedData(
    Data_id SERIAL PRIMARY KEY,
    Pair_name VARCHAR NOT NULL,
    Exchange VARCHAR(100) NOT NULL,
    StoredTime Timestamp DEFAULT NOW(),
    Average_price FLOAT NOT NULL, 
    Min_price FLOAT NOT NULL,
    Max_price FLOAT NOT NULL
);