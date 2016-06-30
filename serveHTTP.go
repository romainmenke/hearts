package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"limbo.services/trace"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func serveHTTP() {

	fmt.Println("server.http.starting")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.json", GetHeartJSON)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.svg", GetHeartSVG)
	router.HandleFunc("/user/{domain}/{user}.json", GetUserJSON)
	router.HandleFunc("/user/{domain}/{user}.svg", GetUserSVG)

	fmt.Println("server.http.ready")

	http.ListenAndServe(":8080", router)

}

func GetUserJSON(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "server.http.getUser.json")
	defer span.Close()

	vars := mux.Vars(r)
	domain := vars["domain"]
	name := vars["user"]

	user, err := db.LoadUser(ctx, domain, name)
	if err != nil || user == nil {
		resp := make(map[string]string)
		resp["message"] = "unknown user"
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(*user)
}

func GetUserSVG(w http.ResponseWriter, r *http.Request) {

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

	heart, err := db.LoadHeart(ctx, domain, user, repo)
	if err != nil || heart == nil {
		resp := make(map[string]string)
		resp["message"] = "unknown repository"
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(*heart)
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

	heart, err := db.LoadHeart(ctx, domain, user, repo)
	if err != nil {
		return
	}

	svgString := heart.SVG()

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
