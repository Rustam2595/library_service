package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Rustam2595/library_service/internal/domain/models"
	"github.com/Rustam2595/library_service/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

func (r *Repository) SaveUser(user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	UID := uuid.NewString()
	_, err := r.conn.Exec(ctx, "INSERT INTO Users(uid, name, email, pass) VALUES($1, $2, $3, $4)",
		UID, user.Name, user.Email, user.Pass)
	if err != nil {
		return "", err
	}
	return UID, nil
}

func (r *Repository) ValidateUser(user models.User) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	row := r.conn.QueryRow(ctx, "SELECT uid, pass FROM Users WHERE email = $1 AND uid = $2 AND deleted_user = false",
		user.Email, user.UID)
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
	rows, err := r.conn.Query(ctx, "SELECT * FROM Users WHERE deleted_user = false")
	if err != nil {
		return nil, err
	}
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UID, &user.Name, &user.Email, &user.Pass, &user.DeletedUser); err != nil {
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
	_, err := r.conn.Exec(ctx,
		"UPDATE Users SET name = $1, email = $2, pass = $3 WHERE uid = $4 AND deleted_user = false",
		user.Name, user.Email, user.Pass, uid)
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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err = transaction.Rollback(ctx); err != nil {
			return
		}
	}()
	if _, err = transaction.Prepare(ctx,
		"update user",
		"UPDATE Users SET deleted_user = true WHERE uid = $1"); err != nil {
		return err
	}
	result, err := transaction.Exec(ctx, "update user", uid)
	if err != nil {
		return ErrUserNotFound
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *Repository) DeleteUsers() error {
	zLog := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	result, err := r.conn.Exec(ctx, "DELETE FROM users WHERE deleted_user = true")
	if err != nil {
		zLog.Error().Err(err).Msg("deleted users failed")
		return err
	}
	deletedCount := result.RowsAffected()
	zLog.Debug().Msgf("%d users deleted!", deletedCount)
	return nil
}

func (r *Repository) GetBooks() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	rows, err := r.conn.Query(ctx, "SELECT * FROM Books WHERE deleted = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.BID, &book.Label, &book.Author, &book.Deleted, &book.UserUID, &book.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if len(books) == 0 {
		return nil, ErrBooksListEmpty
	}
	return books, nil
}

func (r *Repository) GetBookByID(bid string) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	row := r.conn.QueryRow(ctx, "SELECT * FROM Books WHERE bid = $1", bid)
	var book models.Book
	if err := row.Scan(&book.BID, &book.Label, &book.Author, &book.Deleted, &book.UserUID, &book.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, fmt.Errorf("book with id = %s, does not exist", bid)
		}
		return models.Book{}, ErrBookNotFound
	}
	if book.Deleted {
		return models.Book{}, ErrBookWasDeleted
	}
	return book, nil
}

func (r *Repository) GetBookByUID(uid string) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	rows, err := r.conn.Query(ctx,
		"SELECT bid,label,author,deleted,user_uid,created_at FROM Books WHERE deleted = false AND user_uid = $1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//books := make([]models.Book, 0)
	//for rows.Next() {
	//	var book models.Book
	//	if err := rows.Scan(&book.BID, &book.Label, &book.Author, &book.Deleted, &book.UserUID, &book.CreatedAt);
	//	err != nil {
	//		return nil, err
	//	}
	//	books = append(books, book)
	//}
	//if err := rows.Err(); err != nil {
	//	return nil, fmt.Errorf("error iterating rows: %w", err)
	//}
	// Автоматически сканирует все строки в слайс
	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Book])
	if err != nil {
		return nil, fmt.Errorf("failed to collect books: %w", err)
	}
	if len(books) == 0 {
		return nil, ErrBooksListEmpty
	}
	return books, nil
}

func (r *Repository) SaveBook(book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	_, err := r.conn.Exec(ctx, "INSERT INTO Books VALUES($1, $2, $3, $4, $5, $6)",
		uuid.NewString(), book.Label, book.Author, book.Deleted, book.UserUID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteBook(bid string) error {
	zLog := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	transaction, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err = transaction.Rollback(ctx); err != nil {
			return
		}
	}()
	if _, err = transaction.Prepare(ctx,
		"update book",
		"UPDATE Books SET deleted = true WHERE bid = $1"); err != nil {
		return err
	}
	result, err := transaction.Exec(ctx, "update book", bid)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrBookNotFound
	}
	zLog.Debug().Msgf("book id = %s, deleted = %t", bid, true)

	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *Repository) DeleteBooks() error {
	zLog := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	res, err := r.conn.Exec(ctx, "DELETE FROM Books WHERE deleted = true")
	if err != nil {
		zLog.Error().Err(err).Msg("deleted books failed")
		return err
	}
	deletedCount := res.RowsAffected()
	zLog.Debug().Msgf("%d books deleted!", deletedCount)
	return nil
}

func Migrations(dbAddr, migrationsPath string) error {
	migratePath := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(migratePath, dbAddr)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	// Применение всех миграций
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}
	log.Println("Migrations successfully created")
	return nil
}
