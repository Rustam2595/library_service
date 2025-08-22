package storage

import (
	"errors"
	errMess "github.com/Rustam2595/library_service/internal/domain/errors"
)

var ErrInvalidAuthData = errors.New(errMess.InvalidAuthDataError)
var ErrUserNotFound = errors.New(errMess.UserNotFoundError)
var ErrUserListEmpty = errors.New(errMess.UserListEmptyError)
var ErrBookNotFound = errors.New(errMess.BookNotFoundError)
var ErrBooksListEmpty = errors.New(errMess.BooksListEmptyError)
var ErrBookWasDeleted = errors.New(errMess.BookWasDeletedError)
