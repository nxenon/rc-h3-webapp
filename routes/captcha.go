package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/dchest/captcha"
	"net/http"
	"rc-h3-webapp/models"
)

// CaptchaRouteHandler /api/captcha
func CaptchaRouteHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	captchaID := captcha.New()

	var imgBuffer bytes.Buffer

	// Write the CAPTCHA image to the buffer in PNG format
	err := captcha.WriteImage(&imgBuffer, captchaID, 240, 80)
	if err != nil {
		http.Error(w, "Error generating CAPTCHA image", http.StatusInternalServerError)
		return
	}

	imgBase64 := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())

	// Create the response data
	response := models.CaptchaResponse{
		CaptchaID: captchaID,
		Captcha:   imgBase64,
	}

	json.NewEncoder(w).Encode(response)

}
