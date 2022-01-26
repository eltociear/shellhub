package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestGetNamespace(t *testing.T) {
	login, ok := Logins["user2-namespace1"]
    if !ok || login == nil {
        t.Errorf("error")
        return
    }

	httpexpect.New(t, BaseURL).GET(fmt.Sprintf("/namespaces/%s", login.Namespace.TenantID)).
		WithHeader("Authorization", fmt.Sprintf("Bearer %s", login.Token)).
		Expect().
		Status(http.StatusOK)
}
