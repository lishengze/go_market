package jadepoolsaas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCreateAddress(t *testing.T) {
	coin := "ETH"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, coin) {
			t.Errorf("request url = %s; want contain %+v", r.URL.Path, coin)
		}

		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.CreateAddress(coin)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyAddress(t *testing.T) {
	coin := "ETH"
	address := "0x7C3A4d3ff2b92CFDD2eD1a105d5bAc8fAF4008aE"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.VerifyAddress(coin, address)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAssets(t *testing.T) {
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.GetAssets()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBalance(t *testing.T) {
	coin := "ETH"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, coin) {
			t.Errorf("request url = %s; want contain %+v", r.URL.Path, coin)
		}

		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.GetBalance(coin)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrder(t *testing.T) {
	id := "rNXBQGJlw09apVyg4nDo"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.GetOrder(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithdraw(t *testing.T) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	coin := "ETH"
	to := "0x7C3A4d3ff2b92CFDD2eD1a105d5bAc8fAF4008aE"
	value := "0.01"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.Withdraw(id, coin, to, value)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelegate(t *testing.T) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	coin := "IRIS"
	value := "0.01"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.Delegate(id, coin, value)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnDelegate(t *testing.T) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	coin := "IRIS"
	value := "0.01"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.UnDelegate(id, coin, value)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetValidators(t *testing.T) {
	coin := "IRIS"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, coin) {
			t.Errorf("request url = %s; want contain %+v", r.URL.Path, coin)
		}

		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.GetValidators(coin)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStakingInterest(t *testing.T) {
	coin := "IRIS"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, coin) {
			t.Errorf("request url = %s; want contain %+v", r.URL.Path, coin)
		}

		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	fmt.Println(ts.URL)
	ts.URL = "http://127.0.0.1:8092"
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	result, err := app.GetStakingInterest(coin, "2019-09-26")
	fmt.Println(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddUrgentStakingFunding(t *testing.T) {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	coin := "IRIS"
	value := "0.01"
	expiredAt := time.Now().Unix() + 86400*21
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	app := NewAppWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := app.AddUrgentStakingFunding(id, coin, value, expiredAt)
	if err != nil {
		t.Fatal(err)
	}
}
