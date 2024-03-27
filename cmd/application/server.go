package application

import "net/http"

func (a *App) initServer() {
	a.srv = &http.Server{
		Addr:         a.cfg.Port,
		ReadTimeout:  a.cfg.ReadTimeout,
		WriteTimeout: a.cfg.WriteTimeout,
		IdleTimeout:  a.cfg.IdleTimeout,
		Handler:      a.handler,
	}
}
