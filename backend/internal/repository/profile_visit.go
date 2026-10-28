package repository

import (
	"time"
	"yorgut/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProfileVisitRepository struct {
	repository[*model.ProfileVisit]
}

func NewProfileVisitRepository(Db *gorm.DB) *ProfileVisitRepository {
	return &ProfileVisitRepository{
		repository: repository[*model.ProfileVisit]{db: Db},
	}
}

// Record stores a visit at now: a new row for the first visit of the day,
// otherwise it moves that day's updated_at to now. One upsert, so concurrent
// visits can't create two rows.
func (r *ProfileVisitRepository) Record(visitorId uint, visitedId uint, day time.Time, now time.Time) error {
	visit := &model.ProfileVisit{VisitorID: visitorId, VisitedID: visitedId, Day: day}
	visit.CreatedAt = now
	visit.UpdatedAt = now

	result := r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "visited_id"}, {Name: "visitor_id"}, {Name: "day"}},
			DoUpdates: clause.Assignments(map[string]any{"updated_at": now}),
		}).
		Create(visit)

	return result.Error
}

// CountOfVisited counts the visits to the user's profile (one per visitor and
// day) since the given day, or ever when since is nil.
func (r *ProfileVisitRepository) CountOfVisited(visitedId uint, since *time.Time) (int64, error) {
	query := r.db.Model(&model.ProfileVisit{}).Where("visited_id = ?", visitedId)
	if since != nil {
		query = query.Where("day >= ?", *since)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindRecentVisitors returns the last limit users who visited the profile,
// most recent first. Users who turned off "registrar minhas visitas" are left
// out, their earlier visits included.
func (r *ProfileVisitRepository) FindRecentVisitors(visitedId uint, limit int) ([]model.User, error) {
	lastVisits := r.db.
		Model(&model.ProfileVisit{}).
		Select("visitor_id, max(updated_at) as last_visit").
		Where("visited_id = ?", visitedId).
		Group("visitor_id")

	visitors := []model.User{}
	result := r.db.
		Model(&model.User{}).
		Joins("join (?) as visits on visits.visitor_id = users.id", lastVisits).
		Where("users.record_visits = ?", true).
		Order("visits.last_visit desc").
		Limit(limit).
		Find(&visitors)
	if result.Error != nil {
		return nil, result.Error
	}
	return visitors, nil
}
