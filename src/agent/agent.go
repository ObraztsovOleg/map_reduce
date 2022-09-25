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
	Day   string
	Speed float64
	Count int
}

type Data struct {
	User  string  `json:"User"`
	Day   string  `json:"Day"`
	Speed float64 `json:"Speed"`
}

const DIR = "/home/obrol/Downloads/BigData/Практики/1/v1"

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

func send_request(user string, day string, speed float64) {
	var values Data

	values.Day = day
	values.User = user
	values.Speed = speed

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
	).Filter(func(line string) bool {
		return !(strings.Contains(line, "Day") ||
			strings.Contains(line, "obrol-HP-ProBook-430-G5"))
	}).Map(func(line string, ch chan Time_line) {
		var time_line Time_line

		new_line := strings.Split(line, ",")

		for index, elem := range new_line {
			new_line[index] = strings.Trim(elem, "\t ")
		}

		float_speed, err := strconv.ParseFloat(new_line[2], 64)

		if err != nil {
			fmt.Println(err)
		}

		time_line.Speed = float_speed
		time_line.Day = new_line[0]
		time_line.Count = 1

		ch <- time_line
	}).Map(func(src Time_line) flow.KeyValue {
		// fmt.Println(src)
		return flow.KeyValue{src.Day, src}
	}).ReduceByKey(func(x Time_line, y Time_line) Time_line {
		x.Speed = x.Speed + y.Speed
		x.Count = x.Count + y.Count
		return x
	}).Map(func(day string, obj Time_line) {
		speed := obj.Speed / float64(obj.Count)
		send_request(strings.Split(file_name, ".")[2], obj.Day, speed)
	}).Run()
}

func main() {
	// json_file, err := os.Open("time.json")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer json_file.Close()

	if FOLDER != "master" {
		files, err := ioutil.ReadDir(DIR + "/" + FOLDER)
		if err != nil {
			log.Fatalf("unable to read dir: %v", err)
		}
		for _, file := range files {
			new_flow(DIR+"/"+FOLDER+"/"+file.Name(), file.Name())
		}
	}

	// if os.Args[1] == "master" {
	// 	fmt.Println(os.Args)
	// } else {
	// new_flow(DIR+"/h31/userlog.h31.u13.csv", "userlog.h31.u13.csv")
	// }
	// byteValue, _ := ioutil.ReadAll(json_file)
	// var time_line Time_line

	// content, err := json.Marshal(time_line)

	// json.Unmarshal(byteValue, &time_line)

	// for i := 0; i < len(time_line.Time_line); i++ {
	// 	fmt.Println("Speed: " + time_line.Time_line[i].Speed)
	// 	fmt.Println("Number: " + strconv.Itoa(time_line.Time_line[i].Number))
	// }

	// if files[0] == nil {
	// 	log.Fatalf("File dose not exist")
	// }

	// records := read_csv_file(os.Getenv("HOME") + DIR + "/" + files[0].Name())
	// records = records[1:]

	// fmt.Println(records[0])
	// new_flow("/home/obrol/Downloads/BigData/Практики/1/v1/userlog.h31.u18.csv")
}
