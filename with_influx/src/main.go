package main

import (
	"fmt"
	"sort"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SurroundingsPayload struct {
	Number      int       `json:"number"`
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"Temperature"`
	Moisture    float64   `josn:"Moisture"`
	AirPressure float64   `json:"airPressure"`
}

func insertPayload(payload []SurroundingsPayload) {
	var token string
	var bucket string
	var org string
	var dbUrl string

	token = "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNg"
	bucket = "vuoy_monitor"
	org = "iot"
	dbUrl = "http://influxdb:8086"

	fmt.Printf("connectingt to %s , bucket :%s ,org :%s ,token :%s\n", dbUrl, bucket, org, token)
	client := influxdb2.NewClient(dbUrl, token)
	writeAPI := client.WriteAPI(org, bucket)

	// Add this block to listen for errors from the writeAPI
	go func() {
		for err := range writeAPI.Errors() {
			fmt.Println("Error writing to InfluxDB:", err)
		}
	}()
	sort.Slice(payload, func(i, j int) bool {
		return payload[i].Number < payload[j].Number
	})

	surroundings := make([]SurroundingsPayload, len(payload))
	for i, v := range payload {
		surroundings[i] = v
	}

	for i, v := range surroundings {
		fmt.Println(fmt.Printf("%d: number: %d, timestamp: %s, Temperature: %f, Moisture: %f, airPressure: %f", i, v.Number, v.Timestamp, v.Temperature, v.Moisture, v.AirPressure))
		p := influxdb2.NewPointWithMeasurement("vuoy_surroundings").
			AddTag("user", "bar").
			AddField("Temperature", v.Temperature).
			AddField("Moisture", v.Moisture).
			AddField("AirPressure", v.AirPressure).
			SetTime(time.Now())
		fmt.Printf("time: %s\n", v.Timestamp)
		fmt.Printf("timesamp: %s\n", v.Timestamp)
		writeAPI.WritePoint(p)
		defer client.Close()
	}
	return
}

func main() {
	payload := []SurroundingsPayload{
		{
			Number:      1,
			Timestamp:   time.Now(),
			Temperature: 25.5,
			Moisture:    0.8,
			AirPressure: 1013.25,
		},
		{
			Number:      2,
			Timestamp:   time.Now(),
			Temperature: 26.1,
			Moisture:    0.9,
			AirPressure: 1012.75,
		},
		{
			Number:      3,
			Timestamp:   time.Now(),
			Temperature: 25.8,
			Moisture:    0.7,
			AirPressure: 1013.50,
		},
	}

	insertPayload(payload)
}
