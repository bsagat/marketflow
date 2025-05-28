package datafetcher

import (
	"marketflow/internal/domain"
)

type TestMode struct {
	RawDataCh chan []domain.Data // канал где будет храниться все цены за последние 60 секунд
}

var _ domain.DataFetcher = (*TestMode)(nil)

func NewTestModeFetcher() *TestMode {
	return &TestMode{}
}

func (m *TestMode) SetupDataFetcher() (chan map[string]domain.ExchangeData, chan []domain.Data) {
	ch := make(chan map[string]domain.ExchangeData)
	rawDataCh := make(chan []domain.Data)

	// close(ch)
	// close(rawDataCh) // обязательно
	return ch, rawDataCh
}

func (m *TestMode) CheckHealth() error {
	return nil
}

func (m *TestMode) Close() {

}
