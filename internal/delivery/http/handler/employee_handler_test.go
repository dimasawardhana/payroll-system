package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestEmployeeLoginHandler_BadRequest(t *testing.T) {
	r := gin.Default()
	h := &EmployeeHandler{}
	r.POST("/login", h.EmployeeLoginHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
