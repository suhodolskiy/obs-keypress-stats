package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func PersistCount(name string, count int) {
	log.Printf("Data saved to file  [%d key presses] 💾\n", count)

	data := fmt.Sprintf("%s %d", time.Now().Format(time.RFC3339), count)
	os.WriteFile(name, []byte(data), 0644)
}

func ReadCount(name string) int {
	data, err := os.ReadFile(name)
	if err != nil {
		return 0
	}

	text := string(data)
	i := strings.LastIndex(text, " ")
	if i == -1 {
		return 0
	}

	datetimeStr := text[:i]
	num, err := strconv.Atoi(text[i+1:])
	if err != nil {
		return 0
	}

	date, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return 0
	}

	if diff := time.Now().Sub(date); diff > time.Hour || diff < 0 {
		return 0
	}

	log.Println("Data loaded from file 📖")

	return num
}
