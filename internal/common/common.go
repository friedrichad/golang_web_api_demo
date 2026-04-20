package common

import "fmt"

type Error struct {
	Code    string
	Message string
}

var Success = &Error{Code: "200", Message: "Thành công"}
var AuthenticationFail = &Error{Code: "400", Message: "Username hoặc mật khẩu không đúng"}
var SystemError = &Error{Code: "999", Message: "Có lỗi xảy ra, vui lòng quay lại sau"}
var TokenInvalid = &Error{Code: "401", Message: "Token không hợp lệ hoặc đã hết hạn"}
var FileEmpty = &Error{Code: "400", Message: "File không được để trống"}
var CreateBucketFailed = &Error{Code: "500", Message: "Không tạo được bucket"}
var UploadFileFailed = &Error{Code: "500", Message: "Upload file không thành công"}
var RequestInvalid = &Error{Code: "400", Message: "Request không hợp lệ"}
var WrongPassword = &Error{Code: "500", Message: "Mật khẩu không đúng"}

func ParamRequired(name string) *Error {
	return &Error{Code: "400", Message: fmt.Sprintf("%s không được để trống", name)}
}

func ParamInvalid(name string) *Error {
	return &Error{Code: "400", Message: fmt.Sprintf("%s không hợp lệ", name)}
}

func ObjectNotExisted(name string) *Error {
	return &Error{Code: "404", Message: fmt.Sprintf("%s không tồn tại", name)}
}
