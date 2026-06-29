package response

import (
	"log/slog"
	"net/http"
)

func ICSResponse(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(data)); err != nil {
		slog.Error("Write error", "error", err)
	}
}
