package model

type BaseResponse struct {
	Code    int32  `json:"code"`    // 响应码
	Message string `json:"message"` // 响应消息
}

type Code int32

const (
	CodeOK                Code = 0  // 请求成功
	CodeError             Code = 1  // 请求失败
	CodeNotFound          Code = 2  // 请求不存在
	CodeInternal          Code = 3  // 内部错误
	CodeInvalidParam      Code = 4  // 参数错误
	CodeInvalidValue      Code = 5  // 值错误
	CodeInvalidLength     Code = 6  // 长度错误
	CodeInvalidSize       Code = 7  // 大小错误
	CodeInvalidRange      Code = 8  // 范围错误
	CodeInvalidRegex      Code = 9  // 正则错误
	CodeInvalidEmail      Code = 10 // 邮箱错误
	CodeInvalidPhone      Code = 11 // 手机号错误
	CodeInvalidCosKey     Code = 12 // cos key 错误
	CodeInvalidBucketId   Code = 13 // bucket id 错误
	CodeInvalidKey        Code = 14 // key 错误
	CodeInvalidBusinessId Code = 15 // business id 错误
	CodeDownloadFailed    Code = 16 // 下载失败
)

type CodeMessage string

const (
	CodeMessageOK                CodeMessage = "请求成功"
	CodeMessageError             CodeMessage = "请求失败"
	CodeMessageNotFound          CodeMessage = "请求不存在"
	CodeMessageInternal          CodeMessage = "内部错误"
	CodeMessageInvalidParam      CodeMessage = "参数错误"
	CodeMessageInvalidValue      CodeMessage = "值错误"
	CodeMessageInvalidLength     CodeMessage = "长度错误"
	CodeMessageInvalidSize       CodeMessage = "大小错误"
	CodeMessageInvalidRange      CodeMessage = "范围错误"
	CodeMessageInvalidRegex      CodeMessage = "正则错误"
	CodeMessageInvalidEmail      CodeMessage = "邮箱错误"
	CodeMessageInvalidPhone      CodeMessage = "手机号错误"
	CodeMessageInvalidCosKey     CodeMessage = "cos key 错误"
	CodeMessageInvalidBucketId   CodeMessage = "bucket id 错误"
	CodeMessageInvalidKey        CodeMessage = "key 错误"
	CodeMessageInvalidBusinessId CodeMessage = "business id 错误"
	CodeMessageDownloadFailed    CodeMessage = "下载失败"
)
