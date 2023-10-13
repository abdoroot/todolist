package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsAuthApiUser(t *testing.T) {
	// إعداد محيط الاختبار
	r := gin.Default()
	r.Use(IsAuthApiUser())

	r.GET("/secure-endpoint", func(c *gin.Context) {
		c.String(http.StatusOK, "Authorized")
	})

	// اختبار حالة التوكن الصحيحة
	w := performRequest(r, "GET", "/secure-endpoint", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTczMTc4NzksInVzZXJfaWQiOjF9.AFyhx8nYt5ZSrsfARn0klk6y7GM0rPCFPcY_cj3OeK4")
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// اختبار حالة التوكن غير صحيحة
	w = performRequest(r, "GET", "/secure-endpoint", "Bearer invalid_token")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, w.Code)
	}
}

func performRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
