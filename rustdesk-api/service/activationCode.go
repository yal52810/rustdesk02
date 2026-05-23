package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type ActivationCodeService struct {
}

const activationCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const activationCodeLength = 13

func (s *ActivationCodeService) GenerateCode() string {
	code := make([]byte, activationCodeLength)
	max := big.NewInt(int64(len(activationCodeCharset)))
	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			code[i] = activationCodeCharset[i%len(activationCodeCharset)]
			continue
		}
		code[i] = activationCodeCharset[n.Int64()]
	}
	return string(code)
}

func (s *ActivationCodeService) Validate(code string) (*model.ActivationCode, error) {
	var ac model.ActivationCode
	err := DB.Preload("Package.Servers").Where("code = ?", code).First(&ac).Error
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

	return &ac, nil
}

func (s *ActivationCodeService) MarkUsed(ac *model.ActivationCode, userId uint) error {
	now := time.Now()
	ac.UsedBy = userId
	ac.UsedAt = &now
	return DB.Model(ac).Updates(map[string]interface{}{
		"used_by": userId,
		"used_at": now,
	}).Error
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

	query := DB.Model(&model.ActivationCode{}).Preload("Package").Preload("Package.Servers")
	query.Count(&total)

	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&codes).Error
	return codes, total, err
}

func (s *ActivationCodeService) Delete(id uint) error {
	return DB.Delete(&model.ActivationCode{}, id).Error
}

func (s *ActivationCodeService) ValidateAndUse(code string, userId uint) (*model.ActivationCode, error) {
	ac, err := s.Validate(code)
	if err != nil {
		return nil, err
	}
	if err := s.MarkUsed(ac, userId); err != nil {
		return nil, err
	}
	return ac, nil
}

func (s *ActivationCodeService) ResolveAssignedServers(ac *model.ActivationCode) (*uint, *uint) {
	if ac == nil {
		return nil, nil
	}
	if ac.PrimaryServerId != nil || ac.BackupServerId != nil {
		return ac.PrimaryServerId, ac.BackupServerId
	}
	if ac.Package == nil || len(ac.Package.Servers) == 0 {
		return nil, nil
	}

	servers := append([]*model.Server{}, ac.Package.Servers...)
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Priority != servers[j].Priority {
			return servers[i].Priority > servers[j].Priority
		}
		return servers[i].Id < servers[j].Id
	})

	primaryID := &servers[0].Id
	var backupID *uint
	if len(servers) > 1 {
		backupID = &servers[1].Id
	}
	return primaryID, backupID
}

func (s *ActivationCodeService) ApplyToUser(user *model.User, ac *model.ActivationCode, accumulate bool) map[string]interface{} {
	updates := map[string]interface{}{}
	if user == nil || ac == nil {
		return updates
	}

	if accumulate {
		if user.ValidDays != -1 {
			if ac.ValidDays == -1 {
				user.ValidDays = -1
			} else {
				user.ValidDays += ac.ValidDays
			}
			updates["valid_days"] = user.ValidDays
		}
		if ac.DeviceLimit > user.DeviceLimit {
			user.DeviceLimit = ac.DeviceLimit
			updates["device_limit"] = user.DeviceLimit
		}
	} else {
		user.ValidDays = ac.ValidDays
		user.DeviceLimit = ac.DeviceLimit
		updates["valid_days"] = user.ValidDays
		updates["device_limit"] = user.DeviceLimit
	}

	primaryID, backupID := s.ResolveAssignedServers(ac)
	if ac.PackageId != nil {
		user.PackageId = ac.PackageId
		updates["package_id"] = ac.PackageId
	}
	if primaryID != nil || ac.PackageId != nil {
		user.PrimaryServerId = primaryID
		user.RelayServerId = primaryID
		updates["primary_server_id"] = primaryID
		updates["relay_server_id"] = primaryID
	}
	if backupID != nil || ac.PackageId != nil {
		user.BackupServerId = backupID
		updates["backup_server_id"] = backupID
	}

	return updates
}
