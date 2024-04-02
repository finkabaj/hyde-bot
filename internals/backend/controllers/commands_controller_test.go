package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCommandsControllerRoutes(t *testing.T) {
	r := chi.NewRouter()
	controller := CommandsController{}
	controller.RegisterRoutes(r)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expect := "Hello world"

	if rr.Body.String() != expect {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expect)
	}
}
