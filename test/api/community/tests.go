package community

import (
	"fmt"

	"github.com/shellhub-io/shellhub/test/api"
)

func Tests(address string) []api.Test {
	return []api.Test{
		{
			Description: "GET all namespaces from a user namespace",
			Method:      "GET",
			Url:         fmt.Sprintf("http://%s/api/namespaces", address),
			Kind:        "application/json",
			Cases: []api.Case{
				{
					Auth: api.Auth{
						Username: "henry",
						Password: "barreto",
					},
					Expected: 200,
				},
			},
		},
	}
}
