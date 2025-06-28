package handler

import (
	"net/http"

	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/logic"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRedirectLogic(r.Context(), svcCtx)
		resp, err := l.Redirect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if resp.OringinalUrl != "" {
			http.Redirect(w, r, resp.OringinalUrl, http.StatusFound)
			return
		}
		// 理论上不会走到这里，如果走到这里，返回错误信息
		httpx.ErrorCtx(r.Context(), w, err)
	}
}
