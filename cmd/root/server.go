package root

import (
	"encoding/json"
	"net/http"

	"github.com/chinmayb/notepad/pkg/notepad"
	"github.com/go-logr/logr"
)

type handler struct {
	notepad.NotePadServer
	log logr.Logger
}

func (s *handler) createNotepad(w http.ResponseWriter, r *http.Request) {
	in := &notepad.NotePad{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	out, err := s.Create(ctx, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(b); err != nil {
		s.log.Error(err, "error while writing")
	}
}

func (s *handler) listNotepad(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	out, err := s.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(b); err != nil {
		s.log.Error(err, "error while writing")
	}
}
