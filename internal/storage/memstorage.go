package storage

import (
	"github.com/Rustam2595/library_service/internal/domain/models"
	"github.com/google/uuid"
)

type MemStorage struct {
	UsersMap map[string]models.User
	BooksMap map[string]models.Book
}

func New() *MemStorage {
	uMap := make(map[string]models.User)
	bMap := make(map[string]models.Book)
	return &MemStorage{
		UsersMap: uMap,
		BooksMap: bMap,
	}
}

func (ms *MemStorage) SaveUser(user models.User) error {
	uid := uuid.NewString()
	ms.UsersMap[uid] = user
	return nil
}
func (ms *MemStorage) ValidateUser(user models.User) (string, string, error) {
	for uid, value := range ms.UsersMap {
		if value.Email == user.Email {
			if value.Pass != user.Pass {
				return "", "", ErrInvalidAuthData
			}
			return uid, value.Pass, nil
		}
	}
	return "", "", ErrUserNotFound
}
func (ms *MemStorage) GetUsers() ([]models.User, error) {
	var users []models.User
	for uid, e := range ms.UsersMap {
		e.UID = uid
		users = append(users, e)
	}
	if len(users) == 0 {
		return nil, ErrUserListEmpty
	}
	return users, nil
}
func (ms *MemStorage) UpdateUser(uid string, user models.User) error {
	if _, ok := ms.UsersMap[uid]; !ok {
		return ErrUserNotFound
	}
	ms.UsersMap[uid] = user
	return nil
}
func (ms *MemStorage) DeleteUser(uid string) error {
	if _, ok := ms.UsersMap[uid]; !ok {
		return ErrUserNotFound
	}
	delete(ms.UsersMap, uid)
	return nil
}
func (ms *MemStorage) GetBooks() ([]models.Book, error) {
	var books []models.Book
	for bid, e := range ms.BooksMap {
		e.BID = bid
		books = append(books, e)
	}
	if len(books) == 0 {
		return nil, ErrBooksListEmpty
	}
	return books, nil
}

func (ms *MemStorage) GetBookById(bid string) (models.Book, error) {
	if book, ok := ms.BooksMap[bid]; ok {
		book.BID = bid
		return book, nil
	}
	return models.Book{}, ErrBookNotFound
}

func (ms *MemStorage) SaveBook(book models.Book) error {
	nid := uuid.NewString()
	ms.BooksMap[nid] = book
	return nil
}

func (ms *MemStorage) DeleteBook(bid string) error {
	if _, ok := ms.BooksMap[bid]; !ok {
		return ErrBookNotFound
	}
	delete(ms.BooksMap, bid)
	return nil
}
