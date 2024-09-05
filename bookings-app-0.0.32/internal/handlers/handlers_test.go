package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tsawler/bookings-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "Post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "Post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "Post", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Generals Quarterr",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code %d and wanted %d", rr.Code, http.StatusOK)
	}

}

func TestPostReservation(t *testing.T) {
	reqBody := "start_date=2020-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2020-12-2")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/www.form-url-encoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code %d and wanted %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_AvailablityJSON(t *testing.T) {
	reqBody := "start=2020-12-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-12-2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2020-12-2")
	req, _ := http.NewRequest("GET", "/search-availability", strings.NewReader(reqBody))
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-type", "application/www.form-url-encoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)
	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err!= nil {
		t.Error("failed to parse json")
	}
	

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code %d and wanted %d", rr.Code, http.StatusOK)
	}

}
