# 📈 MarketFlow

**MarketFlow** is a server-side application on Go for real-time aggregation and processing of market data. It uses microservice architecture, Redis and PostgreSQL.

## 🚀 Features

- 🧪 Support for Test Mode and Live Mode
- 🔌 Connection to multiple exchanges (exchange1, exchange2, exchange3)
- 🧠 Data caching in Redis
- 💾 Aggregated data storage in PostgreSQL
- 🌐 REST API for data management and retrieval

## ⚙️ Installation and startup

### 1. Clone the repository

```bash
git clone https://github.com/bsagat/marketflow
cd marketflow
```

### 2. Set configuration file

```bash
# Database configs
DB_HOST=db_host
DB_USER=Investor667
DB_PASSWORD=superpassword
DB_NAME=DBName
DB_PORT=5432

# Cache memory configs
CACHE_HOST=cache_host
CACHE_PORT=6379
CACHE_PASSWORD=superPassword
```

### 3. Load exchange images

```bash
make load-images
```

### 4. Запустите сервисы
```bash
make up
```

## API Endpoints

### 📊 Market Data API

1) `GET /prices/latest/{symbol}` – Get the latest price for a given symbol.

2) `GET /prices/latest/{exchange}/{symbol}` – Get the latest price for a given symbol from a specific exchange.

3) `GET /prices/highest/{symbol}` – Get the highest price over a period.

4) `GET /prices/highest/{exchange}/{symbol}` – Get the highest price over a period from a specific exchange.

5) `GET /prices/highest/{symbol}?period={duration}` – Get the highest price within the last {duration} (e.g., the last 1s, 3s, 5s, 10s, 30s, 1m, 3m, 5m).

6) `GET /prices/highest/{exchange}/{symbol}?period={duration}` – Get the highest price within the last {duration} from a specific exchange.

7) `GET /prices/lowest/{symbol}` – Get the lowest price over a period.

8) `GET /prices/lowest/{exchange}/{symbol}` – Get the lowest price over a period from a specific exchange.

9) `GET /prices/lowest/{symbol}?period={duration}` – Get the lowest price within the last {duration}.

10) `GET /prices/lowest/{exchange}/{symbol}?period={duration}` – Get the lowest price within the last {duration} from a specific exchange.

11) `GET /prices/average/{symbol}` – Get the average price over a period.

12) `GET /prices/average/{exchange}/{symbol}` – Get the average price over a period from a specific exchange.

13) `GET /prices/average/{exchange}/{symbol}?period={duration}` – Get the average price within the last {duration} from a specific exchange

### 🔀 Data Mode API

1) `POST /mode/test` – Switch to Test Mode (use generated data).

2) `POST /mode/live` – Switch to Live Mode (fetch data from provided programs).

### ❤️ System Health API
 
1) `GET /health` - Returns system status (e.g., connections, Redis availability)

---

## 📌 Supported Symbols & Exchanges

### 💱 Available Symbols
| Symbol     | Description                  |
|------------|------------------------------|
| `BTCUSDT`  | Bitcoin                      |
| `DOGEUSDT` | Dogecoin                     |
| `TONUSDT`  | Toncoin                      |
| `SOLUSDT`  | Solana                       |
| `ETHUSDT`  | Ethereum                     |

### 🏦 Available Exchanges
| Exchange    | Description          |
|-------------|----------------------|
| `exchange1` | Exchange Simulator 1 |
| `exchange2` | Exchange Simulator 2 |
| `exchange3` | Exchange Simulator 3 |
