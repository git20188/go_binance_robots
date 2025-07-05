package feature

import (
	"fmt"
	"go_binance_futures/feature/api/binance"
	"go_binance_futures/feature/strategy"
	"go_binance_futures/feature/strategy/line"
	"go_binance_futures/lang"
	"go_binance_futures/models"
	"go_binance_futures/notify"
	"go_binance_futures/utils"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/expr-lang/expr"
)

// 测试交易逻辑
func GoTestFeature() {
	coins, _ := GetAllSymbols()
	for _, coin := range coins {
		if coin.Symbol != "FLMUSDT" {
			continue
		}
		// positionSideLong := "LONG"
		// positionSideShort := "SHORT"
		symbol := coin.Symbol
		tickSize := coin.TickSize // 交易金额精度
		stepSize := coin.StepSize // 交易数量精度
		usdt_float64, _ := strconv.ParseFloat(coin.Usdt, 64) // 交易金额
		leverage_float64 := float64(coin.Leverage) // 合约倍数
		buyPrice, _, err := binance.GetDepthAvgPrice(symbol) // 平均买价
		if err == nil {
			buyPrice = utils.GetTradePrecision(buyPrice, tickSize) // 合理精度的价格
			quantity := (usdt_float64 / buyPrice) * leverage_float64  // 购买数量
			quantity = utils.GetTradePrecision(quantity, stepSize) // 合理精度的价格
			logs.Info(fmt.Sprintf("buyPrice:%f, quantity:%f", buyPrice, quantity), buyPrice, quantity)
		}
	}
}

func GoTestLine() {
	// logs.Info(line.BaseCheckCanLongOrShort())
	// return
	logs.Info("start test line")
	systemConfig, err := utils.GetSystemConfig()
	globalLineStrategy := GetLineStrategy(systemConfig.FutureStrategyTrade) // 选币策略
	if err != nil {
		logs.Error("GetSystemConfig:", err)
		return
	}
		
	coins, _ := GetAllSymbols()
	for _, coin := range coins {
		symbol := coin.Symbol
		if symbol != "BTCUSDT" {
			continue
		}
		
		// interval := "4h"
		// limit := 150
		// lines, _ := binance.GetKlineData(symbol, interval, limit)
		// closePrices := line.GetLineClosePrices(lines)
		
		// res, _ := line.CalculateRSI(closePrices, 6)
		// logs.Info(res)
		
		// high, low, close, _ := line.GetLineFloatPrices(lines)
		// logs.Info(high[0], low[0], close[0])
		
		// ma50, _ := line.CalculateSimpleMovingAverage(closePrices, 50)
		// logs.Info(ma50)
		// ema50, _ := line.CalculateExponentialMovingAverage(close, 50)
		// logs.Info(ema50)
		
		// upper1, ma1, lower1 := line.CalculateKeltnerChannels(high, low, close, 50, 2.75)
		// upper2, _, lower2 := line.CalculateKeltnerChannels(high, low, close, 50, 3.75)
		// // logs.Info(upper[0], ma[0], lower[0])
		
		// up, mb, dn, _ := line.CalculateBollingerBands(closePrices, 21, 2.0)
		// logs.Info(up[0], mb[0], dn[0])
		
		coin_line_strategy := globalLineStrategy
		if coin.StrategyType != "global" {
			coin_line_strategy = GetLineStrategy(coin.StrategyType)
		}
		openResult := coin_line_strategy.GetCanLongOrShort(strategy.OpenParams{
			Symbols: coin,
		})
		// logs.Info(symbol, canLang, canShort)
		if openResult.CanLong || openResult.CanShort {
			logs.Info(symbol, openResult.CanLong, openResult.CanShort)
		}
		// logs.Info("count:", index + 1)
	}
	
	logs.Info("end test line")
}

