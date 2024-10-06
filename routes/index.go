package routes

import (
	"fmt"
	"github.com/nxenon/rc-h3-webapp/utils"
	"net/http"
	"strings"
)

func IndexRouteHandler(w http.ResponseWriter, r *http.Request) {
	AppData := utils.LoadEnvFile(".env")
	subStrings := strings.Split(AppData.H3ListenAddr, ":")
	subStringsHost := strings.Split(r.Host, ":")
	altSvcHeaderValue := fmt.Sprintf("h3=\"%s:%s\"", subStringsHost[0], subStrings[1])
	w.Header().Set("Alt-Svc", altSvcHeaderValue)

	http.ServeFile(w, r, "templates/index.html")

}
