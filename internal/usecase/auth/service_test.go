package auth

import (
	"context"
	"testing"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	jwtpkg "github.com/educrm/educrm-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type stubUserRepo struct {
	byLogin        map[string]*domain.User
	byID           map[uuid.UUID]*domain.User
	findByLoginErr error
}

func (s *stubUserRepo) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	if s.findByLoginErr != nil {
		return nil, s.findByLoginErr
	}
	if s.byLogin == nil {
		return nil, nil
	}
	return s.byLogin[login], nil
}

func (s *stubUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

func (s *stubUserRepo) Create(ctx context.Context, u *domain.User) error                     { return nil }
func (s *stubUserRepo) Update(ctx context.Context, u *domain.User) error                     { return nil }
func (s *stubUserRepo) Delete(ctx context.Context, id uuid.UUID) error                       { return nil }
func (s *stubUserRepo) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserRepo) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *stubUserRepo) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}

type stubRefreshRepo struct {
	byHash map[string]uuid.UUID
}

func (s *stubRefreshRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	if s.byHash == nil {
		s.byHash = make(map[string]uuid.UUID)
	}
	s.byHash[tokenHash] = userID
	return nil
}

func (s *stubRefreshRepo) FindValidByHash(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	if s.byHash == nil {
		return uuid.Nil, repository.ErrNotFound
	}
	uid, ok := s.byHash[tokenHash]
	if !ok {
		return uuid.Nil, repository.ErrNotFound
	}
	return uid, nil
}

func (s *stubRefreshRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	if s.byHash != nil {
		delete(s.byHash, tokenHash)
	}
	return nil
}

func (s *stubRefreshRepo) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	if s.byHash == nil {
		return nil
	}
	for h, u := range s.byHash {
		if u == userID {
			delete(s.byHash, h)
		}
	}
	return nil
}

func (s *stubRefreshRepo) Replace(ctx context.Context, userID uuid.UUID, oldHash, newHash string, expiresAt time.Time) error {
	if s.byHash == nil {
		return repository.ErrNotFound
	}
	if _, ok := s.byHash[oldHash]; !ok {
		return repository.ErrNotFound
	}
	delete(s.byHash, oldHash)
	s.byHash[newHash] = userID
	return nil
}

func testJWT(t *testing.T) *jwtpkg.Manager {
	t.Helper()
	return jwtpkg.NewManager("test-secret-key-min-32-chars-ok", time.Hour, "educrm-test")
}

