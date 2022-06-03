package jadepoolsaas

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFundingWallets(t *testing.T) {
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	company := NewCompanyWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := company.GetFundingWallets()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFundingTransfer(t *testing.T) {
	from := "L6RayqPn4jXExW0"
	to := "e5dJyVp8R3B1m4o"
	coin := "ETH"
	value := "0.01"
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	company := NewCompanyWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := company.FundingTransfer(from, to, coin, value)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFundingRecords(t *testing.T) {
	response := map[string]interface{}{}

	queryHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := writeSuccessResponse(w, response)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(queryHandler))
	company := NewCompanyWithAddr(ts.URL, TestAppKey, TestAppSecret)
	_, err := company.GetFundingRecords(1, 10)
	if err != nil {
		t.Fatal(err)
	}
}
