package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func SanitizeInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := makeBodyReusable(r)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Content-Type") == "application/json" && len(body) > 0 {
			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err == nil {
				for key, value := range jsonBody {
					if str, ok := value.(string); ok {
						jsonBody[key] = sanitize(str)
					}
				}
				sanitizedBody, _ := json.Marshal(jsonBody)
				r.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
				r.ContentLength = int64(len(sanitizedBody))
				r.Header.Set("Content-Length", string(r.ContentLength))
			}
		}

		if err := r.ParseForm(); err == nil {
			for key := range r.Form {
				r.Form.Set(key, sanitize(r.Form.Get(key)))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func sanitize(input string) string {
	return strings.TrimSpace(input)
}
