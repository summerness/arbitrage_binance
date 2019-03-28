package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance"
	"math"
	"strconv"
	"time"
)

var (
	apiKey    = ""
	secretKey = ""
)

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

var client = binance.NewClient(apiKey, secretKey)

func getBalance(symbol string) (qty float64) {
	res, err := client.NewGetAccountService().Do(context.Background())
	var count float64
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, each := range res.Balances {
		if each.Asset == symbol {
			count, _ = strconv.ParseFloat(each.Free, 64)
		}
	}
	return count

}

func trend() (state bool) {
	endTime := time.Now().UnixNano() / 1e6
	startTime := time.Now().Add(-time.Minute*60).UnixNano() / 1e6
	fmt.Println(endTime, startTime)
	trades, err := client.NewAggTradesService().
		Symbol("ETHBTC").StartTime(startTime).
		EndTime(endTime).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var startTemp float64 = 0
	for _, t := range trades[:10] {
		tPrice, _ := strconv.ParseFloat(t.Price, 64)
		startTemp = tPrice + startTemp
	}
	startAve := startTemp / 10

	var endTemp float64 = 0
	for _, t := range trades[len(trades)-10:] {
		tPrice, _ := strconv.ParseFloat(t.Price, 64)
		endTemp = tPrice + endTemp
	}
	endAve := endTemp / 10
	fmt.Println(startAve, endAve)

	var a Accuracy = func() float64 { return 0.0000001 }

	if a.Greater(startTemp, endAve) {
		return false
	} else {
		return true
	}

}

func trendSmall() (state int) {
	endTime := time.Now().UnixNano() / 1e6
	startTime := time.Now().Add(-time.Minute*10).UnixNano() / 1e6
	fmt.Println(endTime, startTime)
	trades, err := client.NewAggTradesService().
		Symbol("ETHBTC").StartTime(startTime).
		EndTime(endTime).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	var startTemp float64 = 0
	for _, t := range trades[:10] {
		tPrice, _ := strconv.ParseFloat(t.Price, 64)
		startTemp = tPrice + startTemp
	}
	startAve := startTemp / 10

	var end_temp float64 = 0
	for _, t := range trades[len(trades)-10:] {
		tPrice, _ := strconv.ParseFloat(t.Price, 64)
		end_temp = tPrice + end_temp
	}
	endAve := end_temp / 10
	fmt.Println(startAve, endAve)

	var a Accuracy = func() float64 { return 0.0000001 }

	if a.Greater(startAve, endAve) {
		fmt.Println(-1)
		return -1
	} else if a.Equal(startAve, endAve) {
		fmt.Println(0)
		return 0
	} else if a.Smaller(startAve, endAve) {
		fmt.Println(1)
		return 1
	} else {
		return
	}
}

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

func limitOrder(sellPer float64, buyPer float64) {
	res, err := client.NewDepthService().Symbol("ETHBTC").
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	// 买一 卖一价
	buyPrice := res.Bids[0].Price
	sellPrice := res.Bids[0].Price

	fSellPrice, _ := strconv.ParseFloat(sellPrice, 64)
	fmt.Println(fSellPrice)
	fSellPrice = Round2(fSellPrice*sellPer, 6)
	fBuyPrice, _ := strconv.ParseFloat(buyPrice, 64)
	fmt.Println(fBuyPrice)
	fBuyPrice = Round2(fBuyPrice*buyPer, 6)

	sSellPrice := strconv.FormatFloat(fSellPrice, 'g', -1, 64)
	sBuyPrice := strconv.FormatFloat(fBuyPrice, 'g', -1, 64)
	fmt.Println(sSellPrice, sBuyPrice)

	tradeSize := getBalance("ETH") * 0.98
	tradeSizeCount := Round2(tradeSize, 3)
	sTradeSize := strconv.FormatFloat(tradeSizeCount, 'g', -1, 64)
	fmt.Println(sTradeSize)

	//sell_trade_size := strconv.FormatFloat(Tradesize,'g',-1,64)
	if ((fSellPrice - fBuyPrice) > 0) && (fBuyPrice < 0.00031) {

		//endTime := time.Now().UnixNano()/1e6
		//startTime := time.Now().Add(-time.Hour*1).UnixNano()/1e6
		//trades, err := client.NewAggTradesService().
		//	Symbol("ETHBTC").StartTime(startTime).
		//	EndTime(endTime).Do(context.Background())
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//var prices []float64
		//for _, t := range trades{
		//	t_Price,_ := strconv.ParseFloat(t.Price,64)
		//	prices = append(prices,t_Price)
		//}

		//max_price , min_price := Max_min_arry(prices)
		//var a Accuracy = func() float64 { return 0.0000001 }
		//if (a.Greater(max_price*1.008,f_sell_price)) && (a.Smaller(min_price*0.992,f_buy_price)){
		_, err1 := client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceGTC).Quantity(sTradeSize).Price(sSellPrice).Do(context.Background())
		if err1 != nil {
			fmt.Println(err1)
			return
		}

		_, err2 := client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceGTC).Quantity(sTradeSize).Price(sBuyPrice).Do(context.Background())
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		//}
	}
}

func orderState() (state bool) {
	orders, err := client.NewListOpenOrdersService().Symbol("ETHBTC").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, o := range orders {
		fmt.Println(o)
		return true
	}
	return false
}

func main() {
	counts := 0
	for {
		if !orderState() {
			counts = 0
			time.Sleep(time.Hour * 5)

			if !trend() {
				limitOrder(1.001, 0.995)
			}
			//} else {
			//	limit_order(1.010,1)
			//}
		}
		counts++
		time.Sleep(time.Minute * 1)
		if counts > 5440 {
			orders, _ := client.NewListOpenOrdersService().Symbol("ETHBTC").
				Do(context.Background())
			for _, o := range orders {
				orderId := o.OrderID
				orderSide := o.Side
				tradeSize := o.OrigQuantity
				if orderSide == "SELL" {
					if trendSmall() == -1 {
						_, err := client.NewCancelOrderService().Symbol("ETHBTC").
							OrderID(orderId).Do(context.Background())
						if err == nil {
							client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeSell).
								Type(binance.OrderTypeMarket).Quantity(tradeSize).Do(context.Background())
						}
					}
				} else if tradeSize == "BUY" {
					if trendSmall() == 1 {
						_, err := client.NewCancelOrderService().Symbol("ETHBTC").
							OrderID(orderId).Do(context.Background())
						if err == nil {
							client.NewCreateOrderService().Symbol("ETHBTC").Side(binance.SideTypeBuy).
								Type(binance.OrderTypeMarket).Quantity(tradeSize).Do(context.Background())
						}
					}
				}
			}
		}
	}
}