func GoTestOrder() {
	coins, _ := GetAllSymbols()
	for _, coin := range coins {
		if coin.Symbol != "CKBUSDT" {
			continue
		}
		// positionSideLong := "LONG"
		// positionSideShort := "SHORT"
		symbol := coin.Symbol
		tickSize := coin.TickSize // 交易金额精度
		stepSize := coin.StepSize // 交易数量精度
		usdt_float64, _ := strconv.ParseFloat(coin.Usdt, 64) // 交易金额
		leverage_float64 := float64(coin.Leverage) // 合约倍数
		buyPrice, sellPrice, err := binance.GetDepthAvgPrice(symbol) // 平均买价
		logs.Info(symbol, usdt_float64, leverage_float64, buyPrice, sellPrice)
		if err == nil {
			// 开多
			// buyPrice = utils.GetTradePrecision(buyPrice, tickSize) // 合理精度的价格
			// quantity := (usdt_float64 / buyPrice) * leverage_float64  // 购买数量
			// quantity = utils.GetTradePrecision(quantity, stepSize) // 合理精度的价格
			// result, _ := binance.BuyLimit(symbol, quantity, buyPrice, futures.PositionSideTypeLong)
			// logs.Info(result)
			
			// 开空
			sellPrice = utils.GetTradePrecision(sellPrice, tickSize) // 合理精度的价格
			quantity := (usdt_float64 / sellPrice) * leverage_float64  // 购买数量
			quantity = utils.GetTradePrecision(quantity, stepSize) // 合理精度的价格
			result, _ := binance.SellLimit(symbol, quantity, sellPrice, futures.PositionSideTypeShort)
			logs.Info(result)
		}
	}
}

func GoTestUtil() {
	ma1 := []float64{10, 9, 9.2, 6}
	ma2 := []float64{9, 7, 9, 6.1}
	logs.Info(line.Kdj(ma1, ma2, 4))
}

func GoTestApi() {
	res, _ := binance.GetFundingRateHistory(binance.FundingRateParams{
		Limit: 1000,
		// StartTime: (time.Now().Unix() - 60 * 60 * 12) * 1000,
	})
	logs.Info(utils.ToJson(res))
	logs.Info((time.Now().Unix() - 60 * 60))
}

