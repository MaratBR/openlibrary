package app

import "github.com/MaratBR/openlibrary/internal/store"

type AgeRating string

const (
	AgeRatingUnknown AgeRating = "?"
	AgeRatingG       AgeRating = "G"
	AgeRatingPG      AgeRating = "PG"
	AgeRatingPG13    AgeRating = "PG-13"
	AgeRatingR       AgeRating = "R"
	AgeRatingNC17    AgeRating = "NC-17"
)

var AllRatings []AgeRating = []AgeRating{
	AgeRatingG,
	AgeRatingPG,
	AgeRatingPG13,
	AgeRatingR,
	AgeRatingNC17,
	AgeRatingUnknown,
}

func ageRatingDbValue(r AgeRating) store.AgeRating {
	switch r {
	case AgeRatingG:
		return store.AgeRatingG
	case AgeRatingPG:
		return store.AgeRatingPG
	case AgeRatingPG13:
		return store.AgeRatingPG13
	case AgeRatingR:
		return store.AgeRatingR
	case AgeRatingNC17:
		return store.AgeRatingNC17
	default:
		return store.AgeRatingValue0
	}
}

func ageRatingFromDbValue(v store.AgeRating) AgeRating {
	switch v {
	case store.AgeRatingG:
		return AgeRatingG
	case store.AgeRatingPG:
		return AgeRatingPG
	case store.AgeRatingPG13:
		return AgeRatingPG13
	case store.AgeRatingR:
		return AgeRatingR
	case store.AgeRatingNC17:
		return AgeRatingNC17
	default:
		return AgeRatingUnknown
	}
}

func AsRating(v string) AgeRating {
	switch v {
	case "G":
		return AgeRatingG
	case "PG":
		return AgeRatingPG
	case "PG-13":
		return AgeRatingPG13
	case "R":
		return AgeRatingR
	case "NC-17":
		return AgeRatingNC17
	default:
		return AgeRatingUnknown
	}
}

func (r AgeRating) IsAdult() bool {
	return r == AgeRatingNC17 || r == AgeRatingR
}
