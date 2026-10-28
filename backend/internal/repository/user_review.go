package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserReviewRepository struct {
	repository[*model.UserReview]
}

type UserReviewAverages struct {
	Cool        float64
	Sexy        float64
	Trustworthy float64
	ReviewCount int64
}

func NewUserReviewRepository(Db *gorm.DB) *UserReviewRepository {
	return &UserReviewRepository{
		repository: repository[*model.UserReview]{db: Db},
	}
}

func (u *UserReviewRepository) Upsert(review *model.UserReview) error {
	result := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "reviewer_id"}, {Name: "reviewed_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"cool", "sexy", "trustworthy", "updated_at", "deleted_at"}),
	}).Create(review)

	return translate(result.Error)
}

func (u *UserReviewRepository) FindOfReviewerAndReviewed(reviewerId uint, reviewedId uint) (*model.UserReview, error) {
	var ents []model.UserReview
	result := u.db.
		Where("reviewer_id = ? and reviewed_id = ?", reviewerId, reviewedId).
		Limit(1).
		Find(&ents)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(ents) == 0 {
		return nil, nil
	}

	return &ents[0], nil
}

func (u *UserReviewRepository) AveragesOfUser(reviewedId uint) (*UserReviewAverages, error) {
	var averages UserReviewAverages
	result := u.db.
		Model(&model.UserReview{}).
		Select(
			"COALESCE(AVG(cool), 0) as cool, "+
				"COALESCE(AVG(sexy), 0) as sexy, "+
				"COALESCE(AVG(trustworthy), 0) as trustworthy, "+
				"COUNT(*) as review_count",
		).
		Where("reviewed_id = ?", reviewedId).
		Scan(&averages)

	if result.Error != nil {
		return nil, result.Error
	}

	return &averages, nil
}

// FindOfReviewerForUsers returns the reviews reviewerId wrote about any of
// reviewedIds, in one query.
func (u *UserReviewRepository) FindOfReviewerForUsers(reviewerId uint, reviewedIds []uint) ([]model.UserReview, error) {
	reviews := []model.UserReview{}
	if len(reviewedIds) == 0 {
		return reviews, nil
	}
	result := u.db.
		Where("reviewer_id = ? and reviewed_id in ?", reviewerId, reviewedIds).
		Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}
