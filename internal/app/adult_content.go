package app

// AdultContentService owns audience classification so callers do not encode
// rating boundaries or tag rules. Replace its Fx provider to customize policy.
type AdultContentService interface {
	IsAdult(rating AgeRating, tags []DefinedTagDto) bool
}

type adultContentService struct{}

func NewAdultContentService() AdultContentService {
	return adultContentService{}
}

func (adultContentService) IsAdult(rating AgeRating, tags []DefinedTagDto) bool {
	if rating == AgeRatingNC17 {
		return true
	}
	for _, tag := range tags {
		if tag.IsAdult {
			return true
		}
	}
	return false
}
