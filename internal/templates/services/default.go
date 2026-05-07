package services

import (
    "{{APP_NAME}}/uca"
    "{{APP_NAME}}/uca/context"
)

// To override a default CRUD method, uncomment and implement the function:
/*
func {{NAME}}GET(e *context.RequestEvent) error {
    // Custom logic before
    err := uca.DefaultGET(e, "{{NAME}}")
    // Custom logic after
    return err
}
*/

var _ = uca.DefaultGET
var _ *context.RequestEvent
