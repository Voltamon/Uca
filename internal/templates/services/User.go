package services

import (
    "net/http"

    "github.com/pocketbase/pocketbase/core"
)

type UserResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func UserGET(e *core.RequestEvent) error {
    records, err := e.App.FindAllRecords("User")
    if err != nil || len(records) == 0 {
        return e.JSON(http.StatusNotFound, map[string]string{"error": "no user found"})
    }

    record := records[0]
    return e.JSON(http.StatusOK, UserResponse{
    	Id:   record.Id,
		Name: record.GetString("name"),
		Role: record.GetString("role"),
	})
}

func UserPOST(e *core.RequestEvent) error {
	var body struct {
		Name string `json:"name"`
	}

	err := e.BindBody(&body)
	if err != nil || body.Name == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	collection, err := e.App.FindCollectionByNameOrId("User")
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "collection not found"})
	}

	defaultRole := os.Getenv("UCA_DEFAULT_ROLE")

	record := core.NewRecord(collection)
	record.Set("name", body.Name)
	if defaultRole != "" {
		record.Set("role", defaultRole)
	}

	err = e.App.Save(record)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save user"})
	}

	return e.JSON(http.StatusCreated, UserResponse{
		Id:   record.Id,
		Name: record.GetString("name"),
	})
}
