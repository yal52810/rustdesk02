package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type PackageService struct {
	*BaseService
}

// List 获取套餐列表
func (s *PackageService) List(page, pageSize uint, where func(tx *gorm.DB)) (res *model.PackageList) {
	res = &model.PackageList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := s.db.Model(&model.Package{}).Preload("Servers")
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Order("priority DESC, id ASC").Find(&res.Packages)
	return
}

// Create 创建套餐
func (s *PackageService) Create(pkg *model.Package, serverIds []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 如果设置为新用户默认套餐，先清除旧的默认标记
		if pkg.IsDefaultNewUser {
			if err := tx.Model(&model.Package{}).Where("is_default_new_user = ?", true).
				Update("is_default_new_user", false).Error; err != nil {
				return err
			}
		}
		// 创建套餐
		if err := tx.Create(pkg).Error; err != nil {
			return err
		}
		// 关联服务器
		if len(serverIds) > 0 {
			var servers []*model.Server
			if err := tx.Where("id IN ?", serverIds).Find(&servers).Error; err != nil {
				return err
			}
			if err := tx.Model(pkg).Association("Servers").Append(servers); err != nil {
				return err
			}
		}
		return nil
	})
}

// Update 更新套餐
func (s *PackageService) Update(pkg *model.Package, serverIds []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 如果设置为新用户默认套餐，先清除旧的默认标记
		if pkg.IsDefaultNewUser {
			if err := tx.Model(&model.Package{}).Where("is_default_new_user = ? AND id != ?", true, pkg.Id).
				Update("is_default_new_user", false).Error; err != nil {
				return err
			}
		}
		// 更新套餐基本信息
		if err := tx.Model(pkg).Updates(pkg).Error; err != nil {
			return err
		}
		// 更新关联的服务器
		if serverIds != nil {
			// 清除旧关联
			if err := tx.Model(pkg).Association("Servers").Clear(); err != nil {
				return err
			}
			// 添加新关联
			if len(serverIds) > 0 {
				var servers []*model.Server
				if err := tx.Where("id IN ?", serverIds).Find(&servers).Error; err != nil {
					return err
				}
				if err := tx.Model(pkg).Association("Servers").Append(servers); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Delete 删除套餐
func (s *PackageService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		pkg := &model.Package{}
		pkg.Id = id
		// 清除服务器关联
		if err := tx.Model(pkg).Association("Servers").Clear(); err != nil {
			return err
		}
		// 删除套餐
		return tx.Delete(pkg).Error
	})
}

// GetById 根据ID获取套餐
func (s *PackageService) GetById(id uint) (*model.Package, error) {
	var pkg model.Package
	err := s.db.Preload("Servers").First(&pkg, id).Error
	return &pkg, err
}

// GetActivePackages 获取所有启用的套餐
func (s *PackageService) GetActivePackages() ([]*model.Package, error) {
	var packages []*model.Package
	err := s.db.Where("is_active = ?", true).Preload("Servers").Order("priority DESC, id ASC").Find(&packages).Error
	return packages, err
}

// GetDefaultNewUserPackage 获取新用户注册默认套餐
func (s *PackageService) GetDefaultNewUserPackage() (*model.Package, error) {
	var pkg model.Package
	err := s.db.Where("is_default_new_user = ? AND is_active = ?", true, true).
		Preload("Servers").First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// ClearDefaultNewUserFlag 清除所有套餐的新用户默认标记
func (s *PackageService) ClearDefaultNewUserFlag() error {
	return s.db.Model(&model.Package{}).Where("is_default_new_user = ?", true).
		Update("is_default_new_user", false).Error
}

// GetPackageServers 获取套餐的服务器列表
func (s *PackageService) GetPackageServers(packageId uint) ([]*model.Server, error) {
	var pkg model.Package
	if err := s.db.Preload("Servers").First(&pkg, packageId).Error; err != nil {
		return nil, err
	}
	return pkg.Servers, nil
}
