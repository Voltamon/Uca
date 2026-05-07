package services

import (
    "todo-app/uca"
    "todo-app/uca/context"
)

// To override a default CRUD method, uncomment and implement the function:
/*
func TaskGET(e *context.RequestEvent) error {
    // Custom logic before
    err := uca.DefaultGET(e, "Task")
    // Custom logic after
    return err
}
*/

var _ = uca.DefaultGET
var _ *context.RequestEvent
