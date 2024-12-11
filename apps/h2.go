package apps

import (
	"fmt"
	"github.com/nxenon/rc-h3-webapp/routes"
	"github.com/nxenon/rc-h3-webapp/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

func StartHttp2Server(appData utils.AppData) {
	fmt.Printf("Satrting HTTP/2 Server on https://%s (TCP)\n", appData.H2ListenAddr)

	mux := http.NewServeMux()

	h2Server := &http2.Server{}

	server := &http.Server{
		Addr:    appData.H2ListenAddr,
		Handler: h2c.NewHandler(mux, h2Server),
	}

	routes.HandleRoutes(mux)

	err := server.ListenAndServeTLS(appData.CertPath, appData.KeyPath)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		panic(err)
	}

}
