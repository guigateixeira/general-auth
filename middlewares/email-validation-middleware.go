package middlewares

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/guigateixeira/general-auth/util"
)

func EmailValidatorMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := makeBodyReusable(r)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Content-Type") == "application/json" && len(body) > 0 {
			var requestBody map[string]string
			if err := json.Unmarshal(body, &requestBody); err != nil {
				util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}

			email, exists := requestBody["email"]
			if !exists {
				util.RespondWithError(w, http.StatusBadRequest, "Email field is missing")
				return
			}

			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			re := regexp.MustCompile(emailRegex)
			if !re.MatchString(email) {
				util.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
