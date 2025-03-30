package filemanager

import (
	"bnt/tgmanager"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileManager struct {
	mu        sync.Mutex
	baseDir   string
	orderbook map[string][][]string
	trade     map[string][][]string
}

var FM *FileManager = NewFileManager(".")

func NewFileManager(baseDir string) *FileManager {
	fm := &FileManager{
		baseDir:   baseDir,
		orderbook: make(map[string][][]string),
		trade:     make(map[string][][]string),
	}
	go fm.startDumpRoutine()
	return fm
}

func (fm *FileManager) startDumpRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		fm.dumpData()
	}
}

func (fm *FileManager) dumpData() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	currentDate := time.Now().Format("060102")
	currentMinute := time.Now().Format("1504")
	minute := time.Now().Minute()

	for key, records := range fm.orderbook {
		dirPath := fm.getDirPath(key, "orderbook", currentDate)
		filePath := fm.getFilePath(dirPath, currentMinute)
		fm.writeCSV(filePath, records)
	}
	fmt.Printf("[ob] dump %s %s\n", currentDate, currentMinute)

	for key, records := range fm.trade {
		dirPath := fm.getDirPath(key, "trade", currentDate)
		filePath := fm.getFilePath(dirPath, currentMinute)
		fm.writeCSV(filePath, records)
	}
	fmt.Printf("[tr] dump %s %s\n", currentDate, currentMinute)

	fm.orderbook = make(map[string][][]string)
	fm.trade = make(map[string][][]string)

	if minute%10 == 0 {
		tgMsg := fmt.Sprintf("FileManager Dump completed at %s", currentMinute)
		tgmanager.SendMsg(tgMsg)
	}
}

func (fm *FileManager) writeCSV(filePath string, records [][]string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Printf("Error writing to CSV %s: %v\n", filePath, err)
		}
	}
}

func (fm *FileManager) getDirPath(key, fileType, date string) string {
	var baseDir string
	switch fileType {
	case "orderbook":
		baseDir = "csv_orderbook"
	case "trade":
		baseDir = "csv_trade"
	}
	dirPath := fmt.Sprintf("%s/%s/%s/%s", fm.baseDir, baseDir, key, date)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
	}
	return dirPath
}

func (fm *FileManager) getFilePath(dirPath, currentMinute string) string {
	return fmt.Sprintf("%s/%s.csv", dirPath, currentMinute)
}

func (fm *FileManager) flattenOrderbookSlice(slice []interface{}, limit int) []string {
	var result []string
	count := 0

	for _, item := range slice {
		price := item.([2]string)[0]
		volume := item.([2]string)[1]
		result = append(result, price, volume)
		count++

		if count >= limit {
			break
		}
	}
	return result
}
func (fm *FileManager) PreHandleOrderbook(exchange, market, symbol, ts string, askSlice, bidSlice []interface{}) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	key := fmt.Sprintf("%s/%s/%s", exchange, market, symbol)

	flattenAsks := fm.flattenOrderbookSlice(askSlice, 10)
	flattenBids := fm.flattenOrderbookSlice(bidSlice, 10)

	record := append([]string{ts}, append(flattenAsks, flattenBids...)...)
	fm.orderbook[key] = append(fm.orderbook[key], record)
}

func (fm *FileManager) PreHandleTrade(exchange, market, symbol, ts, priceTrade, tsTrade string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	key := fmt.Sprintf("%s/%s/%s", exchange, market, symbol)
	record := []string{ts, tsTrade, priceTrade}
	fm.trade[key] = append(fm.trade[key], record)
}
