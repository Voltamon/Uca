package services

import (
    "blog-app/uca"
    "blog-app/uca/context"
)

// To override a default CRUD method, uncomment and implement the function:
/*
func PostGET(e *context.RequestEvent) error {
    // Custom logic before
    err := uca.DefaultGET(e, "Post")
    // Custom logic after
    return err
}
*/

var _ = uca.DefaultGET
var _ *context.RequestEvent
