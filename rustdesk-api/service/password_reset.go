package service

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/utils"
)

const (
	passwordResetCodeExpireSeconds = 10 * 60
	passwordResetCooldownSeconds   = 60
)

type PasswordResetState struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Code   string `json:"code"`
	Used   bool   `json:"used"`
}

type PasswordResetService struct{}

func (prs *PasswordResetService) SendCode(email string) error {
	email = normalizeEmail(email)

	if email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}
	if !AllService.MailService.IsConfigured() {
		return errors.New("mail service not configured")
	}

	user, err := prs.findSingleUserByEmail(email)
	if err != nil {
		return err
	}

	var cooldown bool
	if err := global.Cache.Get(prs.cooldownKey(email), &cooldown); err == nil && cooldown {
		return errors.New("please wait 60 seconds before requesting another code")
	}

	code := utils.RandomNumericPassword(6)
	body := fmt.Sprintf(
		"Your RustDesk password reset code is: %s\n\nThis code expires in 10 minutes.\nIf you did not request a password reset, please ignore this email.",
		code,
	)
	if err := AllService.MailService.Send(user.Email, "RustDesk Password Reset Code", body); err != nil {
		return err
	}

	state := PasswordResetState{
		UserID: user.Id,
		Email:  normalizeEmail(user.Email),
		Code:   code,
		Used:   false,
	}
	if err := global.Cache.Set(prs.cacheKey(email), state, passwordResetCodeExpireSeconds); err != nil {
		return err
	}
	if err := global.Cache.Set(prs.cooldownKey(email), true, passwordResetCooldownSeconds); err != nil {
		Logger.Warnf("set password reset cooldown failed: %v", err)
	}
	return nil
}

func (prs *PasswordResetService) ResetWithCode(email, code, newPassword string) error {
	email = normalizeEmail(email)
	code = strings.TrimSpace(code)
	newPassword = strings.TrimSpace(newPassword)

	if email == "" || code == "" || newPassword == "" {
		return errors.New("email, code and new password are required")
	}
	if len(newPassword) < 6 || len(newPassword) > 18 {
		return errors.New("password length must be 6-18 characters")
	}

	user, err := prs.findSingleUserByEmail(email)
	if err != nil {
		return err
	}

	var state PasswordResetState
	if err := global.Cache.Get(prs.cacheKey(email), &state); err != nil {
		return errors.New("verification code expired or not requested")
	}
	if state.Used {
		return errors.New("verification code has already been used")
	}
	if state.Email != email || state.UserID != user.Id {
		return errors.New("verification code mismatch")
	}
	if state.Code != code {
		return errors.New("verification code is incorrect")
	}

	if err := AllService.UserService.UpdatePassword(user, newPassword); err != nil {
		return err
	}

	state.Used = true
	state.Code = ""
	if err := global.Cache.Set(prs.cacheKey(email), state, 300); err != nil {
		Logger.Warnf("mark password reset code used failed: %v", err)
	}
	return nil
}

func (prs *PasswordResetService) cacheKey(email string) string {
	return "password_reset_code:" + email
}

func (prs *PasswordResetService) cooldownKey(email string) string {
	return "password_reset_cooldown:" + email
}

func (prs *PasswordResetService) findSingleUserByEmail(email string) (*model.User, error) {
	var users []model.User
	if err := DB.Where("lower(email) = ?", email).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("no account is bound to this email")
	}
	if len(users) > 1 {
		return nil, errors.New("multiple accounts are bound to this email, please contact support")
	}
	return &users[0], nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
