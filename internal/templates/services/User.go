package services

import (
	"{{APP_NAME}}/uca"
	"{{APP_NAME}}/uca/context"
)

// UserGET and UserPOST are handled automatically by Uca's hidden CRUD.
// To add custom logic, simply implement the functions here.

// Keep these imports for your custom logic:
var _ = uca.DefaultGET
var _ *context.RequestEvent
