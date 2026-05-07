package services

import (
	"{{APP_NAME}}/uca"
	"{{APP_NAME}}/uca/context"
)

// HistoryGET and HistoryPOST are handled automatically by Uca's hidden CRUD.

// Actions like GetChatHistory can still be implemented manually
func GetChatHistory(e *context.RequestEvent) error {
	return uca.DefaultGET(e, "History")
}
