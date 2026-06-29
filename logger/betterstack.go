package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type betterstackCore struct {
	endpoint string
	token    string
	client   *http.Client
	queue    chan map[string]any
	wg       sync.WaitGroup
}

func (c *betterstackCore) drain() {
	defer c.wg.Done()

	batch := make([]map[string]any, 0, 16)
	ticker := time.NewTicker(time.Second)

	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		_ = c.send(batch)
		batch = batch[:0]
	}

	for {
		select {
		case payload, ok := <-c.queue:
			if !ok {
				flush()
				return
			}

			batch = append(batch, payload)
			if len(batch) >= 16 {
				flush()
				ticker.Reset(time.Second)
			}

		case <-ticker.C:
			flush()
		}
	}
}

func (c *betterstackCore) send(batch []map[string]any) error {
	body, err := json.Marshal(batch)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() // nolint: errcheck

	return nil
}

type betterstackHandler struct {
	core  *betterstackCore
	attrs []slog.Attr
}

func newBetterstackHandler(endpoint string, token string, globalAttrs []slog.Attr) *betterstackHandler {
	core := &betterstackCore{
		endpoint: endpoint,
		token:    token,
		client:   &http.Client{Timeout: 5 * time.Second},
		queue:    make(chan map[string]any, 256),
	}

	core.wg.Add(1)
	go core.drain()

	return &betterstackHandler{core: core, attrs: globalAttrs}
}

func (h *betterstackHandler) Close() {
	close(h.core.queue)
	h.core.wg.Wait()
}

func (h *betterstackHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }

func (h *betterstackHandler) Handle(_ context.Context, r slog.Record) error {
	payload := map[string]any{
		"dt":      r.Time.UTC().Format(time.RFC3339Nano),
		"level":   r.Level.String(),
		"message": r.Message,
	}

	for _, a := range h.attrs {
		payload[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		payload[a.Key] = a.Value.Any()
		return true
	})

	select {
	case h.core.queue <- payload:
	default:
	}

	return nil
}

func (h *betterstackHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make([]slog.Attr, len(h.attrs)+len(attrs))

	copy(merged, h.attrs)
	copy(merged[len(h.attrs):], attrs)

	return &betterstackHandler{core: h.core, attrs: merged}
}

func (h *betterstackHandler) WithGroup(_ string) slog.Handler { return h }

type teeHandler struct{ handlers []slog.Handler }

func (t *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range t.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (t *teeHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range t.handlers {
		if h.Enabled(ctx, r.Level) {
			_ = h.Handle(ctx, r.Clone())
		}
	}

	return nil
}

func (t *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}

	return &teeHandler{handlers: handlers}
}

func (t *teeHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		handlers[i] = h.WithGroup(name)
	}

	return &teeHandler{handlers: handlers}
}

func SetupBetterStackSlog(endpoint string, token string, globalAttrs ...slog.Attr) func() {
	stderr := slog.NewTextHandler(os.Stderr, nil)
	if token == "" || endpoint == "" {
		slog.SetDefault(slog.New(stderr))
		slog.Warn("BETTERSTACK_TOKEN or BETTERSTACK_ENDPOINT env vars not set, defaulting to stderr logger")
		return func() {}
	}

	if !strings.HasPrefix(endpoint, "https://") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}

	bs := newBetterstackHandler(endpoint, token, globalAttrs)
	slog.SetDefault(slog.New(&teeHandler{handlers: []slog.Handler{stderr, bs}}))

	return bs.Close
}
