package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/educrm/educrm-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthRequired_table(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := jwt.NewManager("test-secret-for-jwt-middleware-32", time.Hour, "t")
	uid := uuid.New()
	valid, err := mgr.GenerateAccessToken(uid.String(), "teacher")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"not bearer", "Basic x", http.StatusUnauthorized},
		{"bad token", "Bearer not-a-jwt", http.StatusUnauthorized},
		{"valid", "Bearer " + valid, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotUID bool
			r := gin.New()
			r.GET("/p", AuthRequired(mgr), func(c *gin.Context) {
				_, gotUID = UserID(c)
			})
			req := httptest.NewRequest(http.MethodGet, "/p", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && !gotUID {
				t.Fatal("expected user id in context")
			}
		})
	}
}
