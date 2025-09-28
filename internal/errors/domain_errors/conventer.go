package domainerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Return true if log is needed, and error ofc
func ToGRPCStatus(err error) (bool, error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return false, status.Error(codes.NotFound, "user not found")
	case errors.Is(err, ErrEmailAlreadyExists):
		return false, status.Error(codes.AlreadyExists, "email is already registered")
	case errors.Is(err, ErrInvalidInput):
		return false, status.Error(codes.InvalidArgument, "invalid input")
	case errors.Is(err, ErrPasswordMismatch):
		return false, status.Error(codes.InvalidArgument, "passwords do not match")
	case errors.Is(err, ErrInvalidCredentials):
		return false, status.Error(codes.Unauthenticated, "ivalid credentials")
	default:
		return true, status.Error(codes.Internal, "unknown internal error")
	}
}
