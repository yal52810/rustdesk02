package service

import (
	"errors"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

type DeviceSessionService struct {
	*BaseService
}

func (s *DeviceSessionService) Connect(userId uint, deviceId string) error {
	user := AllService.UserService.InfoById(userId)
	if user.Id == 0 {
		return errors.New("user not found")
	}

	s.CleanupInactiveSessions()

	var existingSession model.DeviceSession
	err := s.db.Where("user_id = ? AND device_id = ? AND status = ?",
		userId, deviceId, model.DeviceSessionStatusOnline).First(&existingSession).Error

	if err == nil {
		existingSession.LastActiveAt = time.Now()
		return s.db.Save(&existingSession).Error
	}

	activeCount := s.GetActiveSessionCount(userId)
	if activeCount >= user.DeviceLimit {
		return errors.New("device limit reached")
	}

	session := &model.DeviceSession{
		UserId:       userId,
		DeviceId:     deviceId,
		Status:       model.DeviceSessionStatusOnline,
		LastActiveAt: time.Now(),
	}
	return s.db.Create(session).Error
}

func (s *DeviceSessionService) Heartbeat(userId uint, deviceId string) error {
	return s.db.Model(&model.DeviceSession{}).
		Where("user_id = ? AND device_id = ? AND status = ?",
			userId, deviceId, model.DeviceSessionStatusOnline).
		Update("last_active_at", time.Now()).Error
}

func (s *DeviceSessionService) Disconnect(userId uint, deviceId string) error {
	return s.db.Model(&model.DeviceSession{}).
		Where("user_id = ? AND device_id = ?", userId, deviceId).
		Update("status", model.DeviceSessionStatusOffline).Error
}

func (s *DeviceSessionService) GetActiveSessionCount(userId uint) int {
	var count int64
	s.db.Model(&model.DeviceSession{}).
		Where("user_id = ? AND status = ? AND last_active_at > ?",
			userId, model.DeviceSessionStatusOnline, time.Now().Add(-60*time.Second)).
		Count(&count)
	return int(count)
}

func (s *DeviceSessionService) CleanupInactiveSessions() {
	s.db.Model(&model.DeviceSession{}).
		Where("status = ? AND last_active_at < ?",
			model.DeviceSessionStatusOnline, time.Now().Add(-60*time.Second)).
		Update("status", model.DeviceSessionStatusOffline)
}

func (s *DeviceSessionService) GetUserSessions(userId uint) ([]model.DeviceSession, error) {
	var sessions []model.DeviceSession
	err := s.db.Where("user_id = ? AND status = ?",
		userId, model.DeviceSessionStatusOnline).
		Order("last_active_at desc").Find(&sessions).Error
	return sessions, err
}
