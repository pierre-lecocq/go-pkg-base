package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("JSON encoding error", "error", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func JSONError(w http.ResponseWriter, status int, msg string, err error) {
	if err != nil {
		slog.Error("HTTP Error", "status", status, "msg", msg, "error", err)
	} else {
		slog.Error("HTTP Error", "status", status, "msg", msg)
	}

	JSONResponse(w, status, map[string]string{
		"error":  msg,
		"status": http.StatusText(status),
	})
}
