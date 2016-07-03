package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"limbo.services/trace"

	"github.com/gorilla/mux"
	"github.com/romainmenke/travis"
	"golang.org/x/net/context"
)

func serveHTTP() {

	fmt.Println("server.http.starting")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.json", GetHeartJSON)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.svg", GetHeartSVG)
	router.HandleFunc("/user/{domain}/{user}.json", GetUserJSON)
	router.HandleFunc("/user/{domain}/{user}.svg", GetUserSVG)
	travis.HandleTravisWebHook(router, "/travis/", HandleTravisPayload)

	fmt.Println("server.http.ready")
	fmt.Println("server.tcp.listening on port : 8080")

	http.ListenAndServe(":8080", router)

}

func GetUserJSON(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "server.http.getUser.json")
	defer span.Close()

	vars := mux.Vars(r)
	domain := vars["domain"]
	name := vars["user"]

	userCache, err := cache.LoadUser(ctx, domain, name)
	if err != nil || userCache.User == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(*userCache.User)
	return
}

func GetUserSVG(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "not implemented", http.StatusNotImplemented)
	return

}

func GetHeartJSON(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "server.http.getHeart.json")
	defer span.Close()

	vars := mux.Vars(r)
	domain := vars["domain"]
	user := vars["user"]
	repo := vars["repo"]

	heartCache, err := cache.LoadHeart(ctx, domain, user, repo)
	if err != nil || heartCache.Heart == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(*heartCache.Heart)
}

func GetHeartSVG(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "server.http.getHeart.svg")
	defer span.Close()

	setResponseHeaderSVG(w)

	vars := mux.Vars(r)
	domain := vars["domain"]
	user := vars["user"]
	repo := vars["repo"]

	heartCache, err := cache.LoadHeart(ctx, domain, user, repo)
	if err != nil {
		return
	}

	svgString := heartCache.Heart.SVG()

	fmt.Fprint(w, svgString)

}

func setResponseHeaderSVG(w http.ResponseWriter) {

	cacheSince := time.Now().Format(http.TimeFormat)
	delay := time.Duration(30) * time.Second

	cacheUntil := time.Now().Add(delay).Format(http.TimeFormat)

	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Cache-Control, Expires, Etag, Last-Modified")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)

}

func HandleTravisPayload(payload *travis.PayloadObject) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "server.http.webhook")
	defer span.Close()

	if payload.Branch != "master" {
		span.Log("not on the main branch")
	}

	message := newFromTravis(payload)

	err := update(ctx, cache, message)
	if err != nil {
		span.Error(err)
	}
}
