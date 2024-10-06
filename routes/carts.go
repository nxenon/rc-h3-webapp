package routes

import (
	"encoding/json"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/db"
	"github.com/nxenon/rc-h3-webapp/models"
	"github.com/nxenon/rc-h3-webapp/utils"
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

	userCartId, err := db.GetCartIdByUserId(userId)
	if err != nil {
		http.Error(w, "Error getting Cart ID by User ID!", http.StatusInternalServerError)
		return
	}

	userCart, err := db.GetCartById(userCartId)
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
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get user ID (assumed to be set in context or session)
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	coupon, err := db.GetCouponByValue(request.CouponValue)
	if err != nil {
		http.Error(w, "Invalid Coupon Code!", http.StatusNotFound)
		return
	}

	if coupon.IsValid != 1 {
		http.Error(w, "Invalid Coupon Code!", http.StatusNotFound)
		return
	}

	userCartId, err := db.GetCartIdByUserId(userId)
	if err != nil {
		http.Error(w, "Error getting user cart!", http.StatusInternalServerError)
		return
	}
	userCartObject, err4 := db.GetCartById(userCartId)
	if err4 != nil {
		http.Error(w, "Error getting user cart by id!", http.StatusInternalServerError)
		return
	}

	// update user cart
	newPrice := CalculateDiscountedPrice(userCartObject.CartOverallPrice, coupon.DiscountPercent)
	err3 := db.UpdateCartById(userCartId, newPrice)
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
	userObject, err := db.GetUserObjectByUserId(userId)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"balance": userObject.UserBalance,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
	return
}
