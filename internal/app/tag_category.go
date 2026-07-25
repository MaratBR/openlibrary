package app

import (
	"encoding/json"
	"fmt"

	"github.com/MaratBR/openlibrary/internal/store"
)

// go2tsdef:generate
// go2tsdef:override_type 'other' | 'warning' | 'fandom' | 'rel' | 'reltype' | 'unknown'
type TagsCategory uint8

const (
	TagsCategoryOther TagsCategory = iota
	TagsCategoryWarning
	TagsCategoryFandom
	TagsCategoryRelationship
	TagsCategoryRelationshipType
	TagsCategoryGenre
)

var TagsCategoryList = []TagsCategory{TagsCategoryOther, TagsCategoryWarning, TagsCategoryFandom, TagsCategoryRelationship, TagsCategoryRelationshipType, TagsCategoryGenre}

func (t TagsCategory) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%s"`, t.String()), nil
}

func (t *TagsCategory) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = TagsCategoryFromName(s)
	return nil
}

func (t TagsCategory) String() string {
	switch t {
	case TagsCategoryOther:
		return "other"
	case TagsCategoryWarning:
		return "warning"
	case TagsCategoryFandom:
		return "fandom"
	case TagsCategoryRelationship:
		return "rel"
	case TagsCategoryRelationshipType:
		return "reltype"
	case TagsCategoryGenre:
		return "genre"
	default:
		return "unknown"
	}
}

func dbTagTypeToTagsCategory(t store.TagType) TagsCategory {
	switch t {
	case store.TagTypeFreeform:
		return TagsCategoryOther
	case store.TagTypeWarning:
		return TagsCategoryWarning
	case store.TagTypeFandom:
		return TagsCategoryFandom
	case store.TagTypeRel:
		return TagsCategoryRelationship
	case store.TagTypeReltype:
		return TagsCategoryRelationshipType
	case store.TagTypeGenre:
		return TagsCategoryGenre
	default:
		return TagsCategoryOther
	}
}

func tagsCategoryToDbTagType(cat TagsCategory) store.TagType {
	switch cat {
	case TagsCategoryOther:
		return store.TagTypeFreeform
	case TagsCategoryWarning:
		return store.TagTypeWarning
	case TagsCategoryFandom:
		return store.TagTypeFandom
	case TagsCategoryRelationship:
		return store.TagTypeRel
	case TagsCategoryRelationshipType:
		return store.TagTypeReltype
	case TagsCategoryGenre:
		return store.TagTypeGenre
	default:
		return store.TagTypeFreeform
	}
}

func TagsCategoryFromName(name string) TagsCategory {
	switch name {
	case "other":
		return TagsCategoryOther
	case "warning":
		return TagsCategoryWarning
	case "fandom":
		return TagsCategoryFandom
	case "rel":
		return TagsCategoryRelationship
	case "reltype":
		return TagsCategoryRelationshipType
	case "genre":
		return TagsCategoryGenre
	default:
		return TagsCategoryOther
	}
}
