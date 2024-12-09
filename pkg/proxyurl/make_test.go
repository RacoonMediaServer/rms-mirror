package proxyurl

import (
	"net/url"
	"testing"
)

func TestMake(t *testing.T) {
	baseUrl, err := url.Parse("http://127.0.0.1:8080/proxy/")
	if err != nil {
		t.Error("Unexpected parser fail", err)
	}

	m := Maker{BaseURL: baseUrl}
	u := m.MakeURL("https://image.tmdb.org/t/p/w780/3L8uAS38M7sgehBoZ7ao2FooG2j.jpg")
	if u != "http://127.0.0.1:8080/proxy/?href=aHR0cHM6Ly9pbWFnZS50bWRiLm9yZy90L3Avdzc4MC8zTDh1QVMzOE03c2dlaEJvWjdhbzJGb29HMmouanBn" {
		t.Errorf("Unexpected url maked: %s", u)
	}
}
