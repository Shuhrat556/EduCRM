package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setActor(c *gin.Context, uid uuid.UUID, role string) {
	c.Set(ctxAuthUserID, uid)
	c.Set(ctxAuthRole, role)
}

func TestRequirePermission_table(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()

	tests := []struct {
		name       string
		role       string
		perm       rbac.Permission
		wantStatus int
		nextCalled bool
	}{
		{"admin users.manage", "admin", rbac.PermUsersManage, http.StatusOK, true},
		{"teacher denied users.manage", "teacher", rbac.PermUsersManage, http.StatusForbidden, false},
		{"teacher attendance", "teacher", rbac.PermAttendanceManage, http.StatusOK, true},
		{"student denied attendance manage", "student", rbac.PermAttendanceManage, http.StatusForbidden, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var next bool
			r := gin.New()
			r.GET("/t", func(c *gin.Context) {
				setActor(c, uid, tt.role)
				c.Next()
			}, RequirePermission(tt.perm), func(c *gin.Context) { next = true })
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/t", nil))
			if w.Code != tt.wantStatus {
				t.Fatalf("status %d, want %d", w.Code, tt.wantStatus)
			}
			if next != tt.nextCalled {
				t.Fatalf("next=%v, want %v", next, tt.nextCalled)
			}
		})
	}
}

func TestRequirePermission_missingActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/t", RequirePermission(rbac.PermUsersManage), func(c *gin.Context) {})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/t", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", w.Code)
	}
}

func TestRequireAnyPermission_studentPayments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()
	var next bool
	r := gin.New()
	r.GET("/t", func(c *gin.Context) {
		setActor(c, uid, "student")
		c.Next()
	}, RequireAnyPermission(rbac.PermPaymentsReadOwn, rbac.PermPaymentsStaff), func(c *gin.Context) { next = true })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/t", nil))
	if w.Code != http.StatusOK || !next {
		t.Fatalf("want 200 and next, got %d next=%v", w.Code, next)
	}
}
