package globalMiddleware

import (
	"fmt"
	"market_server/common/globalKey"
	"net/http"
	"strconv"
)

type SetUidToCtxMiddleware struct {
}

func NewSetUidToCtxMiddleware() *SetUidToCtxMiddleware {
	return &SetUidToCtxMiddleware{}
}

func (m *SetUidToCtxMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("header:%+v", r.Header)
		userId, _ := strconv.ParseInt(r.Header.Get(globalKey.XMEMBERID), 10, 64)
		ctx := r.Context()
		fmt.Println("userId:", userId)
		//ctx = context.WithValue(ctx, ctxdata.CtxKeyJwtUserId, userId)

		next(w, r.WithContext(ctx))
	}
}
