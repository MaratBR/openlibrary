package public

import (
	"encoding/json"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
)

type apiControllerComments struct {
	commentsService app.CommentsService
}

func newAPICommentsController(
	commentsService app.CommentsService,
) *apiControllerComments {
	return &apiControllerComments{
		commentsService: commentsService,
	}
}

func (c *apiControllerComments) Register(r chi.Router) {
	r.Get("/comments", c.list)
	r.Get("/comments/replies", c.replies)
	r.Group(func(r chi.Router) {
		r.Use(requiresAuthorizationMiddleware)
		r.Post("/comments/like", c.like)
		r.Post("/comments/add", c.add)
	})
}

func (c *apiControllerComments) list(w http.ResponseWriter, r *http.Request) {
	chapterID, err := olhttp.URLQueryParamInt64(r, "chapterId")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	cursor, _ := olhttp.URLQueryParamInt64(r, "cursor")
	result, err := c.commentsService.GetList(r.Context(), app.GetCommentsQuery{
		ActorUserID: auth.GetNullableUserID(r.Context()),
		ChapterID:   chapterID,
		Limit:       30,
		Cursor:      uint32(cursor),
		Sort:        app.ParseCommentSort(r.URL.Query().Get("sort")),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(map[string]any{
		"cursor": result.Cursor, "nextCursor": result.NextCursor, "comments": result.Comments, "total": result.Total,
	}).Write(w)
}

func (c *apiControllerComments) add(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChapterID       int64  `json:"chapterId,string"`
		ParentCommentID *int64 `json:"parentCommentId,string"`
		Content         string `json:"content"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	command := app.AddCommentCommand{UserID: auth.RequireSession(r.Context()).UserID, ChapterID: input.ChapterID, Content: input.Content}
	if input.ParentCommentID != nil {
		command.ParentCommentID = app.Value(*input.ParentCommentID)
	}
	result, err := c.commentsService.AddComment(r.Context(), command)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(result.Comment).Write(w)
}

func (c *apiControllerComments) like(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	like := olhttp.GetBoolDefault(urlQuery, "like", false)
	commentID, err := olhttp.URLQueryParamInt64(r, "commentId")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	s := auth.RequireSession(r.Context())

	liked, err := c.commentsService.LikeComment(r.Context(), app.LikeCommentCommand{
		CommentID: commentID,
		Like:      like,
		UserID:    s.UserID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(liked).Write(w)
}

func (c *apiControllerComments) replies(w http.ResponseWriter, r *http.Request) {
	cursor, _ := olhttp.URLQueryParamInt64(r, "cursor")
	commentId, _ := olhttp.URLQueryParamInt64(r, "commentId")

	const pageSize = 2
	result, err := c.commentsService.GetReplies(r.Context(), app.GetCommentRepliesQuery{
		ActorUserID: auth.GetNullableUserID(r.Context()),
		Limit:       pageSize + 1,
		Cursor:      uint32(cursor),
		CommentID:   commentId,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	if len(result.Comments) > pageSize {
		result.Comments = result.Comments[:pageSize]
		result.NextCursor = uint32(result.Comments[pageSize-1].CreatedAt.Unix())
	} else {
		result.NextCursor = 0
	}
	olhttp.NewAPIResponse(map[string]any{
		"cursor":     result.Cursor,
		"nextCursor": result.NextCursor,
		"comments":   result.Comments,
	}).Write(w)
}
