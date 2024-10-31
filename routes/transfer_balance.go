package routes

import (
	"encoding/json"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/db"
	"github.com/nxenon/rc-h3-webapp/models"
	"github.com/nxenon/rc-h3-webapp/utils"
	"io"
	"net/http"
	"strings"
	"time"
)

func transferBalanceRouteHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode the JSON body to get the couponValue
	var request models.TransferBalanceRequestModel
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	nanoseconds := time.Now().UnixNano()
	w.Header().Set("x-time", fmt.Sprintf("%d", nanoseconds))

	err = json.Unmarshal(body, &request)
	if err != nil {
		x := fmt.Sprintf("Error decoding JSON: %s", err)
		http.Error(w, x, http.StatusInternalServerError)
		return
	}

	// Get user ID (assumed to be set in context or session)
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = db.TransferMoneyToUser(userId, request.ToUsername, request.Amount)
	if err != nil {
		x := fmt.Sprintf("Error in transfering balance: %s", err)
		http.Error(w, x, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Transfer Successfully Done",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return

}

func transferBalanceFrontRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/transfer_balance.html")
	AppData := utils.LoadEnvFile(".env")
	subStrings := strings.Split(AppData.H3ListenAddr, ":")
	subStringsHost := strings.Split(r.Host, ":")
	altSvcHeaderValue := fmt.Sprintf("h3=\"%s:%s\"", subStringsHost[0], subStrings[1])
	w.Header().Set("Alt-Svc", altSvcHeaderValue)
}
