package app

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log/slog"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/hjson/hjson-go/v4"
	"github.com/jackc/pgx/v5/pgtype"
)

type predefinedTag struct {
	Name        string       `json:"name"`
	IsAdult     bool         `json:"adult"`
	IsSpoiler   bool         `json:"spoiler"`
	Category    TagsCategory `json:"category"`
	Description string       `json:"description"`
}

//go:embed predefined-tags
var predefinedTagsFS embed.FS

func loadPredefinedTags() ([]predefinedTag, error) {
	tags := make([]predefinedTag, 0)

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

		tags = append(tags, pt...)
	}

	return tags, err
}

type tagImportRow struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	IsAdult     bool          `json:"is_adult"`
	IsSpoiler   bool          `json:"is_spoiler"`
	TagType     store.TagType `json:"tag_type"`
	Description string        `json:"description"`
	SynonymOf   pgtype.Int8   `json:"synonym_of"`
	CreatedAt   time.Time     `json:"created_at"`
}

func ImportPredefinedTags(ctx context.Context, queries *store.Queries) error {
	tags, err := loadPredefinedTags()
	if err != nil {
		panic(err)
	}

	tagRows := make([]tagImportRow, len(tags))
	seenNames := map[string]struct{}{}
	fnvHash := fnv.New64()

	for i, tag := range tags {
		if _, ok := seenNames[tag.Name]; ok {
			continue
		}
		seenNames[tag.Name] = struct{}{}
		fnvHash.Write([]byte(tag.Name))
		uint64Hash := fnvHash.Sum64() & ^(uint64(1) << 63)
		id := int64(uint64Hash)
		id = id - id%10 + 0

		tagRows[i] = tagImportRow{
			ID:          id,
			Name:        tag.Name,
			IsAdult:     tag.IsAdult,
			IsSpoiler:   tag.IsSpoiler,
			TagType:     tagsCategoryToDbTagType(tag.Category),
			Description: tag.Description,
			SynonymOf:   pgtype.Int8{Valid: false},
			CreatedAt:   time.Now(),
		}
	}

	{
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		err = queries.RemoveUnusedDefaultTags(ctx, tagNames)
		if err != nil {
			return err
		}
		slog.Debug("called RemoveUnusedDefaultTags")
	}

	return importTags(ctx, queries, tagRows)
}

func importTags(ctx context.Context, queries *store.Queries, tags []tagImportRow) error {
	jsonStr, err := json.Marshal(tags)
	if err != nil {
		panic(err)
	}

	err = queries.ImportTags(ctx, jsonStr)
	if err != nil {
		return err
	}

	return nil

}
