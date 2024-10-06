package middlewares

import (
	"context"
	"github.com/nxenon/rc-h3-webapp/utils"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("jwtToken")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		jwt_claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", jwt_claims.UserId)
		ctx2 := context.WithValue(ctx, "userName", jwt_claims.Username)
		r = r.WithContext(ctx2)
		//username, ok := r.Context().Value("username").(string)

		next(w, r)
	}
}
