# Marketflow

### Практическая часть 
1) Построить архитектуру приложения (hexagonal architecture)
2) Реализация контейнера (healthcheck, graceful shutdown)
    Configuration file добавить
3) Построить ERD для БД шки
4) Настроить подключение к external adapters внутри приложения(redis, real data processing, postgres) 
    Учесть момент с failover reconnect
5) Описать домены, сущности в бизнес логике, типо;
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

6) Реализовать дата парсинг (из provided programs) "думаю самое хардовое"
7) Реализовать API endpoint-ы 
8) Написать help функцию 🗿
9) Тестирование 


### Теоритическая часть 
Изучить паттерны конкурентности
Узнать как взаимодействовать с redis (и зач он вообще здесь нужен)


### Дополнительно
Используем slog для логирования (ВАЖНО: добавляем контекстуальную информацию для лучшей откладки)
Документация кода (комментарий, инструкции к сущностям кода)

### Optional 
