package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type ActivationCodeService struct {
}

func (s *ActivationCodeService) GenerateCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *ActivationCodeService) Create(packageId *uint, validDays, deviceLimit int, expiresAt *time.Time, remark string) (*model.ActivationCode, error) {
	code := &model.ActivationCode{
		Code:        s.GenerateCode(),
		PackageId:   packageId,
		ValidDays:   validDays,
		DeviceLimit: deviceLimit,
		ExpiresAt:   expiresAt,
		Remark:      remark,
	}
	err := DB.Create(code).Error
	return code, err
}

func (s *ActivationCodeService) BatchCreate(packageId *uint, count, validDays, deviceLimit int, expiresAt *time.Time, remark string) ([]*model.ActivationCode, error) {
	codes := make([]*model.ActivationCode, count)
	for i := 0; i < count; i++ {
		codes[i] = &model.ActivationCode{
			Code:        s.GenerateCode(),
			PackageId:   packageId,
			ValidDays:   validDays,
			DeviceLimit: deviceLimit,
			ExpiresAt:   expiresAt,
			Remark:      remark,
		}
	}
	err := DB.Create(&codes).Error
	return codes, err
}

func (s *ActivationCodeService) List(page, pageSize int) ([]model.ActivationCode, int64, error) {
	var codes []model.ActivationCode
	var total int64

	query := DB.Model(&model.ActivationCode{})
	query.Count(&total)

	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&codes).Error
	return codes, total, err
}

func (s *ActivationCodeService) Delete(id uint) error {
	return DB.Delete(&model.ActivationCode{}, id).Error
}

func (s *ActivationCodeService) ValidateAndUse(code string, userId uint) (*model.ActivationCode, error) {
	var ac model.ActivationCode
	err := DB.Where("code = ?", code).First(&ac).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("activation code not found")
		}
		return nil, err
	}

	if ac.UsedBy > 0 {
		return nil, errors.New("activation code already used")
	}

	if ac.ExpiresAt != nil && ac.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("activation code expired")
	}

	now := time.Now()
	ac.UsedBy = userId
	ac.UsedAt = &now
	err = DB.Save(&ac).Error
	return &ac, err
}
