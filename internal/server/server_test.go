package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rustam2595/library_service/internal/domain/models"
	"github.com/Rustam2595/library_service/internal/storage"
	"github.com/Rustam2595/library_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegisterHandler(t *testing.T) {
	srv := Server{
		validator: validator.New(),
	}
	r := gin.Default()
	r.POST("/register", srv.RegisterHandler)
	httpSrv := httptest.NewServer(r)
	type want struct {
		errFlag    bool
		mockFlag   bool
		statusCode int
	}
	testCases := []struct {
		name    string
		method  string
		request string
		user    string
		uid     string
		err     error
		want    want
	}{
		{
			name:    "Test RegisterHandler() func; Case 1:",
			method:  http.MethodPost,
			request: "/register",
			user:    `{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "testUID",
			err:     nil,
			want: want{
				errFlag:    false,
				mockFlag:   true,
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "Test RegisterHandler() func; Case 2:",
			method:  http.MethodPost,
			request: "/register",
			user:    `{"uid":"uid","name":"Sergei","email":"testemailya.ru","pass":"qwerty1234","deleted_user":false}`,
			//uid:     "2testUID",
			//err:     errors.New("validator error"),
			want: want{
				errFlag:    true,
				mockFlag:   false,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Test RegisterHandler() func; Case 3:",
			method:  http.MethodPost,
			request: "/register",
			user:    `{"uid":"uid","email":"testemailya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "",
			err:     errors.New("name required error"),
			want: want{
				errFlag:    true,
				mockFlag:   false,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Test RegisterHandler() func; Case 4:",
			method:  http.MethodPost,
			request: "/register",
			user:    `{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "",
			err:     errors.New("save user error"),
			want: want{
				errFlag:    true,
				mockFlag:   true,
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockStorage(ctrl)
			if tc.want.mockFlag {
				mockRepo.EXPECT().SaveUser(gomock.Any()).Return(tc.uid, tc.err)
				srv.storage = mockRepo
			}
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			req.Body = tc.user
			response, err := req.Send()
			if !tc.want.errFlag {
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Header().Get("Authorization"))
			}
			assert.Equal(t, tc.want.statusCode, response.StatusCode())

		})
	}
}

func TestAuthHandler(t *testing.T) {
	srv := Server{
		validator: validator.New(),
	}
	r := gin.Default()
	r.POST("/auth", srv.AuthHandler)
	httpSrv := httptest.NewServer(r)
	type want struct {
		errFlag    bool
		mockFlag   bool
		statusCode int
	}
	testCases := []struct {
		name    string
		method  string
		request string
		user    string
		uid     string
		pass    string
		err     error
		want    want
	}{
		{
			name:    "Test AuthHandler() func; Case 1: OK",
			method:  http.MethodPost,
			request: "/auth",
			user:    `{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "uid",
			pass:    "qwerty1234",
			err:     nil,
			want: want{
				errFlag:    false,
				mockFlag:   true,
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "Test AuthHandler() func; Case 2: validator",
			method:  http.MethodPost,
			request: "/auth",
			user:    `{"uid":"uid","name":"Sergei","email":"testemailya.ru","pass":"qwerty1234","deleted_user":false}`,
			//uid:     "2testUID",
			//pass:  "pass"
			//err:     errors.New("validator error"),
			want: want{
				errFlag:    true,
				mockFlag:   false,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Test AuthHandler() func; Case 3: comparePass",
			method:  http.MethodPost,
			request: "/auth",
			user:    `{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "uid",
			pass:    "pass",
			err:     nil,
			want: want{
				errFlag:    true,
				mockFlag:   true,
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:    "Test AuthHandler() func; Case 4: validateUser",
			method:  http.MethodPost,
			request: "/auth",
			user:    `{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}`,
			uid:     "",
			pass:    "qwerty1234",
			err:     errors.New("validate user error"),
			want: want{
				errFlag:    true,
				mockFlag:   true,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range testCases {
		passHash, err := bcrypt.GenerateFromPassword([]byte(tc.pass), bcrypt.DefaultCost)
		assert.NoError(t, err)
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockStorage(ctrl)
			if tc.want.mockFlag {
				mockRepo.EXPECT().ValidateUser(gomock.Any()).Return(tc.uid, string(passHash), tc.err)
				srv.storage = mockRepo
			}
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			req.Body = tc.user
			response, err := req.Send()
			if !tc.want.errFlag {
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Header().Get("Authorization"))
			}
			assert.Equal(t, tc.want.statusCode, response.StatusCode())

		})
	}
}

func TestAllUserHandler(t *testing.T) {
	srv := Server{
		validator: validator.New(),
	}
	r := gin.Default()
	r.GET("/get_all_users", srv.AllUsersHandler)
	httpSrv := httptest.NewServer(r)
	type want struct {
		errFlag    bool
		users      string
		statusCode int
	}
	testCases := []struct {
		name    string
		method  string
		request string
		users   []models.User
		err     error
		want    want
	}{
		{
			name:    "Test AuthHandler() func; Case 1: OK",
			method:  http.MethodGet,
			request: "/get_all_users",
			err:     nil,
			users: []models.User{
				{
					UID:         "uid",
					Name:        "Sergei",
					Email:       "testemail@ya.ru",
					Pass:        "qwerty1234",
					DeletedUser: false,
				},
			},
			want: want{
				errFlag:    false,
				users:      `[{"uid":"uid","name":"Sergei","email":"testemail@ya.ru","pass":"qwerty1234","deleted_user":false}]`,
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "Test AuthHandler() func; Case 2: userListEmpty",
			method:  http.MethodGet,
			request: "/get_all_users",
			err:     storage.ErrUserListEmpty,
			users:   nil,
			want: want{
				errFlag:    true,
				users:      "",
				statusCode: http.StatusNoContent,
			},
		},
		{
			name:    "Test AuthHandler() func; Case 3: getUsersError",
			method:  http.MethodGet,
			request: "/get_all_users",
			err:     errors.New("get users error"),
			users:   nil,
			want: want{
				errFlag:    true,
				users:      `{"error":"get users error"}`,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockStorage(ctrl)
			mockRepo.EXPECT().GetUsers().Return(tc.users, tc.err)
			srv.storage = mockRepo
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			response, err := req.Send()
			if !tc.want.errFlag {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want.statusCode, response.StatusCode())
			assert.Equal(t, tc.want.users, string(response.Body()))
		})
	}
}

func TestDeleter(t *testing.T) {
	type want struct {
		err error
	}
	type test struct {
		name string
		want want
	}
	tests := []test{
		{
			name: "Test deleter func; Case 1:",
			want: want{
				err: nil,
			},
		},
		{
			name: "Test deleter func; Case 2:",
			want: want{
				err: fmt.Errorf("test err"),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			ctrl := gomock.NewController(t)
			m := mocks.NewMockStorage(ctrl)
			defer ctrl.Finish()
			m.EXPECT().DeleteBooks().Return(tc.want.err)
			srv := New("0.0.0.0:8080", m)
			for i := 0; i < 2; i++ {
				srv.deleteChan <- i
			}
			go srv.Deleter(ctx)
			for {
				select {
				case err := <-srv.ErrChan:
					assert.Equal(t, tc.want.err, err)
					return
				case <-time.After(time.Second):
					if tc.want.err != nil {
						t.Fatalf("Exp err = %s; actual = nil", tc.want.err)
						return
					}
					return
				}
			}
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	srv := &Server{
		validator: validator.New(),
	}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/update_user/:id", srv.UpdateUserHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()
	type want struct {
		statusCode   int
		expectedBody string
	}
	testCases := []struct {
		name        string
		uid         string
		requestBody string
		mockSetup   func(*mocks.MockStorage)
		want        want
	}{
		{
			name:        "Test UpdateUserHandler() func; Case 1: успешное обновление",
			uid:         "uid",
			requestBody: `{"name":"Updated Name","email":"updated@example.com","pass":"newpassword123"}`,
			mockSetup: func(m *mocks.MockStorage) {
				m.EXPECT().UpdateUser("uid", gomock.Any()).Return(nil).Times(1)
			},
			want: want{
				statusCode:   http.StatusOK,
				expectedBody: "successfully updated",
			},
		},
		{
			name:        "Test UpdateUserHandler() func; Case 2: невалидный json",
			uid:         "uid",
			requestBody: `{"name""Updated Name","email":"updated@example.com","pass":"newpassword123"}`,
			mockSetup: func(m *mocks.MockStorage) {
				m.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			want: want{
				statusCode:   http.StatusBadRequest,
				expectedBody: "error",
			},
		},
		{
			name:        "Test UpdateUserHandler() func; Case 3: err UserNotFound",
			uid:         "Updated Name",
			requestBody: `{"name":"Updated Name","email":"updated@example.com","pass":"newpassword123"}`,
			mockSetup: func(m *mocks.MockStorage) {
				m.EXPECT().UpdateUser("Updated Name", gomock.Any()).Return(storage.ErrUserNotFound).Times(1)
			},
			want: want{
				statusCode:   http.StatusNotFound,
				expectedBody: "not found",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStorage := mocks.NewMockStorage(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockStorage)
			}
			srv.storage = mockStorage
			req := resty.New().R()
			req.Method = http.MethodPut
			req.URL = fmt.Sprintf("%s/update_user/%s", httpSrv.URL, tc.uid) //httpSrv.URL + tc.request
			req.Body = tc.requestBody
			resp, err := req.Send()

			//// 6. Создаём resty клиент
			//client := resty.New()
			//// 7. Формируем правильный URL с подстановкой ID
			//url := fmt.Sprintf("%s/update_user/%s", httpSrv.URL, tc.uid)
			//// 8. Отправляем запрос
			//resp, err := client.R().
			//	SetHeader("Content-Type", "application/json").
			//	SetBody(tc.requestBody).
			//	Put(url) // ✅ Правильный URL с реальным ID

			log.Println(resp.String())
			assert.NoError(t, err)
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Contains(t, resp.String(), tc.want.expectedBody)

		})
	}
}
