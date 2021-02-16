package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main data   `json:"main"`
}

type data struct {
	Temp float64 `json:"temp"`
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {

		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)

	})
	http.ListenAndServe(":8081", nil)

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(cityName string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=cc6a555439eb67b28442b4f9be1e4f64&q=" + cityName)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}
