package services

import (
    "net/http"

    "github.com/pocketbase/pocketbase/core"
)

func {{NAME}}GET(e *core.RequestEvent) error {
    return e.JSON(http.StatusOK, map[string]string{"message": "{{NAME}} service"})
}
