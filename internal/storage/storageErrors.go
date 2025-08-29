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

//	SaveUser(models.User) error
//	ValidateUser(models.User) (string, string, error)
//	GetUsers() ([]models.User, error)
//	UpdateUser(string, models.User) error
//	DeleteUser(string) error
//	GetBooks() ([]models.Book, error)
//	GetBookById(string) (models.Book, error)
//	SaveBook(models.Book) error
//	DeleteBook(string) error
