package rqres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromError(err error) *status.Status {
	logErr := errx.Logic{}
	if errors.As(err, &logErr) {
		return status.Newf(codes.InvalidArgument, "[%d] %s", logErr.Code(), logErr.Error())
	}
	invErr := errx.Invalid{}
	if errors.As(err, &invErr) {
		var sb strings.Builder
		sb.WriteString(invErr.Error())
		for _, ve := range invErr.Errors() {
			sb.WriteString(fmt.Sprintf("[%s] %s", ve.Field, ve.Err.Error()))
		}
		return status.Newf(codes.InvalidArgument, sb.String())
	}
	nfErr := errx.NotFound{}
	if errors.As(err, &nfErr) {
		return status.New(codes.NotFound, invErr.Error())
	}
	base := errx.Base{}
	if errors.As(err, &base) {
		switch base.Kind() {
		case errx.TypePerms:
			return status.New(codes.PermissionDenied, base.Error())
		case errx.TypeFatal:
			return status.New(codes.Internal, base.Error())
		}
	}
	return status.New(codes.InvalidArgument, err.Error())
}
