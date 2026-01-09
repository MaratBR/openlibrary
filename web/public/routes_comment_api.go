package public

import (
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
	r.Route("/comments", func(r chi.Router) {
		r.Use(requiresAuthorizationMiddleware)
		r.Post("/like", c.like)
		r.Get("/replies", c.replies)
	})
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

	result, err := c.commentsService.GetReplies(r.Context(), app.GetCommentRepliesQuery{
		ActorUserID: auth.GetNullableUserID(r.Context()),
		Limit:       20,
		Cursor:      uint32(cursor),
		CommentID:   commentId,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(map[string]any{
		"cursor":     result.Cursor,
		"nextCursor": result.NextCursor,
		"comments":   result.Comments,
	}).Write(w)
}
