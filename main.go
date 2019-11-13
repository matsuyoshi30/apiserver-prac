package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Omikuji struct {
	Name      string    `json:"name"`
	Result    string    `json:"result"`
	ExecuteAt time.Time `json:"execute_at"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	if p == "" {
		p = "matsuyoshi"
	}

	str := "大吉"
	// 1/1 - 3 は常に大吉
	// それ以外はランダム
	d := time.Now()
	_, month, day := d.Date()
	if month == time.January && (day == 1 || day == 2 || day == 3) {
	} else {
		switch rand.Intn(6) {
		case 0:
			str = "大吉"
		case 1, 2:
			str = "中吉"
		case 3, 4:
			str = "小吉"
		case 5:
			str = "凶"
		}
	}

	o := &Omikuji{
		Name:      p,
		Result:    str,
		ExecuteAt: d,
	}

	json, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(json)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
