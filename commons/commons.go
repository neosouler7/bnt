package commons

import (
	"bnt/config"
	"bnt/tgmanager"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// // TODO. btc market에서의 eth를 볼 때... 사토시 단위가 짤린다.
// // return 값 먼저 확인하고, 하지만 의심가는 부분은 여기. parsing하면서 유실되는 것으로 추정. 10^-8까지 필요.
// func GetObTargetPrice(volume string, orderbook interface{}) string {
// 	/*
// 		ask's price should go up, and bid should go down

// 		ask = [[p1, v1], [p2, v2], [p3, v3] ...]
// 		bid = [[p3, v3], [p2, v2], [p1, p1] ...]
// 	*/
// 	currentVolume := 0.0
// 	targetVolume, err := strconv.ParseFloat(volume, 64)
// 	tgmanager.HandleErr("GetObTargetPrice1", err)

// 	obSlice := orderbook.([]interface{})
// 	for _, ob := range obSlice {
// 		obInfo := ob.([2]string)
// 		volume, err := strconv.ParseFloat(obInfo[1], 64)
// 		tgmanager.HandleErr("GetObTargetPrice2", err)

// 		currentVolume += volume
// 		if currentVolume >= targetVolume {
// 			return obInfo[0]
// 		}
// 	}
// 	return obSlice[len(obSlice)-1].([2]string)[0]
// }

// func GetTargetVolumeMap(exchange string) map[string]string {
// 	pairs := config.GetPairs(exchange)
// 	m := make(map[string]string, len(pairs)) // 초기 용량 설정

// 	for _, p := range pairs {
// 		idx1 := strings.Index(p, ":")
// 		idx2 := strings.LastIndex(p, ":")

// 		if idx1 < 0 || idx2 <= idx1 {
// 			log.Printf("Invalid pair format: %s", p)
// 			continue
// 		}

// 		market := p[:idx1]
// 		symbol := p[idx1+1 : idx2]
// 		targetVolume := p[idx2+1:]
// 		m[market+":"+symbol] = targetVolume
// 	}
// 	return m
// }

func GetPairMap(exchange string) map[string]interface{} {
	pairs := config.GetPairs(exchange)
	m := make(map[string]interface{}, len(pairs)) // 초기 용량 설정

	for _, pair := range pairs {
		idx := strings.Index(pair, ":")
		if idx < 0 {
			log.Printf("Invalid pair format: %s", pair)
			continue
		}

		market := pair[:idx]
		symbol := pair[idx+1:]
		m[symbol+market] = map[string]string{"market": market, "symbol": symbol}
	}
	return m
}

func FormatTs(ts string) string {
	tsLen := len(ts)

	if tsLen < 13 {
		var sb strings.Builder
		sb.WriteString(ts)
		sb.WriteString(strings.Repeat("0", 13-tsLen))
		return sb.String()
	} else if tsLen == 13 { // if millisecond
		return ts
	} else {
		return ts[:13]
	}
}

// func Min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

func Bytes2Json(data []byte, i interface{}) {
	r := bytes.NewReader(data)
	err := json.NewDecoder(r).Decode(i)
	tgmanager.HandleErr("Bytes2Json", err)
}

func SetTimeZone(name string) *time.Location {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Seoul"
		fmt.Printf("%s : DEFAULT %s\n", name, tz)
	} else {
		fmt.Printf("%s : SERVER %s\n", name, tz)
	}
	location, _ := time.LoadLocation(tz)
	return location
}
