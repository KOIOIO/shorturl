package logic

import (
	"context"
	"errors"
	"example.com/shorturl/short-url/zero_remake/common/errmsg"
	"example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"
	"strings"

	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/svc"
	"example.com/shorturl/short-url/zero_remake/shorturl_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedirectLogic) Redirect(req *types.RedirectRequest) (resp *types.RedirectResponse, err error) {
	// 1. 处理短链接参数
	shorturl := strings.TrimSpace(req.ShortURL)

	// 2. 检查短链接是否为空
	if shorturl == "" {
		return nil, errors.New("提供的短链为空")
	}

	// 3. 调用RPC服务处理短链接
	rsp, err := l.svcCtx.ShortUrlRpc.HandleShort(l.ctx, &shortUrl.HandleShortRequest{
		Shortcode: shorturl,
	})

	// 4. 处理RPC调用错误
	if err != nil {
		return nil, errors.New("短链服务暂时不可用")
	}

	// 5. 检查短链接处理结果
	if rsp.Code != errmsg.SUCCESS {
		return nil, errors.New("你要访问的短链找不到了")
	}

	// 6. 返回重定向响应
	return &types.RedirectResponse{
		Code:         errmsg.SUCCESS,
		OringinalUrl: rsp.LongUrl,
		Message:      "重定向成功",
	}, nil
}
