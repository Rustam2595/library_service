package server

import (
	"errors"
	"fmt"
	"github.com/Rustam2595/library_service/internal/domain/models"
	"github.com/Rustam2595/library_service/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var secretKey = []byte("VerySecretKey2000")

type Claims struct {
	UserID string //`json:"user_id"`
	//Username string `json:"username"`
	//Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Storage interface {
	SaveUser(models.User) (string, error)
	ValidateUser(models.User) (string, string, error)
	GetUsers() ([]models.User, error)
	UpdateUser(string, models.User) error
	DeleteUser(string) error
	GetBooks() ([]models.Book, error)
	GetBookById(string) (models.Book, error)
	SaveBook(models.Book) error
	DeleteBook(string) error
}
type Server struct {
	host    string
	storage Storage
}

func New(host string, storage Storage) *Server {
	return &Server{
		host:    host,
		storage: storage,
	}
}
func (s *Server) Run() error {
	r := gin.Default()
	//r.Use(gin.Recovery())
	//r.Use(gin.Logger())
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", s.RegisterHandler)
		userGroup.POST("/auth", s.AuthHandler)
		userGroup.GET("/get_all_users", s.AllUsersHandler)
		userGroup.PUT("/update_user/:id", s.UpdateUserHandler)
		userGroup.DELETE("/delete/:id", s.DeleteUserHandler)
	}
	bookGroup := r.Group("/book")
	{
		bookGroup.GET("/all_books", s.AllBooksHandler)
		bookGroup.GET("/:id", s.GetBookByIdHandler)
		bookGroup.POST("/add_book", s.SaveBookHandler)
		bookGroup.DELETE("/delete/:id", s.DeleteBookHandler)
	}
	if err := r.Run(s.host); err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Pass = string(passHash)
	val := validator.New()
	if err := val.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.storage.SaveUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//add JWT
	token, err := createJWT(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("authorization", token)
	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
}

func (s *Server) AuthHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	val := validator.New()
	if err := val.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, pass, err := s.storage.ValidateUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidAuthData) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		//errors.As()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Pass)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Hash password"})
		return
	}

	//add token
	token, err := createJWT(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("authorization", token)
	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully registered, uuid = " + uid + " , PASS = " + pass})
}

func (s *Server) AllUsersHandler(ctx *gin.Context) {
	users, err := s.storage.GetUsers()
	if err != nil {
		if errors.Is(err, storage.ErrUserListEmpty) {
			ctx.String(http.StatusNoContent, "There are no users here!")
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
func (s *Server) UpdateUserHandler(ctx *gin.Context) {
	var user models.User
	uid := ctx.Param("id")
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.storage.UpdateUser(uid, user); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully updated"})
}

func (s *Server) DeleteUserHandler(ctx *gin.Context) {
	uid := ctx.Param("id")
	if err := s.storage.DeleteUser(uid); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			ctx.JSON(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}

func (s *Server) AllBooksHandler(ctx *gin.Context) {
	books, err := s.storage.GetBooks()
	if err != nil {
		if errors.Is(err, storage.ErrBooksListEmpty) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (s *Server) GetBookByIdHandler(ctx *gin.Context) {
	bid := ctx.Param("id")
	book, err := s.storage.GetBookById(bid)
	if err != nil {
		if errors.Is(err, storage.ErrBookNotFound) {
			ctx.JSON(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (s *Server) SaveBookHandler(ctx *gin.Context) {
	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	val := validator.New()
	if err := val.Struct(book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.storage.SaveBook(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Book successfully saved"})
}

func (s *Server) DeleteBookHandler(ctx *gin.Context) {
	bid := ctx.Param("id")
	if err := s.storage.DeleteBook(bid); err != nil {
		if errors.Is(err, storage.ErrBookNotFound) {
			ctx.JSON(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book successfully deleted"})
}

// СОЗДАНИЕ токена (при логине)
func createJWT(UID string) (string, error) {
	// Данные для токена
	claims := Claims{
		UserID: UID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", UID),
		},
	}
	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем секретным ключом
	return token.SignedString(secretKey)
}

func validJWT(tokenString string) (*Claims, error) {
	// Парсим и проверяем токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// Проверяем валидность
	if claim, ok := token.Claims.(*Claims); ok && token.Valid {
		return claim, nil
	}
	return nil, fmt.Errorf("невалидный токен")
}
