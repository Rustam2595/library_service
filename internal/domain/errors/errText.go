package errors

const (
	// InvalidAuthDataError возвращается, когда переданы некорректные данные аутентификации (логин или пароль).
	InvalidAuthDataError = "invalid password"

	// UserNotFoundError указывает, что пользователь с указанными данными не найден в системе.
	UserNotFoundError = "user not found"

	// UserListEmptyError сигнализирует, что в базе пользователей нет ни одной записи.
	UserListEmptyError = "user database is empty"

	// BookNotFoundError возвращается, когда книга с указанным идентификатором не найдена.
	BookNotFoundError = "book not found"

	// BooksListEmptyError означает, что в базе книг отсутствуют какие‑либо записи.
	BooksListEmptyError = "book database is empty"

	// BookWasDeletedError указывает, что запрошенная книга была ранее удалена и недоступна.
	BookWasDeletedError = "the book has been deleted"
)
