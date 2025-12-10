package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	errorCount := 0
	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}
		resp.Body.Close()

		parts := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(parts) != 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			time.Sleep(1 * time.Second)
			continue
		}

		errorCount = 0

		// Load Average
		loadAverage, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		if loadAverage > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", loadAverage) // Изменено с %.2f на %.0f
		}

		// Memory
		totalMem, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}
		usedMem, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			continue
		}
		if totalMem > 0 {
			memUsage := float64(usedMem) / float64(totalMem) * 100
			if memUsage > 80 {
				fmt.Printf("Memory usage too high: %.0f%%\n", memUsage)
			}
		}

		// Disk
		totalDisk, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			continue
		}
		usedDisk, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			continue
		}
		if totalDisk > 0 {
			diskUsage := float64(usedDisk) / float64(totalDisk) * 100
			if diskUsage > 90 {
				// Исправлен расчет: делим на 1024 только один раз для перевода килобайт в мегабайты
				freeDiskMB := (totalDisk - usedDisk) / 1024
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
			}
		}

		// Network
		totalNet, err := strconv.ParseInt(parts[5], 10, 64)
		if err != nil {
			continue
		}
		usedNet, err := strconv.ParseInt(parts[6], 10, 64)
		if err != nil {
			continue
		}
		if totalNet > 0 {
			netUsage := float64(usedNet) / float64(totalNet) * 100
			if netUsage > 90 {
				// Исправлен расчет: делим на 1000 для перевода килобит в мегабиты
				availableNetMbit := float64(totalNet-usedNet) / 1000
				fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", availableNetMbit) // Изменено на %.0f
			}
		}

		time.Sleep(1 * time.Second)
	}
}
