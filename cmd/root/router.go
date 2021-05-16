package root

import (
	"fmt"
	"net/http"

	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

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
	z, err := zap.NewProduction()
	if err != nil {
		return err
	}

	r.HandleFunc("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	}))
	logger := zapr.NewLogger(z)
	h := &handler{notepad.NewNotePadServer(), logger}
	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Notepad APIs"))
	}))
	han := Handler("/api/v1", r, h)
	srv := &http.Server{

		Handler: han,
		Addr:    fmt.Sprintf("%s:%s", addr, port),
	}

	return srv.ListenAndServe()
}

func Handler(prefix string, r *mux.Router, server *handler) http.Handler {
	r.HandleFunc(fmt.Sprintf("%s/notepads", prefix), server.listNotepad).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/notepad", prefix), server.createNotepad).Methods("POST")
	return r
}
