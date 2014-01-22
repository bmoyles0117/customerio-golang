package customerio

import (
	"fmt"
	"net/http"
	"testing"
)

var site_id = "XXX"
var api_key = "YYY"

type dummyDoer struct {
	status         string
	status_code    int
	was_called     bool
	request_called *http.Request
}

func (d *dummyDoer) Do(req *http.Request) (resp *http.Response, err error) {
	d.was_called = true
	d.request_called = req

	return &http.Response{
		Status:     d.status,
		StatusCode: d.status_code,
	}, nil
}

func newDummyDoer(status string, status_code int) *dummyDoer {
	return &dummyDoer{status, status_code, false, nil}
}

func TestNewCustomerIO(t *testing.T) {
	customer_io := NewCustomerIO(site_id, api_key)

	if customer_io.site_id != site_id {
		t.Errorf("Invalid site id was associated to the customer.io client")
	}

	if customer_io.api_key != api_key {
		t.Errorf("Invalid api key was associated to the customer.io client")
	}
}

func TestCustomerIOgetEndpointUrl(t *testing.T) {
	customer_io := NewCustomerIO(site_id, api_key)

	if customer_io.getEndpointUrl("/customer/5") != fmt.Sprintf("https://%s:%s@track.customer.io/api/v1/customer/5", site_id, api_key) {
		t.Errorf("The endpoint url did not match the expectation")
	}
}

func TestCustomerIOIdentify(t *testing.T) {
	dummy_doer := newDummyDoer("200 OK", 200)
	customer_io := NewCustomerIO(site_id, api_key)
	customer_io.SetDoer(dummy_doer)

	if err := customer_io.Identify("5", map[string]string{
		"email": "customer@example.com",
		"name":  "Bob",
		"plan":  "premium",
	}); err != nil {
		t.Errorf("An error was returned by customer io Identify %s", err)
	}

	if !dummy_doer.was_called {
		t.Errorf("Dummy doer was never called, indicating the request would not have actually been invoked.")
	}

	if dummy_doer.request_called.FormValue("email") != "customer@example.com" {
		t.Errorf("Request form data was not populated as expected")
	}
}

func TestCustomerIOSetDoer(t *testing.T) {
	customer_io := NewCustomerIO(site_id, api_key)

	if customer_io.request_doer == nil {
		t.Errorf("The default request doer was not set")
	}

	dummy_doer := newDummyDoer("200 OK", 200)

	customer_io.SetDoer(dummy_doer)

	if customer_io.request_doer != dummy_doer {
		t.Errorf("The request doer was not overwritten properly")
	}

	if dummy_doer.was_called {
		t.Errorf("The dummy doer indicates that it was called before any invocation")
	}
}

func TestCustomerIOstringMapToValues(t *testing.T) {
	values := stringMapToValues(map[string]string{"test": "value"})

	if values.Get("test") != "value" {
		t.Errorf("The data map was not translated into url values properly")
	}
}

func TestCustomerIOTrack(t *testing.T) {
	customer_io := NewCustomerIO(site_id, api_key)
	dummy_doer := newDummyDoer("200 OK", 200)
	customer_io.SetDoer(dummy_doer)

	if err := customer_io.Track("5", "sample_event", map[string]string{
		"said": "hello",
	}); err != nil {
		t.Errorf("An error was returned by customer io Identify %s", err)
	}

	if !dummy_doer.was_called {
		t.Errorf("Dummy doer was never called, indicating the request would not have actually been invoked.")
	}

	if dummy_doer.request_called.FormValue("name") != "sample_event" {
		t.Errorf("Request form data was not populated as expected : %s", dummy_doer.request_called.Form)
	}

	if dummy_doer.request_called.FormValue("data[said]") != "hello" {
		t.Errorf("Request form data was not populated as expected : %s", dummy_doer.request_called.Form)
	}
}
