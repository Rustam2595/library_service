package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rustam2595/library_service/internal/domain/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const ctxTimeout = 2 * time.Second

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepo(ctx context.Context, dbAddr string) (*Repository, error) {
	conn, err := pgxpool.New(ctx, dbAddr)
	if err != nil {
		return nil, err
	}
	return &Repository{
		conn: conn,
	}, nil
}

func (r *Repository) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	_, err := r.conn.Exec(ctx, "INSERT INTO Users(uid, name, email, pass) VALUES($1, $2, $3, $4)", uuid.NewString(), user.Name, user.Email, user.Pass)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ValidateUser(user models.User) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	row := r.conn.QueryRow(ctx, "SELECT uid, pass FROM Users WHERE email = $1 AND uid = $2", user.Email, user.UID)
	var pass, uid string
	if err := row.Scan(&uid, &pass); err != nil {
		return "", "", ErrUserNotFound
	}
	if pass != user.Pass {
		return "", "", ErrInvalidAuthData
	}
	return uid, pass, nil
}

func (r *Repository) GetUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	rows, err := r.conn.Query(ctx, "SELECT * FROM Users")
	if err != nil {
		return nil, err
	}
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UID, &user.Name, &user.Email, &user.Pass); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, ErrUserListEmpty
	}
	return users, nil
}
func (r *Repository) UpdateUser(uid string, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	_, err := r.conn.Exec(ctx, "UPDATE Users SET name = $1, email = $2, pass = $3 WHERE uid = $4)", user.Name, user.Email, user.Pass, uid)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}
func (r *Repository) DeleteUser(uid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	transaction, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)
	if _, err := transaction.Prepare(ctx, "delete user", "DELETE FROM Users WHERE uid = $1"); err != nil {
		return err
	}
	if _, err := transaction.Exec(ctx, "delete user", uid); err != nil {
		return ErrUserNotFound
	}
	return transaction.Commit(ctx)
}

func (r *Repository) GetBooks() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	rows, err := r.conn.Query(ctx, "SELECT * FROM Books")
	if err != nil {
		return nil, err
	}
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.BID, &book.Label, &book.Author, &book.Deleted, &book.User_UID, &book.Created_at); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if len(books) == 0 {
		return nil, ErrBooksListEmpty
	}
	return books, nil
}

func (r *Repository) GetBookById(bid string) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	row := r.conn.QueryRow(ctx, "SELECT * FROM Books WHERE bid = $1", bid)
	var book models.Book
	if err := row.Scan(&book.BID, &book.Label, &book.Author, &book.Deleted, &book.User_UID, &book.Created_at); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, fmt.Errorf("book with id = %s, does not exist", bid)
		}
		return models.Book{}, ErrBookNotFound
	}
	return book, nil
}

func (r *Repository) SaveBook(book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	_, err := r.conn.Exec(ctx, "INSERT INTO Books VALUES($1, $2, $3, $4, $5, $6)", uuid.NewString(), book.Label, book.Author, book.Deleted, book.User_UID, book.Created_at)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteBook(bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	transaction, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)
	if _, err := transaction.Prepare(ctx, "delete book", "DELETE FROM Books WHERE bid = $1"); err != nil {
		return err
	}
	if _, err := transaction.Exec(ctx, "delete book", bid); err != nil {
		return ErrBookNotFound
	}
	return transaction.Commit(ctx)
}

func Migrations(dbAddr, migrationsPath string) error {
	migratePath := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(migratePath, dbAddr)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	// Применение всех миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}
	log.Println("Migrations successfully created")
	return nil
}
