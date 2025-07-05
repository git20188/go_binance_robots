package line

import (
	"go_binance_futures/models"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/beego/beego/v2/client/orm"
)

type Line struct {
	Position string
	High float64
	Low float64
	Close float64
	Open float64
	TradeNum int64
}
type LineData struct {
	MaxIndex int
	MinIndex int
	Line []*Line
}

// 归一化处理k线数据
func normalizationLineData(data []*futures.Kline) (*LineData) {
	maxIndex := 0
	maxPrice := 0.0
	minIndex := 0
	minPrice := 0.0
	line := make([]*Line, len(data))
	for key, item := range data {
		open, _ := strconv.ParseFloat(item.Open, 64)
		high, _ := strconv.ParseFloat(item.High, 64)
		low, _ := strconv.ParseFloat(item.Low, 64)
		close, _ := strconv.ParseFloat(item.Close, 64)
		if key == 0 {
			maxPrice = high
			minPrice = close
		} else {
			if (high > maxPrice) {
				maxPrice = high
				maxIndex = key
			}
			if (low < minPrice) {
				minPrice = low
				minIndex = key
			}
		}
		position := "LONG"
		if close < open {
			position = "SHORT"
		}
		line[key] = &Line{
			Position: position,
			High: high,
			Low: low,
			Close: close,
			Open: open,
			TradeNum: item.TradeNum,
		}
	}
	return &LineData{
		MaxIndex: maxIndex,
		MinIndex: minIndex,
		Line: line,
	}
}

// 获取收盘价列表
func GetClosePrices(data []*Line) ([]float64) {
	clonePrices := make([]float64, len(data))
	for key, line := range data {
		clonePrices[key] = line.Close
	}
	return clonePrices
}

// 从k线获取收盘价列表
func GetLineClosePrices(data []*futures.Kline) ([]float64) {
	clonePrices := make([]float64, len(data))
	for key, item := range data {
		close, _ := strconv.ParseFloat(item.Close, 64)
		clonePrices[key] = close
	}
	return clonePrices
}

// 从k线获取最高价、最低价、收盘价列表
func GetLineFloatPrices(data []*futures.Kline) (high, low, close, open []float64) {
	high = make([]float64, len(data))
	low = make([]float64, len(data))
	close = make([]float64, len(data))
	open = make([]float64, len(data))
	for key, item := range data {
		highPrice, _ := strconv.ParseFloat(item.High, 64)
		lowPrice, _ := strconv.ParseFloat(item.Low, 64)
		closePrice, _ := strconv.ParseFloat(item.Close, 64)
		openPrice, _ := strconv.ParseFloat(item.Open, 64)
		high[key] = highPrice
		low[key] = lowPrice
		close[key] = closePrice
		open[key] = openPrice
	}
	return high, low, close, open
}

func GetLineFloatValues(data []*futures.Kline) (high, low, close, open, amount, qps []float64) {
	high = make([]float64, len(data))
	low = make([]float64, len(data))
	close = make([]float64, len(data))
	open = make([]float64, len(data))
	amount = make([]float64, len(data))
	qps = make([]float64, len(data))
	for key, item := range data {
		highPrice, _ := strconv.ParseFloat(item.High, 64)
		lowPrice, _ := strconv.ParseFloat(item.Low, 64)
		closePrice, _ := strconv.ParseFloat(item.Close, 64)
		openPrice, _ := strconv.ParseFloat(item.Open, 64)
		amountFloat, _ := strconv.ParseFloat(item.QuoteAssetVolume, 64)
		high[key] = highPrice
		low[key] = lowPrice
		close[key] = closePrice
		open[key] = openPrice	
		amount[key] = amountFloat
		qps[key] = amountFloat
		if item.CloseTime - item.OpenTime > 0 {
			qps[key] = amountFloat / float64(item.CloseTime - item.OpenTime)
		}
	}
	return high, low, close, open, qps, amount
}

// 获取某中类型的line的数量是否超过阈值
func getRightLine(data []*Line, position string) bool {
	positionCount := 0
	for _, line := range data {
		if line.Position == position {
			positionCount++
		}
	}
	return len(data) - positionCount <= 2
}

// 获取所有交易的币
func GetAllSymbols() (symbols []*models.Symbols, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("symbols").OrderBy("ID").All(&symbols)
	return symbols, err
}

// 判断 btc 的涨跌是否大于 5%，判断 当前所有币种涨跌数量是否 80%，是否是一个单项行情
func BaseCheckCanLongOrShort() (canLong bool, canShort bool) {
	coins, err := GetAllSymbols()
	if err != nil {
		return false, false
	}
	canLong, canShort = true, true
	riseCount, fallCount := 0, 0
	btcPercentChange := 0.0
	for _, coin := range coins {
		if coin.PercentChange >= 0 {
			riseCount++
		} else {
			fallCount++
		}
		if coin.Symbol == "BTCUSDT" {
			btcPercentChange = coin.PercentChange
		}
	}
	// logs.Info(riseCount, fallCount, btcPercentChange, len(coins))
	if riseCount / len(coins) > 75 {
		// 都在涨，不要做空
		canShort = false
	}
	if fallCount / len(coins) > 75 {
		// 都在跌，不要做多
		canLong = false
	}
	if riseCount / len(coins) > 60 && btcPercentChange > 5 {
		// 60% 的币种都在涨，btc 涨幅大于 5，不要做空
		canShort = false
	}
	if fallCount / len(coins) < 60 && btcPercentChange < -5 {
		// 60% 的币种都在跌，btc 跌幅大于 5，不要做多
		canLong = false
	}
	return canLong, canShort
}

func BaseTrend() float64 {
	o := orm.NewOrm()
	var symbols []models.Symbols
	sql := "SELECT * FROM symbols WHERE symbol = ? OR symbol = ? OR symbol = ? OR symbol = ?"
	o.Raw(sql, "BTCUSDT", "ETHUSDT", "SOLUSDT", "BNBUSDT").QueryRows(&symbols)
	
	basicTrend := 0.0
	for _, v := range symbols {
		if v.Symbol == "BTCUSDT" {
			basicTrend += v.PercentChange * 0.6
		} else if v.Symbol == "ETHUSDT" {
			basicTrend += v.PercentChange * 0.3
		} else if v.Symbol == "SOLUSDT" {
			basicTrend += v.PercentChange * 0.05
		} else if v.Symbol == "BNBUSDT" {
			basicTrend += v.PercentChange * 0.05
		}
	}
	return basicTrend
}
