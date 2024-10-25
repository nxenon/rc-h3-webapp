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

// ProductsFrontRouteHandler Handler for /products which is front of products
func ProductsFrontRouteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/products.html")
	AppData := utils.LoadEnvFile(".env")
	subStrings := strings.Split(AppData.H3ListenAddr, ":")
	subStringsHost := strings.Split(r.Host, ":")
	altSvcHeaderValue := fmt.Sprintf("h3=\"%s:%s\"", subStringsHost[0], subStrings[1])
	w.Header().Set("Alt-Svc", altSvcHeaderValue)
}

// ProductsRouteHandler Handler for /api/products
func ProductsRouteHandler(w http.ResponseWriter, r *http.Request) {
	allProducts, err := db.GetAllProducts()

	if err != nil {
		http.Error(w, "Error loading products!", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(allProducts)
	if err != nil {
		http.Error(w, "Error converting products to json!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func addProductRouteHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Call AddProductToUserCart
	if err := db.AddProductToUserCart(userId, req.ProductId); err != nil {
		http.Error(w, "Failed to add product to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product added to cart successfully"))
}

func removeProductHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.RemoveProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// remove product from user cart by models.RemoveProductRequest.ProductInCartId
	err2 := db.RemoveProductFromCartByUserId(req.ProductInCartUUID, userId)
	if err2 != nil {
		http.Error(w, "Error removing product from cart!", http.StatusInternalServerError)
		return
	}

	// update cart overall price
	//db.UpdateCartById(userCartObject.CartId, userCartObject.CartOverallPrice-userProduct.ProductPrice)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product removed from cart successfully"))
}
