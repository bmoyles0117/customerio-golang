package customerio

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type CustomerIO struct {
	SiteID    string
	ApiKey    string
	Host      string // The target host to transmit requests to
	UrlPrefix string // The URL prefix for the API, following the domain and before the requested path
	Doer      Doer   // This can be overridden for mocking / custom request completion.
}

// A shortcut function to quickly translate a map of strings to satisfy the default
// expectation of url.Values
func stringMapToValues(data map[string]string) *url.Values {
	values := &url.Values{}

	for k, v := range data {
		values.Add(k, v)
	}

	return values
}

// A default constructor that only asks for credential specifics, defaulting to the standard endpoints
func NewCustomerIO(site_id, api_key string) *CustomerIO {
	return &CustomerIO{site_id, api_key, "track.customer.io", "/api/v1", http.DefaultClient}
}

func (cio *CustomerIO) getEndpointUrl(path string) string {
	return fmt.Sprintf("https://%s:%s@%s%s%s", cio.SiteID, cio.ApiKey, cio.Host, cio.UrlPrefix, path)
}

func (cio *CustomerIO) Identify(customer_id string, data map[string]string) error {
	_, err := cio.sendRequest("PUT", cio.getEndpointUrl("/customers/"+customer_id), strings.NewReader(stringMapToValues(data).Encode()))

	return err
}

func (cio *CustomerIO) Delete(customer_id string) error {
	_, err := cio.sendRequest("DELETE", cio.getEndpointUrl("/customers/"+customer_id), nil)

	return err
}

func (cio *CustomerIO) sendRequest(method, url string, body io.Reader) (*http.Response, error) {
	var (
		err  error
		resp *http.Response
		req  *http.Request
	)

	if req, err = http.NewRequest(method, url, body); err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err = cio.Doer.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (cio *CustomerIO) SetDoer(doer Doer) {
	cio.Doer = doer
}

func (cio *CustomerIO) Track(customer_id, name string, data map[string]string) error {
	request_map := map[string]string{
		"name": name,
	}

	// We need to populate a nested url array of data values
	for k, v := range data {
		request_map[fmt.Sprintf("data[%s]", k)] = v
	}

	_, err := cio.sendRequest("POST", cio.getEndpointUrl("/customers/"+customer_id+"/events"), strings.NewReader(stringMapToValues(request_map).Encode()))

	return err
}
