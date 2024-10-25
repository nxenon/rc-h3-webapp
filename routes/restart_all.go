package routes

import (
	"encoding/json"
	"github.com/nxenon/rc-h3-webapp/db"
	"github.com/nxenon/rc-h3-webapp/utils"
	"net/http"
)

func restartAllRouteHandler(w http.ResponseWriter, r *http.Request) {
	appData := utils.LoadEnvFile("")
	db.DoChecks(appData)

	response := map[string]interface{}{
		"success": true,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return
}
