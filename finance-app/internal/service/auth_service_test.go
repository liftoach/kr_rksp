package service

import (
	"context"
	"errors"
	"kr/internal/domain"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	user *domain.User
	err  error
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.user = user
	return m.err
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockJWT struct{}

func (m *mockJWT) Generate(id uuid.UUID) (string, error) {
	return "token", nil
}

func FuzzService_Register(f *testing.F) {
	f.Add("test@example.com", "password123")
	f.Add("", "")
	f.Add("a@a.a", "123")

	f.Fuzz(func(t *testing.T, email, password string) {
		repo := &mockUserRepository{}

		svc := &Service{
			UserRepository: repo,
		}

		err := svc.Register(context.Background(), email, password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if repo.user == nil {
			t.Fatalf("user was not created")
		}

		if repo.user.Email != email {
			t.Fatalf("expected email %q, got %q", email, repo.user.Email)
		}

		if repo.user.Password == password {
			t.Fatalf("password was not hashed")
		}

		if bcrypt.CompareHashAndPassword(
			[]byte(repo.user.Password),
			[]byte(password),
		) != nil {
			t.Fatalf("invalid password hash")
		}
	})
}

func FuzzService_Login(f *testing.F) {
	f.Add("test@example.com", "password123")
	f.Add("", "")
	f.Add("admin", "admin")

	f.Fuzz(func(t *testing.T, email, password string) {
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			t.Fatalf("bcrypt error: %v", err)
		}

		repo := &mockUserRepository{
			user: &domain.User{
				ID:       uuid.New(),
				Email:    email,
				Password: string(hash),
			},
		}

		svc := &Service{
			UserRepository: repo,
			jwt:            &mockJWT{},
		}

		token, err := svc.Login(context.Background(), email, password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token == "" {
			t.Fatalf("expected token")
		}
	})
}

func FuzzService_Login_InvalidPassword(f *testing.F) {
	f.Add("password123", "wrongpassword")

	f.Fuzz(func(t *testing.T, correctPassword, wrongPassword string) {
		if correctPassword == wrongPassword {
			t.Skip()
		}

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(correctPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			t.Fatalf("bcrypt error: %v", err)
		}

		repo := &mockUserRepository{
			user: &domain.User{
				ID:       uuid.New(),
				Email:    "test@example.com",
				Password: string(hash),
			},
		}

		svc := &Service{
			UserRepository: repo,
			jwt:            &mockJWT{},
		}

		_, err = svc.Login(
			context.Background(),
			"test@example.com",
			wrongPassword,
		)

		if err == nil {
			t.Fatalf("expected invalid password error")
		}

		if !errors.Is(err, errors.New("invalid password")) &&
			err.Error() != "invalid password" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
