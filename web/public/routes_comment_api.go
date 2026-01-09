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
