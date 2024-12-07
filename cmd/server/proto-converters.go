package server

import (
	"github.com/MaratBR/openlibrary/cmd/server/olproto"
	"github.com/MaratBR/openlibrary/internal/app"
)

func ageRatingToProto(v app.AgeRating) olproto.ProtoAgeRating {
	switch v {
	case app.AgeRatingG:
		return olproto.ProtoAgeRating_G
	case app.AgeRatingPG:
		return olproto.ProtoAgeRating_PG
	case app.AgeRatingPG13:
		return olproto.ProtoAgeRating_PG13
	case app.AgeRatingR:
		return olproto.ProtoAgeRating_R
	case app.AgeRatingNC17:
		return olproto.ProtoAgeRating_NC17
	default:
		return olproto.ProtoAgeRating_UNKNOWN
	}
}

func tagCategoryToProto(v app.TagsCategory) olproto.ProtoTagsCategory {
	switch v {
	case app.TagsCategoryOther:
		return olproto.ProtoTagsCategory_OTHER
	case app.TagsCategoryWarning:
		return olproto.ProtoTagsCategory_WARNING
	case app.TagsCategoryFandom:
		return olproto.ProtoTagsCategory_FANDOM
	case app.TagsCategoryRelationship:
		return olproto.ProtoTagsCategory_REL
	case app.TagsCategoryRelationshipType:
		return olproto.ProtoTagsCategory_REL_TYPE
	default:
		return olproto.ProtoTagsCategory_OTHER
	}
}
