package app

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"io"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/hjson/hjson-go/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type tagsPopulateService struct {
	queries *store.Queries
	log     *zap.SugaredLogger
}

func startTagsPopulateService(db store.DBTX, lc fx.Lifecycle, cfg *koanf.Koanf, log *zap.SugaredLogger) {
	if cfg.Bool("init.import-predefined-tags") {
		dir := cfg.String("init.import-predefined-tags-path")

		lc.Append(fx.StartHook(func() {
			queries := store.New(db)
			svc := &tagsPopulateService{queries: queries, log: log}
			go svc.importTags(context.Background(), dir)
		}))
	} else {
		log.Debug("init.import-predefined-tags is false, tags import is skipped")
	}

}

func (t *tagsPopulateService) importTags(ctx context.Context, dir string) {
	tags, err := t.loadTags(dir)
	if err != nil {
		slog.Error("failed to load tags", "err", err, "dir", dir)
		return
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
		err = t.queries.RemoveUnusedDefaultTags(ctx, tagNames)
		if err != nil {
			t.log.Errorw("RemoveUnusedDefaultTags failed", "err", err)
		} else {
			t.log.Debugw("called RemoveUnusedDefaultTags")
		}
	}

	t.log.Debugf("importing %d tags", len(tagRows))
	// t.log.Debugw("tags for import", "tags", tagRows)
	err = importTagsBulk(ctx, t.queries, tagRows)
	if err != nil {
		t.log.Errorw("importTagsBulk failed during tags population", "err", err)
	} else {
		t.log.Debug("import of tags finished")
	}
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

func importTagsBulk(ctx context.Context, queries *store.Queries, tags []tagImportRow) error {
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

func (t *tagsPopulateService) loadTags(dir string) ([]predefinedTag, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.log.Errorw("failed to import default tags", "err", err)
		return nil, err
	}

	tags := make([]predefinedTag, 0)

	for _, entry := range entries {
		file, err := os.Open(path.Join(dir, entry.Name()))
		if err != nil {
			t.log.Errorw("failed to open a file with tags", "err", err, "file", entry.Name())
			continue
		}
		fileContent, err := io.ReadAll(file)
		if err != nil {
			t.log.Errorw("failed to read a file with tags", "err", err, "file", entry.Name())
			continue
		}
		var pt []predefinedTag
		err = hjson.Unmarshal(fileContent, &pt)
		if err != nil {
			t.log.Errorw("failed to parse a file with tags", "err", err, "file", entry.Name())
			continue
		}

		tags = append(tags, pt...)
	}

	return tags, nil
}

type predefinedTag struct {
	Name        string       `json:"name"`
	IsAdult     bool         `json:"adult"`
	IsSpoiler   bool         `json:"spoiler"`
	Category    TagsCategory `json:"category"`
	Description string       `json:"description"`
}
