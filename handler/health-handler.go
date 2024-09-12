package handler

import (
	"net/http"

	"github.com/guigateixeira/general-auth/util"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "running", "success": "true"})
}
