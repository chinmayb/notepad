package root

import (
	"fmt"
	"net/http"

	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"

	"github.com/chinmayb/notepad/pkg/notepad"
)

// func HTTPError(body interface{}, err error, w http.ResponseWriter) {
// 	body = &errorutil.ErrorBody{
// 		Error: err.Error(),
// 	}
// 	buf, merr := json.Marshal(body)
// 	if merr != nil {
// 		log.Infof("Failed to marshal error message  %v", merr.Error())
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	_, err := w.Write(buf)
//
// }

func Register(addr string, port string) error {
	r := mux.NewRouter()
	h := &handler{notepad.NewNotePadServer(), zapr}

	r.PathPrefix("/api/v1/").Subrouter()
	han := Handler(r, h)
	srv := &http.Server{
		Handler: han,
		Addr:    fmt.Sprintf("%s:%s", addr, port),
	}
	return srv.ListenAndServe()
}

func Handler(router *mux.Router, server *handler) http.Handler {
	router.HandleFunc("notepad", server.createNotepad).Methods("POST")
	router.HandleFunc("notepads", server.listNotepad).Methods("GET")
	return router
}
