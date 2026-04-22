package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/usecase/auth"
	jwtpkg "github.com/educrm/educrm-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type httpTestUsers struct {
	byLogin map[string]*domain.User
	byID    map[uuid.UUID]*domain.User
}

func (s *httpTestUsers) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	if s.byLogin == nil {
		return nil, nil
	}
	return s.byLogin[login], nil
}
func (s *httpTestUsers) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}
func (s *httpTestUsers) Create(ctx context.Context, u *domain.User) error { return nil }
func (s *httpTestUsers) Update(ctx context.Context, u *domain.User) error { return nil }
func (s *httpTestUsers) Delete(ctx context.Context, id uuid.UUID) error   { return nil }
func (s *httpTestUsers) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (s *httpTestUsers) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *httpTestUsers) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *httpTestUsers) UsernameTaken(ctx context.Context, username string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}

type httpTestRefresh struct {
	byHash map[string]uuid.UUID
}

func (s *httpTestRefresh) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	if s.byHash == nil {
		s.byHash = make(map[string]uuid.UUID)
	}
	s.byHash[tokenHash] = userID
	return nil
}
func (s *httpTestRefresh) FindValidByHash(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	return uuid.Nil, repository.ErrNotFound
}
func (s *httpTestRefresh) DeleteByHash(ctx context.Context, tokenHash string) error { return nil }
func (s *httpTestRefresh) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (s *httpTestRefresh) Replace(ctx context.Context, userID uuid.UUID, oldHash, newHash string, expiresAt time.Time) error {
	return nil
}

func testJWTManager(t *testing.T) *jwtpkg.Manager {
	t.Helper()
	return jwtpkg.NewManager("handler-test-jwt-secret-32chars!!", time.Hour, "h-test")
}

func TestAuthHandler_Login_table(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte("password12345"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	users := &httpTestUsers{
		byLogin: map[string]*domain.User{
			"teacher@school.edu": {
				ID: uid, Email: strPtrHTTP("teacher@school.edu"), PasswordHash: string(hash),
				Role: domain.RoleTeacher, IsActive: true,
			},
		},
		byID: map[uuid.UUID]*domain.User{
			uid: {ID: uid, Email: strPtrHTTP("teacher@school.edu"), PasswordHash: string(hash), Role: domain.RoleTeacher, IsActive: true},
		},
	}
	refresh := &httpTestRefresh{byHash: map[string]uuid.UUID{}}
	svc := auth.NewService(users, refresh, testJWTManager(t), 24*time.Hour)
	h := NewAuthHandler(svc, false)

	tests := []struct {
		name       string
		body       any
		wantStatus int
		checkToken bool
	}{
		{
			name:       "success",
			body:       dto.LoginRequest{Login: "teacher@school.edu", Password: "password12345"},
			wantStatus: http.StatusOK,
			checkToken: true,
		},
		{
			name:       "wrong password",
			body:       dto.LoginRequest{Login: "teacher@school.edu", Password: "wrongpass9"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid json",
			body:       nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/login", h.Login)
			var req *http.Request
			if tt.body == nil {
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`{`)))
			} else {
				b, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status %d, want %d, body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
			if tt.checkToken {
				var env struct {
					Success bool `json:"success"`
					Data    struct {
						AccessToken string `json:"access_token"`
					} `json:"data"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
					t.Fatal(err)
				}
				if env.Data.AccessToken == "" {
					t.Fatal("missing access_token")
				}
			}
		})
	}
}

func TestAuthHandler_Me_withJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtMgr := testJWTManager(t)
	uid := uuid.New()
	users := &httpTestUsers{
		byID: map[uuid.UUID]*domain.User{
			uid: {ID: uid, Role: domain.RoleStudent, IsActive: true},
		},
	}
	svc := auth.NewService(users, &httpTestRefresh{}, jwtMgr, time.Hour)
	h := NewAuthHandler(svc, false)
	tok, err := jwtMgr.GenerateAccessToken(uid.String(), string(domain.RoleStudent))
	if err != nil {
		t.Fatal(err)
	}
	r := gin.New()
	r.GET("/me", middleware.AuthRequired(jwtMgr), h.Me)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
}

func strPtrHTTP(s string) *string { return &s }

func TestAuthHandler_Login_requirePortal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte("password12345"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	users := &httpTestUsers{
		byLogin: map[string]*domain.User{
			"teacher@school.edu": {
				ID: uid, Email: strPtrHTTP("teacher@school.edu"), PasswordHash: string(hash),
				Role: domain.RoleTeacher, IsActive: true,
			},
		},
		byID: map[uuid.UUID]*domain.User{
			uid: {ID: uid, Email: strPtrHTTP("teacher@school.edu"), PasswordHash: string(hash), Role: domain.RoleTeacher, IsActive: true},
		},
	}
	svc := auth.NewService(users, &httpTestRefresh{byHash: map[string]uuid.UUID{}}, testJWTManager(t), 24*time.Hour)
	h := NewAuthHandler(svc, true)

	r := gin.New()
	r.POST("/login", h.Login)

	t.Run("missing portal", func(t *testing.T) {
		b, _ := json.Marshal(dto.LoginRequest{Login: "teacher@school.edu", Password: "password12345"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("wrong portal", func(t *testing.T) {
		b, _ := json.Marshal(dto.LoginRequest{Login: "teacher@school.edu", Password: "password12345", Portal: "student"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("matching portal", func(t *testing.T) {
		b, _ := json.Marshal(dto.LoginRequest{Login: "teacher@school.edu", Password: "password12345", Portal: "teacher"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status %d body=%s", w.Code, w.Body.String())
		}
	})
}
