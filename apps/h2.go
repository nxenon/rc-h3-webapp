package apps

import (
	"fmt"
	"github.com/nxenon/rc-h3-webapp/routes"
	"github.com/nxenon/rc-h3-webapp/utils"
	"net/http"
)

func StartHttp2Server(appData utils.AppData) {
	fmt.Printf("Satrting HTTP/2 Server on https://%s (TCP)\n", appData.H2ListenAddr)

	mux := http.NewServeMux()

	h2Server := http.Server{
		Handler: mux,
		Addr:    appData.H2ListenAddr,
	}

	routes.HandleRoutes(mux)

	h2Server.ListenAndServeTLS(appData.CertPath, appData.KeyPath)

}
