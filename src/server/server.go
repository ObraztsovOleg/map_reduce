package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	chart "github.com/wcharczuk/go-chart/v2"
)

var (
	SERVERPORT = os.Args[1]
)

type Data struct {
	H     string
	Day   int
	Time  string
	Count int
	Speed float64
}

type Choice struct {
	Choice string
}

var (
	h31       = make(map[string]float64)
	h55       = make(map[string]float64)
	h80       = make(map[string]float64)
	h86       = make(map[string]float64)
	count_h31 = make(map[string]int)
	count_h55 = make(map[string]int)
	count_h80 = make(map[string]int)
	count_h86 = make(map[string]int)
)

var h_max = make([]float64, 4)

func plot(x []string, h []float64, title string) {
	var object = make([]chart.Tick, len(h))
	var x_val = []float64{0.0, 1.0, 2.0, 3.0}

	for i := 0; i < len(x); i++ {
		object[i].Value = float64(i)
		object[i].Label = x[i]
	}

	graph := chart.Chart{
		Title: title,
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
				XValues: x_val,
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
	x := []string{"h31", "h55", "h80", "h86"}

	h_max[0] = max_avg(h31, count_h31)
	h_max[1] = max_avg(h55, count_h55)
	h_max[2] = max_avg(h80, count_h80)
	h_max[3] = max_avg(h86, count_h86)

	req_err := json.NewDecoder(r.Body).Decode(&data)
	if req_err != nil {
		fmt.Println(req_err)
	}

	plot(x, h_max, "max_average_h")
	for i := 0; i < len(x); i++ {
		fmt.Printf("%s - %.7f\n", x[i], h_max[i])
	}
}

func max_avg(h map[string]float64, count_h map[string]int) float64 {
	max := 0.0

	for key := range h {
		avg := h[key] / float64(count_h[key])

		if max < avg {
			max = avg
		}
	}

	return max
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	switch ch := data.H; ch {
	case "h31":
		h31[data.Time] = h31[data.Time] + data.Speed
		count_h31[data.Time] = count_h31[data.Time] + data.Count
	case "h55":
		h55[data.Time] = h55[data.Time] + data.Speed
		count_h55[data.Time] = count_h55[data.Time] + data.Count
	case "h80":
		h80[data.Time] = h80[data.Time] + data.Speed
		count_h80[data.Time] = count_h80[data.Time] + data.Count
	case "h86":
		h86[data.Time] = h86[data.Time] + data.Speed
		count_h86[data.Time] = count_h86[data.Time] + data.Count
	default:
		fmt.Println("Error")
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	r.HandleFunc("/plot", create_plot).Methods("GET")

	return r
}

func main() {
	r := newRouter()
	err := http.ListenAndServe(":"+SERVERPORT, r)

	if err != nil {
		fmt.Println(err)
	}
}
