package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

const (
	// Attempts connection
	host       = "localhost:8000"
	healthPath = "http://" + host + "/ping"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/api/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /auth/login.
func TestHTTPDoTranslate(t *testing.T) {
	body := `{
		"email": "bjorndonaldb@gmail.com",
		"password": "0123456789"
	}`
	Test(t,
		Description("Login Success"),
		Post(basePath+"/auth/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".message").Equal("text for translation"),
	)

	body = `{
		"email": "bjorndonaldb@gmail.com",
		"password": "012345678"
	}`
	Test(t,
		Description("Login Fail"),
		Post(basePath+"/auth/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".message").Equal("invalid request body"),
	)
}

// HTTP GET: /ping.
func TestHTTPPing(t *testing.T) {
	Test(t,
		Description("History Success"),
		Get(basePath+"/ping"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"history":[{`),
	)
}
