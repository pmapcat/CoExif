package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
)

type RESTHandler struct {
	Dispatch Dispatcher
	AllowURL string
}

func (h *RESTHandler) GETMeta(w rest.ResponseWriter, r *rest.Request) {
	parsed_url, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		rest.Error(w, "Malformed url: "+r.RequestURI, 400)
	}

	filter_tags := parsed_url.Query()["tags"]
	if filter_tags == nil {
		filter_tags = make([]string, 0)
	}
	fpath := path.Join(h.AllowURL, parsed_url.Path)
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		rest.Error(w, "no such file or directory: "+fpath, 400)
		return
	}
	datum := h.Dispatch.GETMeta(fpath, filter_tags...)
	w.WriteJson(&datum)
}

func (h *RESTHandler) POSTMeta(w rest.ResponseWriter, r *rest.Request) {
	parsed_url, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		rest.Error(w, "Malformed url: "+r.RequestURI, 400)
	}

	fpath := path.Join(h.AllowURL, parsed_url.Path)
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		rest.Error(w, "no such file or directory: "+fpath, 400)
		return
	}

	// if error on path
	stat, err := os.Stat(fpath)
	if err != nil {
		rest.Error(w, err.Error(), 400)
		return
	}
	// if posts in directory
	if stat.IsDir() {
		rest.Error(w, "Posting to directories is not supported: "+fpath, 400)
		return
	}

	payload := StandartJSON{}
	err = r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), 400)
	}
	err = h.Dispatch.POSTMeta(fpath, payload)
	if err != nil {
		rest.Error(w, err.Error(), 400)
	} else {
		h.GETMeta(w, r)
	}
}
func (h *RESTHandler) Destroy() {
	h.Dispatch.Exit()
	os.Exit(1)
}

func (h *RESTHandler) Run(port, uid, pass, allow_url string) {
	h.Dispatch.Init()
	h.AllowURL = allow_url
	// =========  Handle interrupt ==========
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		// kill all subprocesses
		h.Destroy()
		os.Exit(1)
	}()
	// ========= End handle interrupt ========
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.AuthBasicMiddleware{
		Realm: "basic auth",
		Authenticator: func(userId string, password string) bool {
			if userId == uid && password == pass {
				return true
			}
			return false
		},
	})
	router, err := rest.MakeRouter(
		rest.Get("/", h.GETMeta),
		rest.Get("/*", h.GETMeta),
		rest.Post("/*", h.POSTMeta),
	)
	if err != nil {
		log.Println(err)
	}
	api.SetApp(router)
	log.Println(http.ListenAndServe(":"+port, api.MakeHandler()))
	// kill all subprocesses
	h.Dispatch.Exit()
}