func GoTestNotify() {
	// var langText, _ = lang.ReadLangJsonFile("") 
	// logs.Info(utils.ToJson(langText))
	// str := lang.Lang("futures.notice_title")
	// str = str + "1213"
	// logs.Info(str)
	// pusher.FuturesOpenOrder(notify.FuturesOrderParams{
	// 	Title: lang.Lang("futures.open_notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	PositionSide: "long",
	// 	Price: 66666,
	// 	Quantity: 1,
	//  Leverage: 10,
	// 	Status: "success",
	// })
	// pusher.FuturesOpenOrder(notify.FuturesOrderParams{
	// 	Title: lang.Lang("futures.open_notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	PositionSide: "short",
	// 	Price: 66666,
	// 	Quantity: 1,
	//  Leverage: 10,
	// 	Status: "fail",
	// 	Error: "error message error message error message error message",
	// })
	// pusher.FuturesCloseOrder(notify.FuturesOrderParams{
	// 	Title: lang.Lang("futures.close_notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	PositionSide: "short",
	// 	Price: 66666,
	// 	Quantity: 1,
	//  Leverage: 10,
	// 	Profit: 8.2,
	// 	Remarks: lang.Lang("futures.wind_of_change"),
	// 	Status: "success",
	// })
	// pusher.FuturesCloseOrder(notify.FuturesOrderParams{
	// 	Title: lang.Lang("futures.close_notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "sell",
	// 	PositionSide: "long",
	// 	Quantity: 1,
	//  Leverage: 10,
	// 	Profit: 8.2,
	// 	Remarks: "",
	// 	Status: "fail",
	// 	Error: "error message error message error message error message",
	// })
	// pusher.FuturesNotice(notify.FuturesNoticeParams{
	// 	Title: lang.Lang("futures.notice_price_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "sell",
	// 	PositionSide: "long",
	// 	Price: 66666,
	// 	AutoOrder: lang.Lang("futures.yes"),
	// })
	// pusher.FuturesNotice(notify.FuturesNoticeParams{
	// 	Title: lang.Lang("futures.notice_price_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	PositionSide: "long",
	// 	Price: 66666,
	// 	AutoOrder: lang.Lang("futures.no"),
	// })
	// pusher.FuturesListenKlineBase(notify.FuturesListenParams{
	// 	Title: lang.Lang("futures.listen_kline_base_title"),
	// 	Symbol: "BTCUSDT",
	// 	ChangePercent: 0.01,
	// 	Price: 66666,
	// 	Remarks: lang.Lang("futures.fast_up"),
	// })
	// pusher.FuturesListenKlineBase(notify.FuturesListenParams{
	// 	Title: lang.Lang("futures.listen_kline_base_title"),
	// 	Symbol: "BTCUSDT",
	// 	ChangePercent: 0.05,
	// 	Price: 66666,
	// 	Remarks: lang.Lang("futures.fast_down"),
	// })
	// pusher.FuturesListenKlineKc(notify.FuturesListenParams{
	// 	Title: lang.Lang("futures.listen_keltner_channels_title"),
	// 	PositionSide: "long",
	// 	Symbol: "BTCUSDT",
	// 	NowPrice: 66666,
	// 	StopLossPrice: 66000,
	// 	TargetHalfProfitPrice: 67000,
	// 	TargetAllProfitPrice: 68000,
	// 	DesiredPrice: 69000,
	// })
	// pusher.FuturesListenFundingRate(notify.FuturesListenParams{
	// 	Title: lang.Lang("futures.listen_funding_rate_title"),
	// 	Symbol: "BTCUSDT",
	// 	PositionSide: "long",
	// 	FundingRate: 0.12,
	// 	Price: 66666,
	// 	Remarks: lang.Lang("futures.profit_by_funding_rate"),
	// })
	
	// pusher.SpotOrder(notify.SpotOrderParams{
	// 	Title: lang.Lang("spot.notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	Price: 66666,
	// 	Quantity: 1,
	// 	Remarks: lang.Lang("spot.notice_auto_order"),
	// 	Status: "success",
	// 	Error: "",
	// })
	// pusher.SpotOrder(notify.SpotOrderParams{
	// 	Title: lang.Lang("spot.new_coin_rush_notice_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	Price: 66666,
	// 	Quantity: 1,
	// 	Remarks: lang.Lang("spot.new_coin_rush_buy"),
	// 	Status: "success",
	// 	Error: "",
	// })
	// pusher.SpotNotice(notify.SpotNoticeParams{
	// 	Title: lang.Lang("spot.notice_price_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "buy",
	// 	Price: 66666,
	// 	AutoOrder: "yes",
	// })
	// pusher.SpotNotice(notify.SpotNoticeParams{
	// 	Title: lang.Lang("spot.notice_price_title"),
	// 	Symbol: "BTCUSDT",
	// 	Side: "sell",
	// 	Price: 66666,
	// 	AutoOrder: "no",
	// })
	// pusher.SpotListenKlineBase(notify.SpotListenParams{
	// 	Title: lang.Lang("spot.kline_listen_base_title"),
	// 	Symbol: "BTCUSDT",
	// 	ChangePercent: 0.05,
	// 	Price: 66666,
	// 	Remarks: lang.Lang("spot.fast_down"),
	// })
	pusher.SpotListenKlineBase(notify.SpotListenParams{
		Title: lang.Lang("spot.kline_listen_base_title"),
		Symbol: "BTCUSDT",
		ChangePercent: 0.05,
		Price: 66666,
		Remarks: lang.Lang("spot.fast_up"),
	})
}

func GoTestListen() {
	o := orm.NewOrm()
	var coins []models.ListenSymbols
	o.QueryTable("listen_symbols").OrderBy("ID").Filter("enable", 1).Filter("type", 2).All(&coins) // 通知币列表
	
	for _, coin := range coins {
		if coin.Symbol != "BTCUSDT" {
			continue
		}
		logs.Info("listen futures: %s, type: %s ", coin.Symbol, coin.ListenType)
		switch coin.ListenType {
			case "kline_base":
				klineBaseListen(coin)
			case "kline_kc":
				klineKcListen(coin)
			case "custom":
				klineCustomListen(coin)
			default:
				logs.Error("listen type error:", coin.ListenType)
		}
	}
}

func GoTestParse() {
	
	env := map[string]interface{}{
		"ma":   [][]float64{{1, 2, 3}, {4, 5, 6}},
	}
	// https://expr-lang.org/docs/language-definition#variables
	// 可以添加注释，使用 // 或者 /* */
	// 可以写多行的语句
	// 可以定义变量 let x = 1
	// 可以调用(内置和自定义)函数
	// 可以调用(自定义)变量
	code := `
		// 买入条件
		let x = ma[0][0] + ma[1][0];
		x + 2
	`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		logs.Error(err.Error())
	}
	output, err := expr.Run(program, env)
	if err != nil {
		logs.Error(err.Error())
	}
	fmt.Println(output)
}
