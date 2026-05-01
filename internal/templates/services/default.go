package services

import (
    "net/http"

    "{{APP_NAME}}/uca/context"
)

func {{NAME}}GET(e *context.RequestEvent) error {
    return e.JSON(http.StatusOK, map[string]string{"message": "{{NAME}} service"})
}
