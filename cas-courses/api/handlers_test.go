package api

import (
	// "cas-calendar/model"
	// "dbclient"
	cfg "cas-courses/config"
	"cas-courses/service"
	"encoding/json"
	"fmt"
	// "github.com/alicebob/miniredis"
	"net/http"
	"net/http/httptest"
	"testing"
	"webserver"

	. "github.com/smartystreets/goconvey/convey"
)

// mock server, just need change url to mockserver.URL
// Eg: ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         w.Header().Set("Content-Type", "application/json")
//         fmt.Fprintln(w, `{"fake twitter json string"}`)
//     }))
//     defer ts.Close()

//     twitterUrl = ts.URL
//     c := make(chan *twitterResult)
//     go retrieveTweets(c)

//     tweet := <-c
//     if tweet != expected1 {
//         t.Fail()
//     }
//     tweet = <-c
//     if tweet != expected2 {
//         t.Fail()
//     }
// https://stackoverflow.com/questions/16154999/how-to-test-http-calls-in-go-using-httptest
func mockServer(log func(...interface{})) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		log("Send data to:" + r.RequestURI)
		log("Tut token:" + r.Header.Get("Tut"))
		log(r)
		if r.RequestURI == "/v1/campaigns/HOME-CHECK-IN/users/2000266949/profiles/2000266949/devices/Q7KL9iPGkq0hIZkegeFy4E2LkQZI2KGCOR3FpXbxwDK" {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprintln(w, "")
		} else {
			w.WriteHeader(400)
		}
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

// Redis mimic
// mredis, err := miniredis.Run()
// if err != nil {
// 	panic(err)
// }
// defer mredis.Close()

// initTest(mredis.Addr())

// End2end test
func TestGetCalendarE2E(t *testing.T) {
	envs := cfg.GetEnv()
	MysqlClient = &service.CalendarClient{}
	MysqlClient.Init(envs["mysqlUrl"])

	Convey("Given a HTTP request for /calendar", t, func() {
		// postData := make([]byte, 100)
		// req, err := http.NewRequest("POST", "http://example.com", bytes.NewReader(postData))
		// req.Header.Add("User-Agent", "myClient")
		// resp, err := client.Do(req)
		// defer resp.Body.Close()

		req := httptest.NewRequest("GET", "/calendar", nil)
		req.Header.Add("X-Requested-Un", "testschool")

		resp := httptest.NewRecorder()

		Convey("When the request is handled by the Router", func() {
			// http.DefaultServeMux.ServeHTTP(resp, req)
			webserver.NewRouter(&Routes).ServeHTTP(resp, req)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, 200)

				// calendars := make([]model.Calendar, 0, 5)
				// calendars := []model.Calendar{model.Calendar{}}
				calendars := []map[string]interface{}{}
				err := json.Unmarshal(resp.Body.Bytes(), &calendars)
				if err != nil {
					t.Log(err)
				}
				// https://github.com/smartystreets/goconvey/blob/master/examples/assertion_examples_test.go
				So(calendars[0]["id"], ShouldBeGreaterThan, 0)
				// So(account.Name, ShouldEqual, "Person_123")
				t.Log(string(resp.Body.Bytes()))
				t.Log(calendars)

				// b, _ := json.Marshal(calendars)
				// t.Log(string(b))
			})
		})
	})

	// Convey("Given a HTTP request for /accounts/456", t, func() {
	// 	req := httptest.NewRequest("GET", "/accounts/456", nil)
	// 	resp := httptest.NewRecorder()

	// 	Convey("When the request is handled by the Router", func() {
	// 		NewRouter().ServeHTTP(resp, req)

	// 		Convey("Then the response should be a 404", func() {
	// 			So(resp.Code, ShouldEqual, 404)
	// 		})
	// 	})
	// })
}

// Unittest
func TestGetCalendarUT(t *testing.T) {
	// Create a mock instance that implements the IBoltClient interface
	mockService := &service.MockCalendarClient{}

	// Declare two mock behaviours. For "123" as input, return a proper Account struct and nil as error.
	// For "456" as input, return an empty Account object and a real error.
	// mockRepo.On("QueryAccount", "123").Return(model.Account{Id: "123", Name: "Person_123"}, nil)
	// mockRepo.On("QueryAccount", "456").Return(model.Account{}, fmt.Errorf("Some error"))

	mockService.On("GetAllCalendars", 0, 0, "", int64(1), "true").Return([]interface{}{map[string]string{"a": "aa"}})
	mockService.On("GetBaseInfo", "testschool").Return(int64(1), "sname", "sid")

	// Finally, assign mockRepo to the DBClient field (it's in _handlers.go_, e.g. in the same package)
	// DBClient := mockRepo
	MysqlClient = mockService

	Convey("Given a HTTP request for /calendar", t, func() {
		req := httptest.NewRequest("GET", "/calendar?is_visible=true", nil)
		req.Header.Add("X-Requested-Un", "testschool")

		// Convey("Given a HTTP request for /test", t, func() {
		// 	req := httptest.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the Router", func() {
			webserver.NewRouter(&Routes).ServeHTTP(resp, req)

			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, 200)

				// account := model.Account{}
				// json.Unmarshal(resp.Body.Bytes(), &account)
				// So(account.Id, ShouldEqual, "123")
				// So(account.Name, ShouldEqual, "Person_123")
				t.Log(string(resp.Body.Bytes()))
			})
		})
	})
}
