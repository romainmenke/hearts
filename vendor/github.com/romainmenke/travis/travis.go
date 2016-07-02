package travis

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"limbo.services/trace"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var callBack func(payload *PayloadObject)

func HandleTravisWebHook(r *mux.Router, path string, f func(payload *PayloadObject)) {
	callBack = f
	r.HandleFunc(path, ReceiveTravis)
}

func ReceiveTravis(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	span, ctx := trace.New(ctx, "travis.receive")
	defer span.Close()

	var bodyObject BodyObject
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		span.Error(err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &bodyObject)
	if err != nil {
		span.Error(err)
	}

	if callBack == nil {
		span.Error("bad callback")
	}
	callBack(bodyObject.Payload)

}
