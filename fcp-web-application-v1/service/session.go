package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"errors"
)

type SessionService interface {
	GetSessionByEmail(email string) (model.Session, error)
}

type sessionService struct {
	sessionRepo repo.SessionRepository
}

func NewSessionService(sessionRepo repo.SessionRepository) *sessionService {
	return &sessionService{sessionRepo}
}

func (c *sessionService) GetSessionByEmail(email string) (model.Session, error) {
	if email == "" {
        return model.Session{}, errors.New("email cannot be empty")
    }

    session, err := c.sessionRepo.SessionAvailEmail(email)
    if err != nil {
        return model.Session{}, errors.New("failed to get session by email")
    }

    return session, nil
}
