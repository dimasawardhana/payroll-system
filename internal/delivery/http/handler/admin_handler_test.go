package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminLoginHandler_BadRequest(t *testing.T) {
	r := gin.Default()
	h := &AdminHandler{}
	r.POST("/login", h.AdminLoginHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
