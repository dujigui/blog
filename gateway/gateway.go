package gateway

import (
	"github.com/kataras/iris/v12/context"
)

const (
	ADMIN = "/admin"
)

// 401 unauthorized
// 403 forbidden
// https://docs.iris-go.com/iris/request-authentication
func Gateway(ctx context.Context) {
	if ctx.Path() == ADMIN {
		if err := auth(ctx); err != nil {
			ctx.StatusCode(401)
			ctx.WriteString("Unauthorized")
			return
		}
	}

	ctx.Next()
}

func auth(ctx context.Context) error {
	return nil
}
