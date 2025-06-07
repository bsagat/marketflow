# Marketflow

### Практическая часть 
1) Построить архитектуру приложения (hexagonal architecture) ✅
2) Реализация контейнера (healthcheck, graceful shutdown)✅
    Configuration file добавить ✅
4) Настроить подключение к external adapters внутри приложения(redis, real data processing, postgres) ✅
    Учесть момент с failover reconnect✅
5) Описать домены, сущности в бизнес логике, типо;✅
```go
type ExchangeData struct{
    Pair_name     string    // the trading pair name.
    Exchange      string    // the exchange from which the data was received.
    Timestamp     time.Time // the time when the data is stored.
    Average_price float     // the average price of the trading pair over the last minute.
    Min_price     float     // the minimum price of the trading pair over the last minute.
    Max_price     float     // the maximum price of the trading pair over the last minute
} 
```

6) Реализовать дата парсинг (из provided programs) "думаю самое хардовое" ✅
7) Реализовать API endpoint-ы 
    -  Нужно исправить поиск в dataBuffer(lowest,highest), чтобы при запросе по всем биржам (All exchanges) учитывались данные от каждой из них, а не только одной.
    - Сделать более корректные запросы, тестить их на мок данных
8) Написать help функцию 🗿✅
9) Тестирование 


### Теоритическая часть 
Изучить паттерны конкурентности
Узнать как взаимодействовать с redis (и зач он вообще здесь нужен)


### Дополнительно
Используем slog для логирования (ВАЖНО: добавляем контекстуальную информацию для лучшей откладки)
Документация кода (комментарий, инструкции к сущностям кода)

### Option: 
Market Data API

GET /prices/latest/{symbol} – Get the latest price for a given symbol.✅

GET /prices/latest/{exchange}/{symbol} – Get the latest price for a given symbol from a specific exchange.  

GET /prices/highest/{symbol} – Get the highest price over a period. ✅

GET /prices/highest/{exchange}/{symbol} – Get the highest price over a period from a specific exchange.

GET /prices/highest/{symbol}?period={duration} – Get the highest price within the last {duration} (e.g., the last 1s, 3s, 5s, 10s, 30s, 1m, 3m, 5m).

GET /prices/highest/{exchange}/{symbol}?period={duration} – Get the highest price within the last {duration} from a specific exchange.

GET /prices/lowest/{symbol} – Get the lowest price over a period.✅

GET /prices/lowest/{exchange}/{symbol} – Get the lowest price over a period from a specific exchange.

GET /prices/lowest/{symbol}?period={duration} – Get the lowest price within the last {duration}.

GET /prices/lowest/{exchange}/{symbol}?period={duration} – Get the lowest price within the last {duration} from a specific exchange.

GET /prices/average/{symbol} – Get the average price over a period. ✅

GET /prices/average/{exchange}/{symbol} – Get the average price over a period from a specific exchange.

GET /prices/average/{exchange}/{symbol}?period={duration} – Get the average price within the last {duration} from a specific exchange

TEST:
http://localhost:8080/prices/latest/BTCUSDT✅
http://localhost:8080/prices/latest/Exchange1/BTCUSDT✅

сделать чтобы Exchange правильынй показывал? Надо?
http://localhost:8080/prices/highest/BTCUSDT✅

данные есть хз праввильно?:
http://localhost:8080/prices/highest/Exchange1/BTCUSDT✅
Not work in test:
http://localhost:8080/prices/highest/BTCUSDT?period=1s 
http://localhost:8080/prices/highest/Exchange1/BTCUSDT?period=1s


http://localhost:8080/prices/lowest/BTCUSDT

not in test:
http://localhost:8080/prices/lowest/Exchange1/BTCUSDT
http://localhost:8080/prices/lowest/BTCUSDT?period=1s
http://localhost:8080/prices/lowest/Exchange1/BTCUSDT?period=1s

http://localhost:8080/prices/average/BTCUSDT
http://localhost:8080/prices/average/Exchange1/BTCUSDT
http://localhost:8080/prices/average/Exchange1/BTCUSDT?period=1s

localhost:8080/health
Check health

localhost:8080/mode/live
Change test mode to live mode

localhost:8080/prices/latest/Exchange1/BTCUSDT
latest data from specific exchange

localhost:8080/prices/latest/BTCUSDT
latest data from all exchanges


Domain -> health chek -> ConnMs?
Domain -> interfaces are these intrfcs implemeted?  
Getenv?
CacheMem -> Helth_chekc?



1. very first time: 
docker load -i build/exchange_images/exchange1_amd64.tar
docker load -i build/exchange_images/exchange2_amd64.tar
docker load -i build/exchange_images/exchange3_amd64.tar

2. docker-compose -f build/docker-compose.yml up / docker-compose -f build/docker-compose.yml up --build


3. nc 127.0.0.1 40101
    |
    ->to test

4. go run ./cmd



BTCUSDT
DOGEUSDT
TONUSDT
SOLUSDT
ETHUSDT


In test mode ticks goes much faster tahn in live. In live it is is like one per second


testcases
Проверить Датапарсинг если же редис будет отключен
Проверить запросы с таймлайном для макс и мин прайсов
Проверить все запросы если поменять режимы много раз (25 раз допустим) для двух режимов
Проверить все запросы если база данных будет пустой
Проверить на гонку данных


Проверяем все обратно когда приложение упаковано внутри докера
Проверяем все запросы если наше приложение запущено на другом порте