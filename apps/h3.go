package apps

import (
	"crypto/tls"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/routes"
	"github.com/nxenon/rc-h3-webapp/utils"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func StartHttp3Server(appData utils.AppData) {
	fmt.Printf("Starting HTTP/3 Server on https://%s (UDP)\n", appData.H3ListenAddr)

	// LOG FILE
	//file, err := os.OpenFile(AppData.KeyLogFile, os.O_APPEND|os.O_CREATE, 0600)
	//if err != nil {
	//	panic(err)
	//}

	cert, err := tls.LoadX509KeyPair(
		appData.CertPath,
		appData.KeyPath,
	)
	if err != nil {
		panic(err)
	}
	quicConf := &quic.Config{}

	mux := http.NewServeMux()
	server := http3.Server{
		Handler: mux,
		Addr:    appData.H2ListenAddr,
		TLSConfig: http3.ConfigureTLSConfig(&tls.Config{
			//KeyLogWriter: file,
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h3"},
		}),
		QUICConfig: quicConf,
	}

	routes.HandleRoutes(mux)

	server.ListenAndServe()

}
