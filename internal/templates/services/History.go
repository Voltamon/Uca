package services

import (
    "net/http"

    "github.com/pocketbase/pocketbase/core"
)

type Message struct {
    Sender    string `json:"sender"`
    Content   string `json:"content"`
    Timestamp string `json:"timestamp"`
}

func HistoryGET(e *core.RequestEvent) error {
    records, err := e.App.FindAllRecords("History")
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch history"})
    }

    messages := make([]Message, len(records))
    for i, r := range records {
        messages[i] = Message{
            Sender:    r.GetString("sender"),
            Content:   r.GetString("content"),
            Timestamp: r.GetDateTime("created").String(),
        }
    }

    return e.JSON(http.StatusOK, messages)
}

func HistoryPOST(e *core.RequestEvent) error {
    var body struct {
        Sender  string `json:"sender"`
        Content string `json:"content"`
    }

    err := e.BindBody(&body)
    if err != nil || body.Sender == "" || body.Content == "" {
        return e.JSON(http.StatusBadRequest, map[string]string{"error": "sender and content are required"})
    }

    collection, err := e.App.FindCollectionByNameOrId("History")
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{"error": "collection not found"})
    }

    record := core.NewRecord(collection)
    record.Set("sender", body.Sender)
    record.Set("content", body.Content)

    err = e.App.Save(record)
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save message"})
    }

    return e.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}

func GetChatHistory(e *core.RequestEvent) error {
    records, err := e.App.FindAllRecords("History")
    if err != nil {
        return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch history"})
    }

    messages := make([]Message, len(records))
    for i, r := range records {
        messages[i] = Message{
            Sender:    r.GetString("sender"),
            Content:   r.GetString("content"),
            Timestamp: r.GetDateTime("created").String(),
        }
    }

    return e.JSON(http.StatusOK, messages)
}
