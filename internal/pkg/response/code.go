package response

import "racent.com/pkg/i18n"

type responseCode struct {
	Code int
	Msg  string
}

const (
	ErrJwtTokenExpiredMaxRefresh = -3
	ErrJwtTokenExpired           = -2
	ErrIsGuest                   = -1
	Success                      = 0
	ErrDefault                   = 1
	ErrValidate                  = 2
	ErrAttemptLogin              = 3
	ErrUserExist                 = 4
	ErrGuest                     = 5
	ErrJwtTokenMalformed         = 6
	ErrJwtTokenInvalid           = 7
	ErrJwtTokenEmpty             = 8
	ErrJwtHeaderMalformed        = 9
	ErrOriginPasswordInvalid     = 10
	ErrUnauthorized              = 401
	ErrNotFound                  = 404
	ErrUnprocessableEntity       = 422
	ErrDNSTxtRecordMissing       = 800
	ErrApiSignatureInvalid       = 999
)

func CodeText(code int) string {
	switch code {
	case ErrIsGuest:
		return i18n.Obj.Trans("err_is_guest")
	case Success:
		return i18n.Obj.Trans("success")
	case ErrDefault:
		return i18n.Obj.Trans("err_default")
	case ErrValidate:
		return i18n.Obj.Trans("err_validate")
	case ErrAttemptLogin:
		return i18n.Obj.Trans("err_attempt_login")
	case ErrUserExist:
		return i18n.Obj.Trans("err_user_exist")
	case ErrGuest:
		return i18n.Obj.Trans("err_guest")
	case ErrJwtTokenExpired:
		return i18n.Obj.Trans("err_jwt_token_expired")
	case ErrJwtTokenExpiredMaxRefresh:
		return i18n.Obj.Trans("err_jwt_token_expired_max_refresh")
	case ErrJwtTokenMalformed:
		return i18n.Obj.Trans("err_jwt_token_malformed")
	case ErrJwtTokenInvalid:
		return i18n.Obj.Trans("err_jwt_token_invalid")
	case ErrJwtTokenEmpty:
		return i18n.Obj.Trans("err_jwt_token_empty")
	case ErrJwtHeaderMalformed:
		return i18n.Obj.Trans("err_jwt_header_malformed")
	case ErrOriginPasswordInvalid:
		return i18n.Obj.Trans("err_origin_password_invalid")

	case ErrUnauthorized:
		return i18n.Obj.Trans("err_unauthorized")
	case ErrNotFound:
		return i18n.Obj.Trans("err_not_found")
	case ErrDNSTxtRecordMissing:
		return "err_dns_txt_record_missing"
	case ErrUnprocessableEntity:
		return i18n.Obj.Trans("err_unprocessable_entity")
	case ErrApiSignatureInvalid:
		return i18n.Obj.Trans("err_api_signature_invalid")

	default:
		return ""
	}
}
