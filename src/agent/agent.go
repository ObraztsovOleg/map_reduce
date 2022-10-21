package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chrislusf/glow/flow"
)

var (
	FOLDER  = os.Args[1]
	WEBPORT = os.Args[2]
	WEBHOST = os.Args[3]
)

type Time_line struct {
	Day    int64
	Speed  float64
	Time   string
	Count  int
	Filter bool
}

type Data struct {
	H     string  `json:"H"`
	Day   int     `json:"Day"`
	User  int     `json:"User"`
	Time  string  `json:"Time"`
	Speed float64 `json:"Speed"`
}

const DIR = "/home/obrol/Downloads/"

func read_csv_file(file_path string) [][]string {
	f, err := os.Open(file_path)
	if err != nil {
		log.Fatal("Unable to read input file "+file_path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+file_path, err)
	}

	return records
}

func send_request(h string, day int, user int, time string, speed float64) {
	var values Data

	values.Day = day
	values.Speed = speed
	values.Time = time
	values.H = h
	values.User = user

	json_data, err := json.Marshal(values)

	res, err := http.Post("http://"+WEBHOST+":"+WEBPORT+"/",
		"application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
		fmt.Println(res)
	}
}

func new_flow(path string, file_name string) {
	flow.New().TextFile(
		path, 3,
	).Map(func(line string, ch chan Time_line) {
		var time_line Time_line
		new_line := strings.Split(line, ",")

		for index, elem := range new_line {
			new_line[index] = strings.Trim(elem, "\t ")
		}

		speed, err := strconv.ParseFloat(new_line[2], 64)
		if err != nil {
			time_line.Filter = false
			ch <- time_line
			return
		}

		time, err := strconv.ParseInt(new_line[1], 10, 64)
		if err != nil {
			time_line.Filter = false
			ch <- time_line
			return
		}

		day, err := strconv.ParseInt(new_line[0], 10, 64)
		if err != nil {
			time_line.Filter = false
			ch <- time_line
			return
		}

		if day == 0 {
			time_line.Filter = false
			ch <- time_line
			return
		} else {
			time_line.Filter = true
			time_line.Speed = speed
			time_line.Time = strconv.FormatInt((day-1)*86400+time, 10)
			time_line.Day = day
			time_line.Count = 1

			ch <- time_line
		}
	}).Filter(func(src Time_line) bool {
		return src.Filter
	}).Map(func(src Time_line) flow.KeyValue {
		return flow.KeyValue{src.Time, src} // max average speed
	}).ReduceByKey(func(x Time_line, y Time_line) Time_line {
		x.Speed = x.Speed + y.Speed
		x.Count = x.Count + y.Count
		return x
	}).Map(func(time string, obj Time_line) {
		avg_speed := obj.Speed / float64(obj.Count)
		str := strings.Split(file_name, ".")
		user_str := strings.Split(str[2], "u")

		user, err := strconv.ParseInt(user_str[1], 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		send_request(str[1], int(obj.Day), int(user), obj.Time, avg_speed)
	}).Run()
}

func main() {
	if FOLDER != "master" {
		files, err := ioutil.ReadDir(DIR + "/" + FOLDER)
		if err != nil {
			log.Fatalf("unable to read dir: %v", err)
		}
		for _, file := range files {
			fmt.Println(file.Name() + " processing...")
			new_flow(DIR+"/"+FOLDER+"/"+file.Name(), file.Name())
		}
	}

	fmt.Println("Done")
}
