package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIReponse(t *testing.T) {
	testCases := []struct {
		url          string
		expectedCode int
	}{
		{
			"/tasks/?usr=ronaldo",
			http.StatusOK,
		},
		{
			"/tasks/?usr=ronaldo&status=andamento,finalizado",
			http.StatusOK,
		},
		{
			"/tasks/complete/?usr=ronaldo&id=1",
			http.StatusOK,
		},
		{
			"/tasks/restore/?usr=ronaldo&id=1",
			http.StatusOK,
		},
		{
			"/tasks/?notparsed=test",
			http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", tc.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		hf := http.HandlerFunc(tasksHandler)
		hf.ServeHTTP(rr, req)
		if s := rr.Code; s != tc.expectedCode {
			body, _ := ioutil.ReadAll(rr.Body)
			t.Logf("failed endpoint '%s' request, got status code '%d', message '%s'", tc.url, s, string(body))
			t.Fail()
		}
	}
}
