package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	testcases := []struct {
		inputParam     string
		server         Server
		expectedName   string
		isRandom       bool
		expectedResult string
	}{
		{
			inputParam:     "",
			server:         Server{},
			expectedName:   "matsuyoshi",
			isRandom:       true,
			expectedResult: "",
		},
		{
			inputParam: "",
			server: Server{
				GetTime: func() time.Time { return time.Now() },
			},
			expectedName:   "matsuyoshi",
			isRandom:       true,
			expectedResult: "",
		},
		{
			inputParam:     "gopher",
			server:         Server{},
			expectedName:   "gopher",
			isRandom:       true,
			expectedResult: "",
		},
		{
			inputParam: "",
			server: Server{
				GetTime: func() time.Time { return time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC) },
			},
			expectedName:   "matsuyoshi",
			isRandom:       false,
			expectedResult: "大吉",
		},
		{
			inputParam: "",
			server: Server{
				GetTime: func() time.Time { return time.Date(2009, time.January, 2, 0, 0, 0, 0, time.UTC) },
			},
			expectedName:   "matsuyoshi",
			isRandom:       false,
			expectedResult: "大吉",
		},
		{
			inputParam: "",
			server: Server{
				GetTime: func() time.Time { return time.Date(2009, time.January, 3, 0, 0, 0, 0, time.UTC) },
			},
			expectedName:   "matsuyoshi",
			isRandom:       false,
			expectedResult: "大吉",
		},
	}

	for _, tt := range testcases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		if tt.inputParam != "" {
			param := r.URL.Query()
			param.Add("p", tt.inputParam)
			r.URL.RawQuery = param.Encode()
		}

		tt.server.handler(w, r)
		rw := w.Result()
		defer rw.Body.Close()

		if rw.StatusCode != http.StatusOK {
			t.Fatal("unexpected status code")
		}

		o := Omikuji{}
		if err := json.NewDecoder(rw.Body).Decode(&o); err != nil {
			t.Fatal("failed to decode response")
		}

		if tt.expectedName != o.Name {
			t.Fatalf("unexpected response name: %s", o.Name)
		}

		if !tt.isRandom && tt.expectedResult != o.Result {
			t.Fatalf("unexpected response result: %s", o.Result)
		}
	}
}
