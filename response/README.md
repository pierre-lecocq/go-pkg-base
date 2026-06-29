# response

HTTP response writers for JSON APIs and iCalendar payloads.

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    user, err := store.Find(r.Context(), id)
    if err != nil {
        response.JSONError(w, http.StatusNotFound, "user not found", err)
        return
    }

    response.JSONResponse(w, http.StatusOK, user)
}

func getCalendar(w http.ResponseWriter, r *http.Request) {
    ics, err := calendar.Build(r.Context())
    if err != nil {
        response.JSONError(w, http.StatusInternalServerError, "calendar error", err)
        return
    }

    response.ICSResponse(w, http.StatusOK, ics)
}
```
