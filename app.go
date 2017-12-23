package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/signcontrol"
	"gitlab.com/swarmfund/hgate/config"
	"gitlab.com/swarmfund/hgate/server"
	"gitlab.com/swarmfund/hgate/tx"
)

type App struct {
	Conf *config.GateConfig
	Sub  *tx.Submitter
	Log  *logan.Entry
}

var proxy *httputil.ReverseProxy

func NewApp(configPath string) (app *App, err error) {
	app = new(App)

	app.Conf, err = config.InitConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("can not initialize config: %s", err.Error())
	}

	app.Log = logan.New().Level(app.Conf.LogLevel).
		WithField("application", "hgate")

	app.Sub, err = tx.NewSubmitter(app.Conf.RawHorizonURL, app.Conf.AccountID, app.Conf.KP)
	if err != nil {
		app.Log.WithError(err).Error("can not initialize submitter")
		return nil, err
	}

	return
}

func (app *App) Serve() {
	proxy = newProxy(app.Conf.HorizonURL)
	mux := server.NewMux()
	mux.Post("/send_payment", initPaymentHandler(app))
	mux.HandleFunc("/", app.redirectHandler)

	fmt.Println("HGate server listening at:" + app.Conf.Port)
	err := http.ListenAndServe("localhost:"+app.Conf.Port, mux)
	log.Fatal(err)
}

func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	lEntry := app.Log.WithField("service", "redirect_handler")
	lEntry.WithField("path", r.URL.Path).Info("Started request")

	err := signcontrol.SignRequest(r, app.Conf.KP)
	if err != nil {
		lEntry.WithError(err).Error("SignRequest failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
	lEntry.WithField("path", r.URL.Path).Info("Finished request")
}

func newProxy(targetURL *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	superDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		superDirector(req)
		req.Host = targetURL.Host
		req.Header.Set("TestStuff", targetURL.Host)
	}
	return proxy
}
