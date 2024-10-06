package routes

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/nxenon/rc-h3-webapp/db"
	"github.com/nxenon/rc-h3-webapp/models"
	"github.com/nxenon/rc-h3-webapp/utils"
	"net/http"
	"strings"
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

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginRequest.CaptchaValue == "111" {
		// for testing purposes
	} else {
		if !captcha.VerifyString(loginRequest.CaptchaID, loginRequest.CaptchaValue) {
			response := map[string]interface{}{
				"success": false,
				"message": "Invalid CAPTCHA",
			}
			json.NewEncoder(w).Encode(response)
			w.WriteHeader(http.StatusForbidden)
			return
		}
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
