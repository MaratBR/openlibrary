package app

import (
	"github.com/MaratBR/openlibrary/internal/app/idgen"
	"github.com/MaratBR/openlibrary/internal/commonutil"
)

func GenID() int64 { return idgen.GenID() }

func genOpaqueID() string {
	s, err := commonutil.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	return s
}
