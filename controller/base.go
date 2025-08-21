package controller

import (
	"goumang-master/errorcode"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func returnErrJson(ctx *gin.Context, code int, extraMsg ...string) {
	extraMsgStr := strings.Join(extraMsg, ";")
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg: func() string {
			if len(extraMsg) > 0 {
				return extraMsgStr
			}
			return errorcode.CodeMsg(code)
		}(),
	})
}

func returnSuccessJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: errorcode.Success,
		Msg:  errorcode.CodeMsg(errorcode.Success),
		Data: data,
	})
}

func paramsValidator(ctx *gin.Context, paramsPtr interface{}) (err error) {
	if err = ctx.ShouldBind(paramsPtr); err != nil {
		returnErrJson(ctx, errorcode.ErrParams, err.Error())
	}
	return
}
