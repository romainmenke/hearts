package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"limbo.services/trace"
	"limbo.services/trace/dev"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/romainmenke/hearts/pkg/fakedb"
	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var (
	db *fakedb.FakeDB
)

func main() {

	serveR()

	serveH()

}

func serveR() {

	trace.DefaultHandler = dev.NewHandler(nil)

	db = fakedb.New("/go/src/app/db/", "/go/src/app/db/")

	fmt.Println("Starting Hearts")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	fmt.Printf("Listening on port : %s", port)
	fmt.Println("")

	s := grpc.NewServer()

	wercker.RegisterNotificationServiceServer(s, &server{})

	fmt.Println("Ready to serve clients")

	s.Serve(lis)

}

type server struct{}

func (s *server) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {

	if s == nil {
		return nil, errors.New("fubar")
	}

	if in == nil || in.Git == nil || in.Build == nil {
		return &wercker.WerckerResponse{Success: false}, errors.New("nil message")
	}

	heart, err := s.heart(ctx, in)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = db.SaveSVG(ctx, *heart)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = s.user(ctx, in, heart)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = db.Persist(ctx)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	return &wercker.WerckerResponse{Success: true}, nil
}

func (s *server) heart(ctx context.Context, message *wercker.WerckerMessage) (*fakedb.Heart, error) {

	span, ctx := trace.New(ctx, "Update Heart")
	defer span.Close()

	heart, err := db.LoadHeart(ctx, message.Git.Domain, message.Git.Owner, message.Git.Repository)
	if err != nil || heart == nil {
		span.Error(err)

		heart, err = s.newHeart(ctx, message)
		if err != nil {
			return nil, span.Error(err)
		}
		return heart, nil
	}

	if message.Result.Result == true && heart.LastBuild == false {
		heart.Count++
		if heart.Count > 3 {
			heart.Count = 3
		}
	} else if message.Result.Result == false {
		heart.Count--
		if heart.Count < 0 {
			heart.Count = 0
		}
	}

	heart.LastBuild = message.Result.Result

	err = db.SaveObject(ctx, heart)
	if err != nil {
		return nil, span.Error(err)
	}

	return heart, nil
}

func (s *server) newHeart(ctx context.Context, message *wercker.WerckerMessage) (*fakedb.Heart, error) {

	span, ctx := trace.New(ctx, "New Heart")
	defer span.Close()

	pass := message.Result.Result
	heart := &fakedb.Heart{
		Count:     3,
		LastBuild: pass,
		Domain:    message.Git.Domain,
		Owner:     message.Git.Owner,
		Repo:      message.Git.Repository,
	}

	err := db.SaveObject(ctx, heart)
	if err != nil {
		return nil, span.Error(err)
	}

	return heart, nil
}

func (s *server) user(ctx context.Context, message *wercker.WerckerMessage, heart *fakedb.Heart) error {
	userLog := fmt.Sprintf("Update User : %s", message.Build.User)
	span, ctx := trace.New(ctx, userLog)
	defer span.Close()

	user, err := db.LoadUser(ctx, message.Git.Domain, message.Build.User)
	if err != nil || user == nil {
		span.Error(err)

		err = s.newUser(ctx, message)
		if err != nil {
			return span.Error(err)
		}
		return nil
	}

	if message.Result.Result == true {
		user.Exp++
		user.Streak++
	} else {
		user.Streak = 0
	}

	if heart.Count == 0 {
		user.Deaths++
	}

	err = db.SaveObject(ctx, user)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func (s *server) newUser(ctx context.Context, message *wercker.WerckerMessage) error {

	span, ctx := trace.New(ctx, "New User")
	defer span.Close()

	pass := message.Result.Result
	user := &fakedb.User{
		Domain: message.Git.Domain,
		Name:   message.Build.User,
		Level:  0,
		Exp:    0,
		Streak: 0,
		Deaths: 0,
	}

	if pass {
		user.Exp++
		user.Streak++
	}

	err := db.SaveObject(ctx, user)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func serveH() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.json", GetUserJSON)
	router.HandleFunc("/heart/{domain}/{user}/{repo}.svg", GetUserSVG)
	router.HandleFunc("/user/{domain}/{user}.json", GetHeartJSON)
	router.HandleFunc("/user/{domain}/{user}.svg", GetHeartSVG)

	http.ListenAndServe(":8080", router)

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
