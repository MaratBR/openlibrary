package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hjson/hjson-go/v4"
)

type TagsCategory uint8

const (
	TagsCategoryOther TagsCategory = iota
	TagsCategoryWarning
	TagsCategoryFandom
	TagsCategoryRelationship
	TagsCategoryRelationshipType
)

func tagsCategoryName(cat TagsCategory) string {
	switch cat {
	case TagsCategoryOther:
		return "other"
	case TagsCategoryWarning:
		return "warning"
	case TagsCategoryFandom:
		return "fandom"
	case TagsCategoryRelationship:
		return "relationship"
	case TagsCategoryRelationshipType:
		return "relationshipType"
	default:
		return "unknown"
	}
}

func tagsCategoryFromName(name string) TagsCategory {
	switch name {
	case "other":
		return TagsCategoryOther
	case "warning":
		return TagsCategoryWarning
	case "fandom":
		return TagsCategoryFandom
	case "relationship":
		return TagsCategoryRelationship
	case "relationshipType":
		return TagsCategoryRelationshipType
	default:
		return TagsCategoryOther
	}
}

func (t TagsCategory) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, tagsCategoryName(t))), nil
}

func (t *TagsCategory) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = tagsCategoryFromName(s)
	return nil
}

type TagDto struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	IsAdult   bool         `json:"isAdult"`
	IsSpoiler bool         `json:"isSpoiler"`
	Category  TagsCategory `json:"category"`
	// If set to true, the tag is officially defined and can be used in search.
	IsDefined bool `json:"isDefined"`
}

type predefinedTag struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	IsAdult   bool         `json:"isAdult"`
	IsSpoiler bool         `json:"isSpoiler"`
	Category  TagsCategory `json:"category"`
}

//go:embed predefined-tags
var predefinedTagsFS embed.FS

func LoadPredefinedTags() (map[string]TagDto, error) {
	var tags []TagDto

	entries, err := predefinedTagsFS.ReadDir("predefined-tags")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		var pt []predefinedTag
		file, err := predefinedTagsFS.Open(fmt.Sprintf("predefined-tags/%s", entry.Name()))
		if err != nil {
			return nil, err
		}
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		err = hjson.Unmarshal(fileContent, &pt)
		if err != nil {
			return nil, err
		}

		for _, t := range pt {
			id := t.ID
			if id == "" {
				id = t.Name
			}
			tags = append(tags, TagDto{
				ID:        id,
				Name:      t.Name,
				IsAdult:   t.IsAdult,
				IsSpoiler: t.IsSpoiler,
				Category:  t.Category,
				IsDefined: true,
			})
		}
	}

	tagsMap := map[string]TagDto{}

	for _, tag := range tags {
		tagsMap[tag.Name] = tag
	}

	return tagsMap, err
}

type TagsService struct {
	predefinedTags map[string]TagDto
}

func NewTagsService() *TagsService {
	predefinedTags, err := LoadPredefinedTags()
	if err != nil {
		panic(err)
	}
	return &TagsService{predefinedTags: predefinedTags}
}

func (t *TagsService) FindTags(names []string) []TagDto {
	tags := make([]TagDto, 0)

	for _, name := range names {
		tag, ok := t.predefinedTags[name]
		if ok {
			tags = append(tags, tag)
		}
	}

	return tags
}
