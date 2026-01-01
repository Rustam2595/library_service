package storage

import (
	"errors"

	errMess "github.com/Rustam2595/library_service/internal/domain/errors"
)

// ErrInvalidAuthData возвращается, когда переданы некорректные учётные данные (логин/пароль).
var ErrInvalidAuthData = errors.New(errMess.InvalidAuthDataError)

// ErrUserNotFound сигнализирует, что пользователь с указанными параметрами не найден.
var ErrUserNotFound = errors.New(errMess.UserNotFoundError)

// ErrUserListEmpty означает, что в хранилище отсутствуют какие‑либо пользователи.
var ErrUserListEmpty = errors.New(errMess.UserListEmptyError)

// ErrBookNotFound возвращается, когда книга с указанным идентификатором не найдена.
var ErrBookNotFound = errors.New(errMess.BookNotFoundError)

// ErrBooksListEmpty указывает, что в хранилище нет ни одной книги.
var ErrBooksListEmpty = errors.New(errMess.BooksListEmptyError)

// ErrBookWasDeleted означает, что запрошенная книга была ранее удалена и недоступна.
var ErrBookWasDeleted = errors.New(errMess.BookWasDeletedError)
