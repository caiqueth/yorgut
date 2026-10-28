package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserFanRepository struct {
	repository[*model.UserFan]
}

func NewUserFanRepository(Db *gorm.DB) *UserFanRepository {
	return &UserFanRepository{
		repository: repository[*model.UserFan]{db: Db},
	}
}

// Add is idempotent: becoming a fan twice keeps a single row.
func (u *UserFanRepository) Add(fan *model.UserFan) error {
	result := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fan_id"}, {Name: "idol_id"}},
		DoNothing: true,
	}).Create(fan)

	return translate(result.Error)
}

// Remove hard-deletes the row, so becoming a fan again doesn't collide with a
// soft-deleted row on the unique index.
func (u *UserFanRepository) Remove(fanId uint, idolId uint) error {
	result := u.db.
		Unscoped().
		Where("fan_id = ? and idol_id = ?", fanId, idolId).
		Delete(&model.UserFan{})

	return result.Error
}

func (u *UserFanRepository) CountOfIdol(idolId uint) (int64, error) {
	var count int64
	result := u.db.
		Model(&model.UserFan{}).
		Where("idol_id = ?", idolId).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (u *UserFanRepository) IsFanOf(fanId uint, idolId uint) (bool, error) {
	var count int64
	result := u.db.
		Model(&model.UserFan{}).
		Where("fan_id = ? and idol_id = ?", fanId, idolId).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// FindFansOfIdol returns one page of the users who are fans of idolId, the
// newest fans first, with profiles, plus how many fans they have.
func (u *UserFanRepository) FindFansOfIdol(idolId uint, page int, pageSize int) ([]model.User, int64, error) {
	count, err := u.CountOfIdol(idolId)
	if err != nil {
		return nil, 0, err
	}

	fans := []model.User{}
	result := u.db.
		Joins("join user_fans on user_fans.fan_id = users.id and user_fans.deleted_at is null").
		Where("user_fans.idol_id = ?", idolId).
		Order("user_fans.created_at desc, users.id").
		Scopes(newPaginate(page, pageSize).paginatedResult).
		Preload("Profile").
		Find(&fans)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return fans, count, nil
}

// FindIdolIdsOfFan returns which of idolIds fanId is a fan of, in one query.
func (u *UserFanRepository) FindIdolIdsOfFan(fanId uint, idolIds []uint) ([]uint, error) {
	ids := []uint{}
	if len(idolIds) == 0 {
		return ids, nil
	}
	result := u.db.
		Model(&model.UserFan{}).
		Where("fan_id = ? and idol_id in ?", fanId, idolIds).
		Pluck("idol_id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}
