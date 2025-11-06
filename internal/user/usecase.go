package user

import (
    "time"

    "majoo-case1-rest-api/config"
    "majoo-case1-rest-api/internal/security"
)

type Usecase struct {
    repo   *Repository
    secret []byte
}

func NewUsecase(repo *Repository, cfg config.Config) *Usecase {
    return &Usecase{repo: repo, secret: []byte(cfg.JWTSecret)}
}

func (u *Usecase) Register(username, email, password string) (User, string, error) {
    exists, err := u.repo.ExistsByEmailOrUsername(email, username)
    if err != nil {
        return User{}, "", err
    }
    if exists {
        return User{}, "", ErrConflict
    }
    hash, err := security.HashPassword(password)
    if err != nil {
        return User{}, "", err
    }
    id, err := u.repo.Create(username, email, hash)
    if err != nil {
        return User{}, "", err
    }
    token, err := security.GenerateToken(u.secret, id, username, email, 24*time.Hour)
    if err != nil {
        return User{}, "", err
    }
    return User{ID: id, Username: username, Email: email}, token, nil
}

func (u *Usecase) Login(email, password string) (User, string, error) {
    user, err := u.repo.GetByEmail(email)
    if err != nil {
        return User{}, "", err
    }
    if !security.CheckPasswordHash(password, user.PasswordHash) {
        return User{}, "", ErrUnauthorized
    }
    token, err := security.GenerateToken(u.secret, user.ID, user.Username, user.Email, 24*time.Hour)
    if err != nil {
        return User{}, "", err
    }
    return user, token, nil
}

var (
    ErrUnauthorized = fmtErr("unauthorized")
    ErrConflict     = fmtErr("conflict")
)

type fmtErr string

func (e fmtErr) Error() string { return string(e) }


