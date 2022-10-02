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
	User  int
	H     string
	Day   string
	Speed float64
}

type Choice struct {
	Choice string
}

var (
	h31 = make([]float64, 7)
	h55 = make([]float64, 7)
	h80 = make([]float64, 7)
	h86 = make([]float64, 7)
)

var (
	h31_u = make([]float64, 20)
	h55_u = make([]float64, 20)
	h80_u = make([]float64, 20)
	h86_u = make([]float64, 20)
)

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

func plot(x []float64, h []float64, title string) {
	graph := chart.Chart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
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

	f, err_f := os.Create("../images/plot_" + title + ".png")
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
	x_u := []float64{
		1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0,
		11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0}

	req_err := json.NewDecoder(r.Body).Decode(&data)
	if req_err != nil {
		fmt.Println(req_err)
	}

	switch ch := data.Choice; ch {
	case "h31":
		plot(x, h31, "week_h31")
	case "h55":
		plot(x, h55, "week_h55")
	case "h80":
		plot(x, h80, "week_h80")
	case "h86":
		plot(x, h86, "week_h86")
	case "h31_u":
		plot(x_u, h31_u, "user_h31")
	case "h55_u":
		plot(x_u, h55_u, "user_h55")
	case "h80_u":
		plot(x_u, h80_u, "user_h80")
	case "h86_u":
		plot(x_u, h86_u, "user_h86")
	default:
		fmt.Println("Error")
	}
}

func addition(data Data, h *[]float64) {
	i, err := strconv.Atoi(data.Day)
	if err != nil {
		fmt.Println(err)
	}

	(*h)[i-1] += data.Speed / float64(7)
}

func addition_u(data Data, h *[]float64) {
	(*h)[data.User-1] += data.Speed / float64(20)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	switch ch := data.H; ch {
	case "h31":
		addition(data, &h31)
		addition_u(data, &h31_u)
	case "h55":
		addition(data, &h55)
		addition_u(data, &h55_u)
	case "h80":
		addition(data, &h80)
		addition_u(data, &h80_u)
	case "h86":
		addition(data, &h86)
		addition_u(data, &h86_u)
	default:
		fmt.Println("Error")
	}

	fmt.Println("h31", h31)
	fmt.Println("h55", h55)
	fmt.Println("h80", h80)
	fmt.Println("h86", h86)
	fmt.Println("u", data.User)
}
