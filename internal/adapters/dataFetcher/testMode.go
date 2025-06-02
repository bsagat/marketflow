package datafetcher

import (
	"marketflow/internal/domain"
	"math/rand"
	"time"
)

type TestMode struct {
	stop chan struct{}
}

var _ domain.DataFetcher = (*TestMode)(nil)

func NewTestModeFetcher() *TestMode {
	return &TestMode{stop: make(chan struct{})}
}

func (m *TestMode) SetupDataFetcher() (chan map[string]domain.ExchangeData, chan []domain.Data, error) {
	ch := make(chan map[string]domain.ExchangeData, 100)
	ch2 := make(chan []domain.Data, 100)
	pairs := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}
	exchanges := []string{"Exchange1", "Exchange2", "Exchange3"}
	basePrices := map[string]float64{
		"BTCUSDT": 60000.0, "DOGEUSDT": 0.15, "TONUSDT": 5.0, "SOLUSDT": 150.0, "ETHUSDT": 3000.0,
	}

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-m.stop:
				close(ch)
				close(ch2)
				// close(m.messageChan)
				return
			case <-ticker.C:
				data := make(map[string]domain.ExchangeData)
				rawData := make([]domain.Data, 0, len(pairs)*len(exchanges))
				for _, ex := range exchanges {
					for _, pair := range pairs {
						// Generate random price fluctuation (±15%)
						price := basePrices[pair] * (1 + (rand.Float64()-0.5)*0.3)
						data[ex+"-"+pair] = domain.ExchangeData{
							Pair_name:     pair,
							Exchange:      ex,
							Timestamp:     time.Now(),
							Average_price: price, // Use Average_price as per your struct
							Min_price:     price, // Set Min/Max same for real-time update
							Max_price:     price,
						}
						rawData = append(rawData, domain.Data{
							ExchangeName: ex,
							Symbol:       pair,
							Price:        price,
							Timestamp:    time.Now().UnixNano() / int64(time.Millisecond),
						})
					}
				}
				ch <- data
				ch2 <- rawData

			}
		}
	}()

	// close(ch)
	return ch, ch2, nil
}

func (m *TestMode) CheckHealth() error {
	return nil // is it nado voobshe?
}

func (m *TestMode) Close() {
	close(m.stop)
}
