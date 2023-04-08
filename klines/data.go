package klines

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type KLineEntry struct {
	OpenTime                 int64   `db:"open_time" json:"open_time"`
	OpenPrice                float64 `db:"open_price" json:"open_price"`
	LowPrice                 float64 `db:"low_price" json:"low_price"`
	HighPrice                float64 `db:"high_price" json:"high_price"`
	ClosePrice               float64 `db:"close_price" json:"close_price"`
	Volume                   float64 `db:"volume" json:"volume"`
	CloseTime                int64   `db:"close_time" json:"close_time"`
	QuoteAssetVolume         float64 `db:"quote_asset_volume" json:"quote_asset_volume"`
	NumTrades                int64   `db:"num_trades" json:"num_trades"`
	TakerBuyBaseAssetVolume  float64 `db:"taker_buy_base_asset_volume" json:"taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume float64 `db:"taker_buy_quote_asset_volume" json:"taker_buy_quote_asset_volume"`
}
type KLineData struct {
	Symbol string
	Period time.Duration
	Data   []KLineEntry
}

func NewKLineData(symbol string, period time.Duration) *KLineData {
	return &KLineData{
		Symbol: symbol,
		Period: period,
		Data:   []KLineEntry{},
	}
}

func FetchKLineDataFromDB(db *sqlx.DB, symbol string, startTime int64, endTime int64) (*KLineData, error) {
	kd := NewKLineData(symbol, time.Minute)
	var data []KLineEntry
	query := `SELECT open_time, open_price, high_price, low_price, close_price, volume, close_time, quote_asset_volume, num_trades, taker_buy_base_asset_volume, taker_buy_quote_asset_volume 
	FROM binance.klines 
	WHERE symbol = $1 AND period = $2 AND open_time BETWEEN $3 AND $4 
	ORDER BY open_time ASC`
	err := db.Select(&data, query, kd.Symbol, "1m", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kline data from database: %w", err)
	}
	kd.Data = data
	return kd, nil
}

func FetchKLineDataFromDBSincePeriodsBefore(db *sqlx.DB, symbol string, minutesBefore int64) (*KLineData, error) {
	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	startTime := endTime - minutesBefore*60000
	return FetchKLineDataFromDB(db, symbol, startTime, endTime)
}
