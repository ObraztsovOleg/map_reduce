package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	chart "github.com/wcharczuk/go-chart"
)

var (
	SERVERPORT = os.Args[1]
)

type Data struct {
	User  string
	Day   string
	Speed float64
}

var time_line = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	r.HandleFunc("/week_plot", create_plot).Methods("POST")

	return r
}

func main() {
	r := newRouter()
	err := http.ListenAndServe(":"+SERVERPORT, r)

	if err != nil {
		fmt.Println(err)
	}
}

func create_plot(w http.ResponseWriter, r *http.Request) {
	x := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: time_line,
			},
		},
	}

	f, err_f := os.Create("../images/week_plot.png")
	if err_f != nil {
		fmt.Println(err_f)
	}

	defer f.Close()

	err := graph.Render(chart.PNG, f)
	if err != nil {
		fmt.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	i, err := strconv.Atoi(data.Day)
	if err != nil {
		fmt.Println(err)
	}

	time_line[i-1] += data.Speed / float64(20)

	fmt.Println(time_line)
}
