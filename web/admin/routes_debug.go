package admin

import "github.com/MaratBR/openlibrary/internal/app"

type debugActionDescriptor struct {
	Func func() string
}

type debugController struct {
	fullReindexService app.BookReindexService
	actions            map[string]debugActionDescriptor
}

func newDebugController(fullReindexService app.BookReindexService) *debugController {
	c := &debugController{
		fullReindexService: fullReindexService,
		actions: map[string]debugActionDescriptor{
			"books:elastic:reindex": {
				Func: func() string {
					err := fullReindexService.ScheduleReindexAll()
					if err == nil {
						return ""
					}

					return err.Error()
				},
			},
		},
	}
	return c
}
