package domain

import (
	pkghttp "auth_service/pkg/http"
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidToken     = errors.New("invalid token")
	ErrMissingParameter = errors.New("missing uuid parameter")
)

func HandleError(err error, r any) pkghttp.Response {
	switch {
	case err == nil:
		return pkghttp.OK(r)
	case errors.Is(err, ErrNotFound):
		return pkghttp.NotFound(err)
	case errors.Is(err, ErrUnauthorized):
		return pkghttp.Unauthorized(err)
	case errors.Is(err, ErrInvalidToken):
		return pkghttp.Unauthorized(err)
	case errors.Is(err, ErrMissingParameter):
		return pkghttp.BadRequest(err)
	default:
		return pkghttp.Unknown(err)
	}
}
