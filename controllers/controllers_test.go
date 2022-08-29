package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInputSites(t *testing.T) {
	jsoninput := []byte(`{"websites":["www.google.com","www.facebook.com","www.fakewebsite1.com","https://linkedin.com/",
	"https://instagram.com/"]}`)
	req, err := http.NewRequest("POST", "/websites", bytes.NewBuffer(jsoninput))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(AddSites)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code: %v ,expected %v", status, http.StatusOK)
	}

	expected := "AddedWebsites added"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v ,expected %v",
			rr.Body.String(), expected)
	}
}

func TestGetAllSitesStatus(t *testing.T) {
	time.Sleep(2 * time.Second)
	req, err := http.NewRequest("GET", "/websites", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetSiteStatus)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code: %v ,expected %v", status, http.StatusOK)
	}

	expected := `{"https://instagram.com/":"working","https://linkedin.com/":"working","www.facebook.com":"not working","www.fakewebsite1.com":"not working","www.google.com":"not working"}`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v ,expected %v",
			rr.Body.String(), expected)
	}

}

func TestGetSingleSiteStatus(t *testing.T) {
	testlinks := []struct {
		name           string
		link           string
		want           string
		wantstatuscode int
	}{
		{"link is stored", "www.google.com", "\"not working\"", 200},
		{"link is not stored", "www.instantt.com", "www.instantt.com is not stored in the database", 404},
	}

	for _, test := range testlinks {

		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/websites", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetSiteStatus)
			q := req.URL.Query()
			q.Add("link", test.link)
			req.URL.RawQuery = q.Encode()
			handler.ServeHTTP(rr, req)
			if rr.Code != test.wantstatuscode {
				t.Errorf("handler returned status code: %v ,expected %v", rr.Code, http.StatusOK)
			}
			expected := test.want
			if rr.Body.String() != string(expected) {
				t.Errorf("handler returned unexpected body: got %s ,expected %s",
					rr.Body.String(), expected)
			}

		})

	}
}
