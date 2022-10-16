package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	chart "github.com/wcharczuk/go-chart/v2"
)

var (
	SERVERPORT = os.Args[1]
)

type Data struct {
	H     string
	Day   int
	Time  int
	User  int
	Speed float64
}

type Choice struct {
	Choice string
}

var (
	h31_max = make([]float64, 7)
	h55_max = make([]float64, 7)
	h80_max = make([]float64, 7)
	h86_max = make([]float64, 7)
)

func max(h []float64) float64 {
	max := 0.0
	for i := 0; i < len(h); i++ {
		if max < h[i] {
			max = h[i]
		}
	}

	return max
}

func plot(x []float64, h []float64, title string) {
	var object = make([]chart.Tick, len(h))

	for i := 0; i < len(x); i++ {
		s := strconv.FormatFloat(x[i], 'g', 2, 64)
		object[i].Value = x[i]
		object[i].Label = s
	}

	graph := chart.Chart{
		Title: "Maximum is " + strconv.FormatFloat(max(h), 'f', 6, 64),
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		XAxis: chart.XAxis{
			Ticks: object,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    title,
				XValues: x,
				YValues: h,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	path := "../images/plot_" + title + ".png"

	f, err_f := os.Create(path)

	if err_f != nil {
		fmt.Println(err_f)
	}

	defer f.Close()

	err := graph.Render(chart.PNG, f)
	if err != nil {
		fmt.Println(err)
	}
}

func create_plot(w http.ResponseWriter, r *http.Request) {
	var data Choice
	x := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}

	req_err := json.NewDecoder(r.Body).Decode(&data)
	if req_err != nil {
		fmt.Println(req_err)
	}

	switch ch := data.Choice; ch {
	case "h31":
		plot(x, h31_max, "max_average_h31")
	case "h55":
		plot(x, h55_max, "max_average_h55")
	case "h80":
		plot(x, h80_max, "max_average_h80")
	case "h86":
		plot(x, h86_max, "max_average_h86")
	default:
		fmt.Println("Error")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	switch ch := data.H; ch {
	case "h31":
		if h31_max[data.Day-1] < data.Speed {
			h31_max[data.Day-1] = data.Speed
		}
	case "h55":
		if h55_max[data.Day-1] < data.Speed {
			h55_max[data.Day-1] = data.Speed
		}
	case "h80":
		if h80_max[data.Day-1] < data.Speed {
			h80_max[data.Day-1] = data.Speed
		}
	case "h86":
		if h86_max[data.Day-1] < data.Speed {
			h86_max[data.Day-1] = data.Speed
		}
	default:
		fmt.Println("Error")
	}

}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	r.HandleFunc("/plot", create_plot).Methods("POST")

	return r
}

func main() {
	r := newRouter()
	err := http.ListenAndServe(":"+SERVERPORT, r)

	if err != nil {
		fmt.Println(err)
	}
}
