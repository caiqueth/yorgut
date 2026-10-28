package repository

import (
	"yorgut/internal/model"

	"gorm.io/gorm"
)

// Hard deletes for account deletion (UserRepository.DeleteForever,
// CommunityRepository.DeleteForever). Every table references its parents
// without ON DELETE CASCADE, so rows go children first; soft-deleted rows
// still hold their foreign keys, so every delete is Unscoped.

// deleteForever deletes the rows of value's table matching query and args.
func deleteForever(db *gorm.DB, value any, query string, args ...any) error {
	return db.Unscoped().Where(query, args...).Delete(value).Error
}

// deletePollsForever deletes the polls matching query and args, with their
// options and votes.
func deletePollsForever(db *gorm.DB, query string, args ...any) error {
	pollIds := db.Unscoped().Model(&model.Poll{}).Select("id").Where(query, args...)
	var ids []uint
	if err := pollIds.Pluck("id", &ids).Error; err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if err := deleteForever(db, &model.PollVote{}, "poll_id in ?", ids); err != nil {
		return err
	}
	if err := deleteForever(db, &model.PollOption{}, "poll_id in ?", ids); err != nil {
		return err
	}
	return deleteForever(db, &model.Poll{}, "id in ?", ids)
}

// deleteTopicsForever deletes the given topics and every reply under them,
// at any depth.
func deleteTopicsForever(db *gorm.DB, topicIds []uint) error {
	all := append([]uint{}, topicIds...)
	for parents := topicIds; len(parents) > 0; {
		var replies []uint
		if err := db.Unscoped().Model(&model.Topic{}).Where("parent_topic_id in ?", parents).Pluck("id", &replies).Error; err != nil {
			return err
		}
		all = append(all, replies...)
		parents = replies
	}
	if len(all) == 0 {
		return nil
	}
	// Unlink them first, so the order of the delete doesn't matter.
	if err := db.Unscoped().Model(&model.Topic{}).Where("id in ?", all).Update("parent_topic_id", nil).Error; err != nil {
		return err
	}
	return deleteForever(db, &model.Topic{}, "id in ?", all)
}
