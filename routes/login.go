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

// LoginFrontRouteHandler /login -> only front
func LoginFrontRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
	AppData := utils.LoadEnvFile(".env")
	subStrings := strings.Split(AppData.H3ListenAddr, ":")
	subStringsHost := strings.Split(r.Host, ":")
	altSvcHeaderValue := fmt.Sprintf("h3=\"%s:%s\"", subStringsHost[0], subStrings[1])
	w.Header().Set("Alt-Svc", altSvcHeaderValue)
}

// LoginRouteHandler /api/login
func LoginRouteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest models.LoginRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	nanoseconds := time.Now().UnixNano()
	w.Header().Set("x-time", fmt.Sprintf("%d", nanoseconds))

	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		x := fmt.Sprintf("Error decoding JSON: %s", err)
		http.Error(w, x, http.StatusNotFound)
		return
	}

	isUserRateLimited, err := db.IsUserRateLimitedByUsername(loginRequest.Username)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"message": "rate limit check failed",
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isUserRateLimited {
		response := map[string]interface{}{
			"success": false,
			"message": "User is Rate Limited! - Wait 15 Mins!",
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(429)
		return
	}

	userObject, err := db.GetUserObjectByUsername(loginRequest.Username)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		response := map[string]interface{}{
			"success": false,
			"message": "Invalid Username or Password",
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err2 := db.IncreaseUserRateLimitByUsername(userObject.Username)
	if err2 != nil {
		response := map[string]interface{}{
			"success": false,
			"message": "Error in updating rate limit",
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isPasswordCorrect := utils.CheckPassword(loginRequest.Password, userObject.UserHashedPassword)
	if isPasswordCorrect {

		jwtToken, err := utils.GenerateJWT(userObject.Username, userObject.UserId)
		if err != nil {
			response := map[string]interface{}{
				"success": false,
				"message": "Could not generate JWT token",
			}
			json.NewEncoder(w).Encode(response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err3 := db.ResetUserRateLimitByUsername(userObject.Username)
		if err3 != nil {
			response := map[string]interface{}{
				"success": false,
				"message": "Error in resetting username rate limit to 0 after success login",
			}
			json.NewEncoder(w).Encode(response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Login Success",
			"jwt":     jwtToken,
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"success": false,
		"message": "Invalid Username or Password",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusUnauthorized)
	return

}
