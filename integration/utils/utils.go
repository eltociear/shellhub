package utils

import (
        "strings"
        "log"
)

func ServerRunning(logs string) bool {
    log.Println(logs)
    return strings.Contains(logs, "shellhub-api-1  | ⇨ http server started on [::]:8080")
}
