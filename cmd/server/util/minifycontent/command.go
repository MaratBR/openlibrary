package minifycontent

import (
	"context"
	"fmt"

	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

const batchSize = 1000

func Run(ctx context.Context, db dal.DB, log *zap.SugaredLogger) error {
	queries := store.New(db)
	if err := processBookSummaries(ctx, db, queries, log); err != nil {
		return err
	}
	return processChapterContent(ctx, db, queries, log)
}

func processBookSummaries(ctx context.Context, db dal.DB, queries *store.Queries, log *zap.SugaredLogger) error {
	var afterID int64
	var total, updated int

	markup := content.NewDefaultEngine()

	log.Infow("started processing book summaries")
	for {
		books, err := queries.CLI_Util_MinifyContent_ListBookSummaries(ctx, store.CLI_Util_MinifyContent_ListBookSummariesParams{
			ID:    afterID,
			Limit: batchSize,
		})
		if err != nil {
			return fmt.Errorf("list books after ID %d: %w", afterID, err)
		}
		if len(books) == 0 {
			break
		}

		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin book batch: %w", err)
		}
		txQueries := queries.WithTx(tx)
		for _, book := range books {
			processed, err := markup.Clean(book.Summary)
			if err != nil {
				_ = tx.Rollback(ctx)
				return fmt.Errorf("process book %d summary: %w", book.ID, err)
			}
			if processed.Sanitized != book.Summary {
				if err := txQueries.CLI_Util_MinifyContent_UpdateBookSummary(ctx, store.CLI_Util_MinifyContent_UpdateBookSummaryParams{
					ID:      book.ID,
					Summary: processed.Sanitized,
				}); err != nil {
					_ = tx.Rollback(ctx)
					return fmt.Errorf("update book %d summary: %w", book.ID, err)
				}
				updated++
			}
			afterID = book.ID
			total++
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit book batch ending at ID %d: %w", afterID, err)
		}

		log.Infow("processed book summary batch", "processed", total, "updated", updated)
	}

	log.Infow("finished processing book summaries", "processed", total, "updated", updated)
	return nil
}

func processChapterContent(ctx context.Context, db dal.DB, queries *store.Queries, log *zap.SugaredLogger) error {
	var afterID int64
	var total, updated int

	markup := content.NewDefaultEngine()

	log.Infow("started processing chapter content")
	for {
		chapters, err := queries.CLI_Util_MinifyContent_ListChapters(ctx, store.CLI_Util_MinifyContent_ListChaptersParams{
			ID:    afterID,
			Limit: batchSize,
		})
		if err != nil {
			return fmt.Errorf("list chapters after ID %d: %w", afterID, err)
		}
		if len(chapters) == 0 {
			break
		}

		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin chapter batch: %w", err)
		}
		txQueries := queries.WithTx(tx)
		for _, chapter := range chapters {
			processed, err := markup.Clean(chapter.Content)
			if err != nil {
				_ = tx.Rollback(ctx)
				return fmt.Errorf("process chapter %d: %w", chapter.ID, err)
			}
			if processed.Sanitized != chapter.Content || processed.Words != chapter.Words {
				if err := txQueries.CLI_Util_MinifyContent_UpdateChapter(ctx, store.CLI_Util_MinifyContent_UpdateChapterParams{
					ID:      chapter.ID,
					Content: processed.Sanitized,
					Words:   processed.Words,
				}); err != nil {
					_ = tx.Rollback(ctx)
					return fmt.Errorf("update chapter %d: %w", chapter.ID, err)
				}
				updated++
			}
			afterID = chapter.ID
			total++
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit chapter batch ending at ID %d: %w", afterID, err)
		}

		log.Infow("processed chapter batch", "processed", total, "updated", updated)
	}

	log.Infow("finished processing chapters", "processed", total, "updated", updated)
	return nil
}
