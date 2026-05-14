package common

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var Success = &Error{Code: "200", Message: "Thành công"}
var AuthenticationFail = &Error{Code: "400", Message: "Username hoặc mật khẩu không đúng"}
var SystemError = &Error{Code: "999", Message: "Có lỗi xảy ra, vui lòng quay lại sau"}
var TokenInvalid = &Error{Code: "401", Message: "Token không hợp lệ hoặc đã hết hạn"}
var FileEmpty = &Error{Code: "400", Message: "File không được để trống"}
var FileError = &Error{Code: "400", Message: "Có sự cố khi đăng tải file"}
var CreateBucketFailed = &Error{Code: "500", Message: "Không tạo được bucket"}
var UploadFileFailed = &Error{Code: "500", Message: "Upload file không thành công"}
var RequestInvalid = &Error{Code: "400", Message: "Request không hợp lệ"}
var WrongPassword = &Error{Code: "500", Message: "Mật khẩu không đúng"}
var NotFound = &Error{Code: "404", Message: "Không tìm thấy dữ liệu"}
var AlreadyLoggedIn = &Error{Code: "400", Message: "Trình duyệt đã có người đăng nhập, vui lòng đăng xuất trước khi đăng nhập bằng tài khoản khác"}
var AccountLocked = &Error{Code: "403", Message: "Tài khoản đang bị khóa do đăng nhập sai quá nhiều lần, vui lòng thử lại sau"}
var UserForbidden = &Error{Code: "403", Message: "Bạn không có quyền truy phê duyệt đơn với người có cấp cao hơn mình"}
func ParamRequired(name string) *Error {
	return &Error{Code: "400", Message: fmt.Sprintf("%s không được để trống", name)}
}

func ParamInvalid(name string) *Error {
	return &Error{Code: "400", Message: fmt.Sprintf("%s không hợp lệ", name)}
}

func ObjectNotExisted(name string) *Error {
	return &Error{Code: "404", Message: fmt.Sprintf("%s không tồn tại", name)}
}
