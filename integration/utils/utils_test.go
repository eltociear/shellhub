package utils

import (
        "testing"
       )

func TestServerRunning(t *testing.T) {
        logs := `
          shellhub-api-1  |    ____    __
          shellhub-api-1  |   / __/___/ /  ___
          shellhub-api-1  |  / _// __/ _ \/ _ \
          shellhub-api-1  | /___/\__/_//_/\___/ v4.6.1
          shellhub-api-1  | High performance, minimalist Go web framework
          shellhub-api-1  | https://echo.labstack.com
          shellhub-api-1  | ____________________________________O/_______
          shellhub-api-1  |                                     O\
          shellhub-api-1  | ⇨ http server started on [::]:8080
      `
      val := ServerRunning(logs) 
      if !val {
          t.Errorf("error")
      }
}
