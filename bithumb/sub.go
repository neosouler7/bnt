package bithumb

import (
	"bnt/commons"
	"bnt/filemanager"
	"fmt"
	"strconv"
	"strings"
)

func getMST(rJson map[string]interface{}) (string, string, string) {
	pair := rJson["code"].(string) // only websocket
	var pairInfo = strings.Split(pair, "-")
	var market, symbol = strings.ToLower(pairInfo[0]), strings.ToLower(pairInfo[1])

	tsFloat := int(rJson["timestamp"].(float64))
	ts := commons.FormatTs(strconv.Itoa(tsFloat))

	return market, symbol, ts
}

func SetOrderbook(exchange string, rJson map[string]interface{}) {
	market, symbol, ts := getMST(rJson)
	orderbooks := rJson["orderbook_units"].([]interface{})

	var askSlice, bidSlice []interface{}
	for _, ob := range orderbooks {
		o := ob.(map[string]interface{})
		ask := [2]string{fmt.Sprintf("%f", o["ask_price"]), fmt.Sprintf("%f", o["ask_size"])}
		bid := [2]string{fmt.Sprintf("%f", o["bid_price"]), fmt.Sprintf("%f", o["bid_size"])}
		askSlice = append(askSlice, ask)
		bidSlice = append(bidSlice, bid)
	}

	filemanager.FM.PreHandleOrderbook(
		exchange,
		market,
		symbol,
		ts,
		askSlice,
		bidSlice,
	)
}

func SetTrade(exchange string, rJson map[string]interface{}) {
	market, symbol, ts := getMST(rJson)

	priceTradeFloat := int(rJson["trade_price"].(float64))
	priceTrade := commons.FormatTs(strconv.Itoa(priceTradeFloat))
	tsTradeFloat := int(rJson["trade_timestamp"].(float64))
	tsTrade := commons.FormatTs(strconv.Itoa(tsTradeFloat))

	filemanager.FM.PreHandleTrade(
		exchange,
		market,
		symbol,
		ts,
		priceTrade,
		tsTrade,
	)
}
