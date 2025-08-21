package errorcode

const (
	Success = 0

	ErrServiceException = 400

	ErrParams = 10001
)

func CodeMsg(code int) (msg string) {
	switch code {
	case Success:
		msg = "成功"
	case ErrServiceException:
		msg = "服务异常"
	case ErrParams:
		msg = "请求参数错误"
	default:
	}
	return
}
