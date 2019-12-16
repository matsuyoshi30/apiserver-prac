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

type Server struct {
	GetTime func() time.Time // for test
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	var d time.Time
	if s.GetTime == nil {
		d = time.Now()
	} else {
		d = s.GetTime()
	}
	_, month, day := d.Date()

	// get request parameter
	p := r.FormValue("p")
	if p == "" {
		p = "matsuyoshi"
	}

	// set result
	var str string
	if month == time.January && (day == 1 || day == 2 || day == 3) {
		str = "大吉"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(o); err != nil {
		http.Error(w, "fail to encode result", http.StatusInternalServerError)
		return
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("Open http://localhost:8080/")

	s := Server{}
	http.HandleFunc("/", s.handler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
