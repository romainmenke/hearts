package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
)

func serveHTTP() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.json", GetUserJSON)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.svg", GetUserSVG)
	router.HandleFunc("/user/{domain}/{user}.json", GetHeartJSON)
	router.HandleFunc("/user/{domain}/{user}.svg", GetHeartSVG)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetUserJSON(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	vars := mux.Vars(r)
	domain := vars["domain"]
	name := vars["name"]

	user, err := db.LoadUser(ctx, domain, name)
	if err != nil {
		resp := make(map[string]string)
		resp["message"] = "unknown user"
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(user)
}

func GetUserSVG(w http.ResponseWriter, r *http.Request) {

	return

}

func GetHeartJSON(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	vars := mux.Vars(r)
	domain := vars["domain"]
	user := vars["user"]
	repo := vars["repo"]

	heart, err := db.LoadHeart(ctx, domain, user, repo)
	if err != nil {
		resp := make(map[string]string)
		resp["message"] = "unknown repository"
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(heart)
}

func GetHeartSVG(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

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

	w.Header().Add("Content-Type", "image/svg")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Cache-Control, Expires, Etag, Last-Modified")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")

}
