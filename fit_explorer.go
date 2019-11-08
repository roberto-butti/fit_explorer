package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/tormoder/fit"
)

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func convertToMinutes(m float64) string {
	minutes, secondsdecimal := math.Modf(m)
	seconds := secondsdecimal * 60
	return fmt.Sprintf("%02.0f:%02.0f", minutes, seconds)
}

func cInt(i int) string {
	v := strconv.Itoa(i)
	return v
}

func cInt8(i int8) string {
	v := strconv.Itoa(int(i))
	return v
}
func cUint8(i uint8) string {
	v := strconv.Itoa(int(i))
	return v
}
func cDecimal(f float64) string {
	v := strconv.FormatFloat(f, 'f', -1, 64)
	return v
}
func cTime(t time.Time) string {
	v := t.String()
	return v
}
func cDuration(d time.Duration) string {
	v := d.String()
	return v
}

func main() {
	// Read our FIT test file data
	testFilePtr := flag.String("file", "file.fit", "the FIT file name")
	flag.Parse()
	fmt.Println("I will open FIT file:", *testFilePtr)

	testData, err := ioutil.ReadFile(*testFilePtr)
	if err != nil {
		log.Fatalln("error reading fit file:", err)
		return
	}

	// Decode the FIT file data
	fit, err := fit.Decode(bytes.NewReader(testData))
	if err != nil {
		log.Fatalln("error decoding fir file:", err)
		return
	}

	w := csv.NewWriter(os.Stdout)

	// Inspect the TimeCreated field in the FileId message
	log.Println("TIME CREATED:", fit.FileId.TimeCreated)

	// Inspect the dynamic Product field in the FileId message
	log.Println("PRODUCT:", fit.FileId.GetProduct())

	// Inspect the FIT file type
	log.Println("FIT TYPE:", fit.Type())

	// Get the actual activity
	activity, err := fit.Activity()
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
		return
	}

	idx := 0
	t0 := time.Now()
	for _, record := range activity.Records {
		idx++
		if idx == 1 {
			t0 = record.Timestamp
		}

		x := []string{
			cInt(idx),
			cDecimal(record.PositionLat.Degrees()),
			cDecimal(record.PositionLong.Degrees()),
			cInt8(record.Temperature),
			cTime(record.Timestamp),
			cDuration(record.Timestamp.Sub(t0)),
			cUint8(record.Cadence),
			cDecimal(record.GetDistanceScaled() / 1000),
			"km",
			cDecimal(record.GetAltitudeScaled()),
			"m",
			cUint8(record.HeartRate),
			"HR bpm",
			convertToMinutes(16.666666666667 / record.GetSpeedScaled()),
			"min/km",
			cDecimal(record.GetSpeedScaled() * 3.6),
			"km/h",
		}
		w.Write(x)
	}
	w.Flush()

	// Print the sport of the first Session message
	for _, session := range activity.Sessions {
		log.Println("SESSION SPORT:", session.Sport)
	}

}
