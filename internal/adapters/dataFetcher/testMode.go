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

func (m *TestMode) SetupDataFetcher() (chan map[string]domain.ExchangeData, chan []domain.Data, error){
	ch := make(chan map[string]domain.ExchangeData)
	ch2 := make(chan []domain.Data)
	pairs := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}
	exchanges := []string{"exchange1", "exchange2", "exchange3"}
	basePrices := map[string]float64{
		"BTCUSDT": 60000.0, "DOGEUSDT": 0.15, "TONUSDT": 5.0, "SOLUSDT": 150.0, "ETHUSDT": 3000.0,
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-m.stop:
				close(ch)
				return
			case <-ticker.C:
				data := make(map[string]domain.ExchangeData)
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
					}
				}
				ch <- data

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