func TestLogin_table(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	active := &domain.User{
		ID: uid, Email: strPtr("u@test.edu"), PasswordHash: string(hash),
		Role: domain.RoleTeacher, IsActive: true,
	}
	inactive := &domain.User{
		ID: uuid.New(), Email: strPtr("off@test.edu"), PasswordHash: string(hash),
		Role: domain.RoleStudent, IsActive: false,
	}

	tests := []struct {
		name      string
		users     *stubUserRepo
		login     string
		password  string
		wantKind  apperror.Kind
		wantOK    bool
	}{
		{
			name: "success",
			users: &stubUserRepo{
				byLogin: map[string]*domain.User{"u@test.edu": active},
				byID:    map[uuid.UUID]*domain.User{uid: active},
			},
			login:    "u@test.edu",
			password: "correct-password",
			wantOK:   true,
		},
		{
			name: "wrong password",
			users: &stubUserRepo{
				byLogin: map[string]*domain.User{"u@test.edu": active},
			},
			login:     "u@test.edu",
			password:  "wrong",
			wantKind:  apperror.KindUnauthorized,
		},
		{
			name:     "unknown user",
			users:    &stubUserRepo{byLogin: map[string]*domain.User{}},
			login:    "nobody@test.edu",
			password: "x",
			wantKind: apperror.KindUnauthorized,
		},
		{
			name: "inactive",
			users: &stubUserRepo{
				byLogin: map[string]*domain.User{"off@test.edu": inactive},
			},
			login:     "off@test.edu",
			password:  "correct-password",
			wantKind:  apperror.KindUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refresh := &stubRefreshRepo{byHash: map[string]uuid.UUID{}}
			svc := NewService(tt.users, refresh, testJWT(t), 24*time.Hour)
			pair, err := svc.Login(ctx, tt.login, tt.password)
			if tt.wantOK {
				if err != nil {
					t.Fatal(err)
				}
				if pair.AccessToken == "" || pair.RefreshToken == "" {
					t.Fatal("expected tokens")
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			ae, ok := apperror.AsError(err)
			if !ok || ae.Kind != tt.wantKind {
				t.Fatalf("got %v (%T), want kind %s", err, err, tt.wantKind)
			}
		})
	}
}

func TestRefresh_table(t *testing.T) {
	ctx := context.Background()
	jwtMgr := testJWT(t)
	uid := uuid.New()
	u := &domain.User{ID: uid, Role: domain.RoleAdmin, IsActive: true}

	tests := []struct {
		name     string
		raw      string
		refresh  *stubRefreshRepo
		users    *stubUserRepo
		wantKind apperror.Kind
		wantOK   bool
	}{
		{
			name:     "empty",
			raw:      "  ",
			wantKind: apperror.KindValidation,
		},
		{
			name:     "unknown token",
			raw:      "deadbeef",
			refresh:  &stubRefreshRepo{byHash: map[string]uuid.UUID{}},
			wantKind: apperror.KindUnauthorized,
		},
		{
			name: "success rotates",
			raw:  "opaque-refresh-token",
			refresh: func() *stubRefreshRepo {
				h := hashRefreshToken("opaque-refresh-token")
				return &stubRefreshRepo{byHash: map[string]uuid.UUID{h: uid}}
			}(),
			users: &stubUserRepo{byID: map[uuid.UUID]*domain.User{uid: u}},
			wantOK: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.users == nil {
				tt.users = &stubUserRepo{}
			}
			if tt.refresh == nil {
				tt.refresh = &stubRefreshRepo{}
			}
			svc := NewService(tt.users, tt.refresh, jwtMgr, 24*time.Hour)
			pair, err := svc.Refresh(ctx, tt.raw)
			if tt.wantOK {
				if err != nil {
					t.Fatal(err)
				}
				if pair.AccessToken == "" || pair.RefreshToken == "" {
					t.Fatal("expected new pair")
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			ae, ok := apperror.AsError(err)
			if !ok || ae.Kind != tt.wantKind {
				t.Fatalf("got %v, want kind %s", err, tt.wantKind)
			}
		})
	}
}

func TestLogout_BearerRevokesAll(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	jwtMgr := testJWT(t)
	tok, err := jwtMgr.GenerateAccessToken(uid.String(), string(domain.RoleStudent))
	if err != nil {
		t.Fatal(err)
	}
	h := hashRefreshToken("r1")
	refresh := &stubRefreshRepo{byHash: map[string]uuid.UUID{h: uid}}
	svc := NewService(&stubUserRepo{}, refresh, jwtMgr, time.Hour)
	if err := svc.Logout(ctx, "Bearer "+tok, ""); err != nil {
		t.Fatal(err)
	}
	if len(refresh.byHash) != 0 {
		t.Fatalf("expected sessions cleared, still %d", len(refresh.byHash))
	}
}

func TestMe_notFound(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	svc := NewService(
		&stubUserRepo{byID: map[uuid.UUID]*domain.User{}},
		&stubRefreshRepo{},
		testJWT(t),
		time.Hour,
	)
	_, err := svc.Me(ctx, uid)
	if err == nil {
		t.Fatal("expected error")
	}
	ae, ok := apperror.AsError(err)
	if !ok || ae.Kind != apperror.KindNotFound {
		t.Fatalf("got %v", err)
	}
}

func strPtr(s string) *string { return &s }
