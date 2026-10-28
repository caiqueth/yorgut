package service

import (
	"time"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/rs/zerolog/log"
)

// recentVisitorsShown is how many visitors the home page lists.
const recentVisitorsShown = 10

// ProfileVisitService records who opens whose profile, for the home page's
// profile views and recent visitors.
type ProfileVisitService struct {
	profileVisitRepository repository.ProfileVisitRepository
	userRepository         repository.UserRepository
}

func NewProfileVisitService(
	profileVisitRepository repository.ProfileVisitRepository,
	userRepository repository.UserRepository,
) *ProfileVisitService {
	return &ProfileVisitService{
		profileVisitRepository: profileVisitRepository,
		userRepository:         userRepository,
	}
}

// dayOf returns local midnight of t's day.
func dayOf(t time.Time) time.Time {
	t = t.In(time.Local)
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// RecordVisit records that visitorId opened visitedId's profile. Visiting
// your own profile doesn't count, and nothing is recorded for visitors who
// turned off "registrar minhas visitas". A failure is only logged: it must
// not break showing the profile.
func (s *ProfileVisitService) RecordVisit(visitorId uint, visitedId uint) {
	if visitorId == visitedId {
		return
	}
	recordVisits, err := s.userRepository.RecordsVisits(visitorId)
	if err != nil {
		log.Warn().Err(err).Uint("visitorId", visitorId).Msg("could not read the visit setting; not recording the visit")
		return
	}
	if !recordVisits {
		return
	}
	now := time.Now()
	if err := s.profileVisitRepository.Record(visitorId, visitedId, dayOf(now), now); err != nil {
		log.Warn().Err(err).Uint("visitorId", visitorId).Uint("visitedId", visitedId).Msg("could not record profile visit")
	}
}

// FindProfileViews counts the visits to the user's profile (one per visitor
// and day), ever and in the last 7 days (today included).
func (s *ProfileVisitService) FindProfileViews(userId uint, now time.Time) (*transport.ProfileViews, error) {
	allTime, err := s.profileVisitRepository.CountOfVisited(userId, nil)
	if err != nil {
		return nil, err
	}

	weekAgo := dayOf(now).AddDate(0, 0, -6)
	lastSevenDays, err := s.profileVisitRepository.CountOfVisited(userId, &weekAgo)
	if err != nil {
		return nil, err
	}

	return &transport.ProfileViews{AllTime: uint(allTime), LastSevenDays: uint(lastSevenDays)}, nil
}

// FindRecentVisitors returns the last users who visited the profile, most
// recent first. Users who turned off "registrar minhas visitas" don't see
// their visitors (nobody sees theirs), so they get an empty list.
func (s *ProfileVisitService) FindRecentVisitors(userId uint) ([]transport.RecentVisitor, error) {
	recordVisits, err := s.userRepository.RecordsVisits(userId)
	if err != nil {
		return nil, err
	}
	if !recordVisits {
		return []transport.RecentVisitor{}, nil
	}

	visitors, err := s.profileVisitRepository.FindRecentVisitors(userId, recentVisitorsShown)
	if err != nil {
		return nil, err
	}
	return transport.RecentVisitorsFromUsers(visitors), nil
}
