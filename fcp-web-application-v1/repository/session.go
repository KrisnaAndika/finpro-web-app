package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) *sessionsRepo {
	return &sessionsRepo{db}
}

func (u *sessionsRepo) AddSessions(session model.Session) error {
	var count int64
    u.db.Model(&model.Session{}).Where("token = ?", session.Token).Count(&count)
    if count > 0 {
        return errors.New("session already exists")
    }
    if result := u.db.Create(&session); result.Error != nil {
        return fmt.Errorf("failed to create session: %w", result.Error)
    }

    return nil
}

func (u *sessionsRepo) DeleteSession(token string) error {
	var count int64
    u.db.Model(&model.Session{}).Where("token = ?", token).Count(&count)
    if count == 0 {
        return errors.New("session not found")
    }

    if result := u.db.Where("token = ?", token).Delete(&model.Session{}); result.Error != nil {
        return fmt.Errorf("failed to delete session: %w", result.Error)
    }

    return nil
}

func (u *sessionsRepo) UpdateSessions(session model.Session) error {
	if result := u.db.Table("sessions").Where("email = ?", session.Email).Updates(session); result.Error != nil {
        return fmt.Errorf("failed to update session: %w", result.Error)
    }

    return nil
}

func (u *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	var session model.Session

    if err := u.db.Where("email = ?", email).First(&session).Error; err != nil {
        return session, fmt.Errorf("session unavailable: %w", err)
    }

    return session, nil
}

func (u *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session

    if err := u.db.Where("token = ?", token).First(&session).Error; err != nil {
        return model.Session{}, fmt.Errorf("session unavailable: %w", err)
    }

    return session, nil
}

func (u *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
