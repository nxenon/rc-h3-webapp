package routes

import (
	"github.com/nxenon/rc-h3-webapp/middlewares"
	"net/http"
)

var routesWithHandlers = make(map[string]func(http.ResponseWriter, *http.Request))

func SetRoutes() {
	routesWithHandlers["/"] = IndexRouteHandler
	routesWithHandlers["/api/login_check"] = LoginCheckRouteHandler
	routesWithHandlers["/login"] = LoginFrontRouteHandler // Login front
	routesWithHandlers["/api/login"] = LoginRouteHandler
	routesWithHandlers["/api/captcha"] = CaptchaRouteHandler
	// Authenticated Routes
	routesWithHandlers["/products"] = middlewares.AuthMiddleware(ProductsFrontRouteHandler) // Products
	routesWithHandlers["/api/products"] = middlewares.AuthMiddleware(ProductsRouteHandler)
	routesWithHandlers["/api/products/add"] = middlewares.AuthMiddleware(addProductRouteHandler)
	routesWithHandlers["/api/product/remove"] = middlewares.AuthMiddleware(removeProductHandler)
	routesWithHandlers["/api/cart"] = middlewares.AuthMiddleware(myCartRouteHandler)
	routesWithHandlers["/cart"] = middlewares.AuthMiddleware(cartFrontRouteHandler) // Cart Front
	routesWithHandlers["/api/cart/apply_coupon"] = middlewares.AuthMiddleware(applyCouponRouteHandler)
	routesWithHandlers["/api/balance"] = middlewares.AuthMiddleware(getUserBalanceHandler)

	routesWithHandlers["/api/restart_all"] = middlewares.AuthMiddleware(restartAllRouteHandler)
}

func HandleRoutes(mux *http.ServeMux) {
	for key, value := range routesWithHandlers {
		mux.HandleFunc(key, value)
	}
}
