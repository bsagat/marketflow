package datafetcher

import "marketflow/internal/domain"

type TestMode struct{}

var _ domain.DataFetcher = (*TestMode)(nil)

func (m *TestMode) SetupDataFetcher() chan map[string]domain.ExchangeData {
	ch := make(chan map[string]domain.ExchangeData)

	return ch
}

func (m *TestMode) CheckHealth() error {
	return nil
}

func (m *TestMode) Close() {

}
