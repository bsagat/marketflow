package domain

import "errors"

var (
	ErrInvalidExchangeVal = errors.New("exchange value is invalid , must be (Exchange1, Exchange2, Exchange3, All)")
	ErrInvalidMetricVal   = errors.New("metric value is invalid , must be (highest, lowest, latest, average)")
	ErrInvalidSymbolVal   = errors.New("symbol value is invalid , must be (BTCUSDT, DOGEUSDT, TONUSDT, ETHUSDT, SOLUSDT)")
	ErrInvalidModeVal     = errors.New("mode value is invalid, must be (test or live)")
	ErrEmptyMetricVal     = errors.New("metric value is empty")
	ErrEmptyExchangeVal   = errors.New("exchange value is empty")
	ErrEmptySymbolVal     = errors.New("symbol value is empty")
)
