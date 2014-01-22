package customerio_test

import (
	"github.com/bmoyles0117/customerio-golang"
)

func ExampleIdentify() {
	customer_io := customerio.NewCustomerIO("XXX", "YYY")

	if err := customer_io.Identify("1", map[string]string{
		"email": "bryan.moyles@teltechcorp.com",
	}); err != nil {
		fmt.Println("Received an error while identifying a customer : %s", err)
	}

	if err := customer_io.Track("1", "signed_up", map[string]string{
		"plan": "best plan ever",
	}); err != nil {
		fmt.Println("Received an error while tracking a customer : %s", err)
	}
}
