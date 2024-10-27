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
)

func cartFrontRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/cart.html")
	AppData := utils.LoadEnvFile(".env")
	subStrings := strings.Split(AppData.H3ListenAddr, ":")
	subStringsHost := strings.Split(r.Host, ":")
	altSvcHeaderValue := fmt.Sprintf("h3=\"%s:%s\"", subStringsHost[0], subStrings[1])
	w.Header().Set("Alt-Svc", altSvcHeaderValue)
}

func myCartRouteHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userCart, err := db.GetCartByUserId(userId)
	if err != nil {
		http.Error(w, "Error getting Cart by cart id!", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(userCart)
	if err != nil {
		http.Error(w, "Error converting cart to json!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func applyCouponRouteHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode the JSON body to get the couponValue
	var request models.ApplyCouponRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		x := fmt.Sprintf("Error decoding JSON: %s", err)
		http.Error(w, x, http.StatusNotFound)
		return
	}

	// Get user ID (assumed to be set in context or session)
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	coupon, err := db.GetCouponByValue(request.CouponValue)
	invalidCouponMsg := "Invalid Coupon Code!"
	if err != nil {
		http.Error(w, invalidCouponMsg, http.StatusNotFound)
		return
	}

	if coupon.IsValid != 1 {
		http.Error(w, invalidCouponMsg, http.StatusNotFound)
		return
	}

	userCartObject, err4 := db.GetCartByUserId(userId)
	if err4 != nil {
		http.Error(w, "Error getting user cart by user id!", http.StatusInternalServerError)
		return
	}

	// update user cart
	newPrice := CalculateDiscountedPrice(userCartObject.CartOverallPrice, coupon.DiscountPercent)
	userCartObject.CartOverallPrice = newPrice

	err3 := db.UpdateCartByUserId(userId, userCartObject)
	if err3 != nil {
		http.Error(w, "Error Updating User Cart Overall price!", http.StatusInternalServerError)
		return
	}

	// make coupon code invalid
	err2 := db.InsertOrUpdateCouponCode(false, coupon.DiscountPercent, coupon.CouponValue)
	if err2 != nil {
		http.Error(w, "Error Updating Coupon Code!", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Coupon Applied Successfully",
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return
}

func CalculateDiscountedPrice(originalPrice int, discountPercent int) int {
	discount := (originalPrice * discountPercent) / 100
	newPrice := originalPrice - discount
	return newPrice
}

func getUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userBalance, err := db.GetUserBalanceFromRedisDb(userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"balance": userBalance,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return
}
