package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mhaatha/go-bookshelf/internal/config"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockBookService struct {
	// CreateNewBook
	CreateMockRequest  web.CreateBookRequest
	CreateMockResponse web.CreateBookResponse

	// GetAllBooks
	GetAllMockQuery    web.QueryParamsGetBooks
	GetAllMockResponse []web.GetBookResponse

	// GetBookById
	GetByIdMockPathValue web.PathParamsGetBook
	GetByIdMockResponse  web.GetBookResponse

	// UpdateBookById
	UpdateByIdMockPathValue web.PathParamsUpdateBook
	UpdateByIdMockRequest   web.UpdateBookRequest
	UpdateByIdMockResponse  web.UpdateBookResponse

	// DeleteBookById
	DeleteByIdMockPathValue web.PathParamsDeleteBook

	MockError error
}

func (m *MockBookService) CreateNewBook(ctx context.Context, request web.CreateBookRequest) (web.CreateBookResponse, error) {
	m.CreateMockRequest = request

	if m.MockError != nil {
		return web.CreateBookResponse{}, m.MockError
	}

	return m.CreateMockResponse, nil
}

func (m *MockBookService) GetAllBooks(ctx context.Context, queries web.QueryParamsGetBooks) ([]web.GetBookResponse, error) {
	m.GetAllMockQuery = queries

	if m.MockError != nil {
		return nil, m.MockError
	}

	return m.GetAllMockResponse, nil
}

func (m *MockBookService) GetBookById(ctx context.Context, pathValues web.PathParamsGetBook) (web.GetBookResponse, error) {
	m.GetByIdMockPathValue = pathValues

	if m.MockError != nil {
		return web.GetBookResponse{}, m.MockError
	}

	return m.GetByIdMockResponse, nil
}

func (m *MockBookService) UpdateBookById(ctx context.Context, pathValues web.PathParamsUpdateBook, request web.UpdateBookRequest) (web.UpdateBookResponse, error) {
	m.UpdateByIdMockPathValue = pathValues
	m.UpdateByIdMockRequest = request

	if m.MockError != nil {
		return web.UpdateBookResponse{}, m.MockError
	}

	return m.UpdateByIdMockResponse, nil
}

func (m *MockBookService) DeleteBookById(ctx context.Context, pathValues web.PathParamsDeleteBook) error {
	m.DeleteByIdMockPathValue = pathValues

	if m.MockError != nil {
		return m.MockError
	}

	return nil
}

func TestBookCreateHandler(t *testing.T) {
	t.Run("create book with complete data", func(t *testing.T) {
		bookRequest := web.CreateBookRequest{
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
			Status:        "completed",
			CompletedDate: "2025-09-29",
		}
		expectedServiceResponse := web.CreateBookResponse{
			Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
			Status:        "completed",
			CompletedDate: "2025-09-29",
		}

		mockService := &MockBookService{
			CreateMockResponse: expectedServiceResponse,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Check status code
		if res.Code != http.StatusCreated {
			t.Errorf("expected status code of %d but got %d", http.StatusCreated, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebSuccessResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body message
		if actualResponseBody.Message != "Book created successfully" {
			t.Errorf("expected '%s' as response message but got '%s'", "Book created successfully", actualResponseBody.Message)
		}

		// Check response body data
		val, ok := actualResponseBody.Data.(map[string]interface{})
		if ok {
			if val["id"] != expectedServiceResponse.Id {
				t.Errorf("expected id '%s' but got '%s'", expectedServiceResponse.Id, val["id"])
			}

			if val["name"] != expectedServiceResponse.Name {
				t.Errorf("expected name '%s' but got '%s'", expectedServiceResponse.Name, val["name"])
			}

			if int(val["total_page"].(float64)) != expectedServiceResponse.TotalPage {
				t.Errorf("expected total_page '%d' but got '%d'", expectedServiceResponse.TotalPage, val["total_page"])
			}

			if val["author_id"] != expectedServiceResponse.AuthorId {
				t.Errorf("expected author_id '%s' but got '%s'", expectedServiceResponse.AuthorId, val["author_id"])
			}

			if val["photo_key"] != expectedServiceResponse.PhotoKey {
				t.Errorf("expected photo_key '%s' but got '%s'", expectedServiceResponse.PhotoKey, val["photo_key"])
			}

			if val["status"] != expectedServiceResponse.Status {
				t.Errorf("expected status '%s' but got '%s'", expectedServiceResponse.Status, val["status"])
			}

			if val["completed_date"] != expectedServiceResponse.CompletedDate {
				t.Errorf("expected completed_date '%s' but got '%s'", expectedServiceResponse.CompletedDate, val["completed_date"])
			}
		} else {
			t.Errorf("val should be true but got false")
		}

		// Check actual request body that has been passed to service
		if !reflect.DeepEqual(mockService.CreateMockRequest, bookRequest) {
			t.Errorf("expected %+v as request body but got %v", bookRequest, mockService.CreateMockRequest)
		}
	})

	t.Run("create book with invalid name", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "required",
				BookRequest: web.CreateBookRequest{
					Name:          "",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "name",
				ErrMessage: "name is required",
			},
			{
				Name: "minimum length",
				BookRequest: web.CreateBookRequest{
					Name:          "La",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "name",
				ErrMessage: "name must be at least 3 characters",
			},
			{
				Name: "maximum length",
				BookRequest: web.CreateBookRequest{
					Name:          "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "name",
				ErrMessage: "name must be at most 255 characters",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with invalid total_page", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "required",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "total_page",
				ErrMessage: "total_page is required",
			},
			{
				Name: "minimum number",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     -1,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "total_page",
				ErrMessage: "total_page must be at least 1 characters",
			},
			{
				Name: "maximum number",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     12001,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "total_page",
				ErrMessage: "total_page must be at most 12000 characters",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with invalid author_id", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "required",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "author_id",
				ErrMessage: "author_id is required",
			},
			{
				Name: "uuid",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "InvalidUUID",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "author_id",
				ErrMessage: "'InvalidUUID' is not a valid UUID",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with invalid photo_key", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "required",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "photo_key",
				ErrMessage: "photo_key is required",
			},
			{
				Name: "minimum length",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "a.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "photo_key",
				ErrMessage: "photo_key must be at least 6 characters",
			},
			{
				Name: "maximum length",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.jpg",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "photo_key",
				ErrMessage: "photo_key must be at most 255 characters",
			},
			{
				Name: "valid photo key",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "InvalidPhotoKey",
					Status:        "completed",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "photo_key",
				ErrMessage: "'InvalidPhotoKey' is not a valid photo key",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with invalid status", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "required",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "status",
				ErrMessage: "status is required",
			},
			{
				Name: "invalid book status",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "InvalidStatus",
					CompletedDate: "2025-09-29",
				},
				ErrField:   "status",
				ErrMessage: "the valid value for this field are only 'completed', 'reading', and 'plan_to_read'",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with invalid completed_date", func(t *testing.T) {
		cases := []struct {
			Name        string
			BookRequest web.CreateBookRequest
			ErrField    string
			ErrMessage  string
		}{
			{
				Name: "datetime",
				BookRequest: web.CreateBookRequest{
					Name:          "Laut Bercerita",
					TotalPage:     379,
					AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
					PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
					Status:        "completed",
					CompletedDate: "InvalidDateTime",
				},
				ErrField:   "completed_date",
				ErrMessage: "use YYYY-MM-DD for valid datetime",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				bookRequest := c.BookRequest
				expectedServiceError := validate.Struct(bookRequest)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})

	t.Run("create book with the same name and author", func(t *testing.T) {
		bookRequest := web.CreateBookRequest{
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
			Status:        "completed",
			CompletedDate: "2025-09-29",
		}
		expectedServiceError := []appError.ErrAggregate{
			appError.ErrAggregate{
				Field:   "name",
				Message: "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists",
			},
		}

		mockService := &MockBookService{
			MockError: appError.NewAppError(
				http.StatusBadRequest,
				expectedServiceError,
				nil,
			),
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "name" {
					t.Errorf("expected error field is %s but got %s", "name", val["field"])
				}

				if val["message"] != "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists" {
					t.Errorf("expected error message is %s but got %s", "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}
	})

	t.Run("create book with not found author_id", func(t *testing.T) {
		bookRequest := web.CreateBookRequest{
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
			Status:        "completed",
			CompletedDate: "2025-09-29",
		}
		expectedServiceError := []appError.ErrAggregate{
			appError.ErrAggregate{
				Field:   "author_id",
				Message: "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found",
			},
		}

		mockService := &MockBookService{
			MockError: appError.NewAppError(
				http.StatusNotFound,
				expectedServiceError,
				nil,
			),
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/books", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		handler.Create(res, req)

		// Check status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code of %d but got %d", http.StatusNotFound, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "author_id" {
					t.Errorf("expected error field is %s but got %s", "author_id", val["field"])
				}
				if val["message"] != "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found" {
					t.Errorf("expected error message is %s but got %s", "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}
	})

	t.Run("create book with invalid JSON payload", func(t *testing.T) {
		cases := []struct {
			Name               string
			InvalidJSONPayload string
			ErrMessage         string
		}{
			{
				Name:               "invalid type name field",
				InvalidJSONPayload: `{"name": 1, "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: name",
			},
			{
				Name:               "invalid type total_page field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": "379", "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: total_page",
			},
			{
				Name:               "invalid type author_id field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": 123, "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: author_id",
			},
			{
				Name:               "invalid type photo_key field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": 123, "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: photo_key",
			},
			{
				Name:               "invalid type status field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": 1, "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: status",
			},
			{
				Name:               "invalid type completed_date field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": 2025}`,
				ErrMessage:         "Invalid JSON type for field: completed_date",
			},
			{
				Name:               "invalid JSON payload",
				InvalidJSONPayload: `{"name": "Laut Bercerita" "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": 2025-10-29}`,
				ErrMessage:         "Invalid JSON payload",
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				mockService := &MockBookService{}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(c.InvalidJSONPayload))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()

				handler.Create(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				val, ok := actualResponseBody.Errors.(string)
				if ok {
					if val != c.ErrMessage {
						t.Errorf("expected %s but got %s", c.ErrMessage, val)
					}
				} else {
					t.Error("val should be true but got false")
				}
			})
		}
	})
}

func TestBookGetAllHandler(t *testing.T) {
	t.Run("get all books", func(t *testing.T) {
		expectedServiceResponse := []web.GetBookResponse{
			{
				Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
				Name:          "Laut Bercerita",
				TotalPage:     379,
				AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
				PhotoURL:      "http://127.0.0.1:9000/book-images/ac0a9b20-2e77-4905-a665-3006763d1935.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
				Status:        "completed",
				CompletedDate: "2025-09-29",
			},
			{
				Id:            "f200a4c1-a141-44a0-9c9d-0b035016e2f9",
				Name:          "Sebuah Seni Untuk Bersikap Bodo Amat",
				TotalPage:     246,
				AuthorId:      "8b970b2a-09d4-450c-8bb8-83da50392d6d",
				PhotoURL:      "http://127.0.0.1:9000/book-images/67f99bdf-4c43-4200-b5d0-a7adbe125f97.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
				Status:        "completed",
				CompletedDate: "2025-11-29",
			},
		}

		mockService := &MockBookService{
			GetAllMockResponse: expectedServiceResponse,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
		res := httptest.NewRecorder()

		handler.GetAll(res, req)

		// Check status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code of %d but got %d", http.StatusOK, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebSuccessResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body message
		if actualResponseBody.Message != "Success get all books" {
			t.Errorf("expected %s as response message but got %s", "Success get all books", actualResponseBody.Message)
		}

		// Check response body data
		dataList, ok := actualResponseBody.Data.([]interface{})
		if ok {
			// First data from JSON array
			val, ok := dataList[0].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[0].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[0].Id, val["id"])
				}

				if val["name"] != expectedServiceResponse[0].Name {
					t.Errorf("expected %s as name but got %s", expectedServiceResponse[0].Name, val["name"])
				}

				if int(val["total_page"].(float64)) != expectedServiceResponse[0].TotalPage {
					t.Errorf("expected %d as total_page but got %d", expectedServiceResponse[0].TotalPage, val["total_page"])
				}

				if val["author_id"] != expectedServiceResponse[0].AuthorId {
					t.Errorf("expected %s as author_id but got %s", expectedServiceResponse[0].AuthorId, val["author_id"])
				}

				if val["photo_url"] != expectedServiceResponse[0].PhotoURL {
					t.Errorf("expected %s as photo_url but got %s", expectedServiceResponse[0].PhotoURL, val["photo_url"])
				}

				if val["status"] != expectedServiceResponse[0].Status {
					t.Errorf("expected %s as status but got %s", expectedServiceResponse[0].Status, val["status"])
				}

				if val["completed_date"] != expectedServiceResponse[0].CompletedDate {
					t.Errorf("expected %s as completed_date but got %s", expectedServiceResponse[0].CompletedDate, val["completed_date"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			// Second data from JSON array
			val, ok = dataList[1].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[1].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[1].Id, val["id"])
				}

				if val["name"] != expectedServiceResponse[1].Name {
					t.Errorf("expected %s as name but got %s", expectedServiceResponse[1].Name, val["name"])
				}

				if int(val["total_page"].(float64)) != expectedServiceResponse[1].TotalPage {
					t.Errorf("expected %d as total_page but got %d", expectedServiceResponse[1].TotalPage, val["total_page"])
				}

				if val["author_id"] != expectedServiceResponse[1].AuthorId {
					t.Errorf("expected %s as author_id but got %s", expectedServiceResponse[1].AuthorId, val["author_id"])
				}

				if val["photo_url"] != expectedServiceResponse[1].PhotoURL {
					t.Errorf("expected %s as photo_url but got %s", expectedServiceResponse[1].PhotoURL, val["photo_url"])
				}

				if val["status"] != expectedServiceResponse[1].Status {
					t.Errorf("expected %s as status but got %s", expectedServiceResponse[1].Status, val["status"])
				}

				if val["completed_date"] != expectedServiceResponse[1].CompletedDate {
					t.Errorf("expected %s as completed_date but got %s", expectedServiceResponse[1].CompletedDate, val["completed_date"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("dataList should be true but got false")
		}
	})

	t.Run("get books by query parameter", func(t *testing.T) {
		cases := []struct {
			Name                    string
			QueryParams             web.QueryParamsGetBooks
			QueryStringURL          string
			ExpectedServiceResponse []web.GetBookResponse
		}{
			{
				Name: "'name' query parameter",
				QueryParams: web.QueryParamsGetBooks{
					Name: "Laut",
				},
				QueryStringURL: "/api/v1/books?name=Laut",
				ExpectedServiceResponse: []web.GetBookResponse{
					{
						Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
						Name:          "Laut Bercerita",
						TotalPage:     379,
						AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
						PhotoURL:      "http://127.0.0.1:9000/book-images/ac0a9b20-2e77-4905-a665-3006763d1935.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
						Status:        "completed",
						CompletedDate: "2025-09-29",
					},
				},
			},
			{
				Name: "'status' query parameter",
				QueryParams: web.QueryParamsGetBooks{
					Status: "completed",
				},
				QueryStringURL: "/api/v1/books?status=completed",
				ExpectedServiceResponse: []web.GetBookResponse{
					{
						Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
						Name:          "Laut Bercerita",
						TotalPage:     379,
						AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
						PhotoURL:      "http://127.0.0.1:9000/book-images/ac0a9b20-2e77-4905-a665-3006763d1935.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
						Status:        "completed",
						CompletedDate: "2025-09-29",
					},
					{
						Id:            "f200a4c1-a141-44a0-9c9d-0b035016e2f9",
						Name:          "Sebuah Seni Untuk Bersikap Bodo Amat",
						TotalPage:     246,
						AuthorId:      "8b970b2a-09d4-450c-8bb8-83da50392d6d",
						PhotoURL:      "http://127.0.0.1:9000/book-images/67f99bdf-4c43-4200-b5d0-a7adbe125f97.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
						Status:        "completed",
						CompletedDate: "2025-11-29",
					},
				},
			},
			{
				Name: "'author_name' query parameter",
				QueryParams: web.QueryParamsGetBooks{
					AuthorName: "Leila",
				},
				QueryStringURL: "/api/v1/books?author_name=Leila",
				ExpectedServiceResponse: []web.GetBookResponse{
					{
						Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
						Name:          "Laut Bercerita",
						TotalPage:     379,
						AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
						PhotoURL:      "http://127.0.0.1:9000/book-images/ac0a9b20-2e77-4905-a665-3006763d1935.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
						Status:        "completed",
						CompletedDate: "2025-09-29",
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				expectedQueries := c.QueryParams
				expectedServiceResponse := c.ExpectedServiceResponse

				mockService := &MockBookService{
					GetAllMockResponse: expectedServiceResponse,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodGet, c.QueryStringURL, nil)
				res := httptest.NewRecorder()

				handler.GetAll(res, req)

				// Check status code
				if res.Code != http.StatusOK {
					t.Errorf("expected status code of %d but got %d", http.StatusOK, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebSuccessResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				// Check response body data
				dataList, ok := actualResponseBody.Data.([]interface{})
				if ok {
					for i, data := range dataList {
						val, ok := data.(map[string]interface{})
						if ok {
							if val["id"] != expectedServiceResponse[i].Id {
								t.Errorf("expected %s as id but got %s", expectedServiceResponse[i].Id, val["id"])
							}

							if val["name"] != expectedServiceResponse[i].Name {
								t.Errorf("expected %s as name but got %s", expectedServiceResponse[i].Name, val["name"])
							}

							if int(val["total_page"].(float64)) != expectedServiceResponse[i].TotalPage {
								t.Errorf("expected %d as total_page but got %d", expectedServiceResponse[i].TotalPage, val["total_page"])
							}

							if val["author_id"] != expectedServiceResponse[i].AuthorId {
								t.Errorf("expected %s as author_id but got %s", expectedServiceResponse[i].AuthorId, val["author_id"])
							}

							if val["photo_url"] != expectedServiceResponse[i].PhotoURL {
								t.Errorf("expected %s as photo_url but got %s", expectedServiceResponse[i].PhotoURL, val["photo_url"])
							}

							if val["status"] != expectedServiceResponse[i].Status {
								t.Errorf("expected %s as status but got %s", expectedServiceResponse[i].Status, val["status"])
							}

							if val["completed_date"] != expectedServiceResponse[i].CompletedDate {
								t.Errorf("expected %s as completed_date but got %s", expectedServiceResponse[i].CompletedDate, val["completed_date"])
							}
						} else {
							t.Error("val should be true but got false")
						}
					}
				} else {
					t.Error("dataList should be true but got false")
				}

				// Check actual queries params that has been parsed in service
				if !reflect.DeepEqual(mockService.GetAllMockQuery, expectedQueries) {
					t.Errorf("expected %+v as query params but got %+v", expectedQueries, mockService.GetAllMockQuery)
				}
			})
		}
	})

	t.Run("get books by invalid query parameter", func(t *testing.T) {
		cases := []struct {
			Name               string
			InvalidQueryParams web.QueryParamsGetBooks
			QueryString        string
			ErrField           string
			ErrMessage         string
		}{
			{
				Name: "minimum 'name' length",
				InvalidQueryParams: web.QueryParamsGetBooks{
					Name: "Ab",
				},
				QueryString: "/api/v1/books?name=Ab",
				ErrField:    "name",
				ErrMessage:  "name must be at least 3 characters",
			},
			{
				Name: "maximum 'name' length",
				InvalidQueryParams: web.QueryParamsGetBooks{
					Name: "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
				},
				QueryString: "/api/v1/books?name=Di%20tengah%20derasnya%20arus%20teknologi%20modern%20kemampuan%20manusia%20untuk%20beradaptasi%20berpikir%20kritis%20dan%20berinovasi%20menjadi%20penentu%20utama%20dalam%20menghadapi%20tantangan%20global%20yang%20terus%20berkembang%20tanpa%20henti%20di%20segala%20bidang%20kehidupan%20manusia%20saat%20ini%20terutama%20dalam%20bidang%20teknologi.",
				ErrField:    "name",
				ErrMessage:  "name must be at most 255 characters",
			},
			{
				Name: "invalid 'status' enum",
				InvalidQueryParams: web.QueryParamsGetBooks{
					Status: "InvalidStatus",
				},
				QueryString: "/api/v1/books?status=InvalidStatus",
				ErrField:    "status",
				ErrMessage:  "the valid value for this field are only 'completed', 'reading', and 'plan_to_read'",
			},
			{
				Name: "minimum 'author_name' length",
				InvalidQueryParams: web.QueryParamsGetBooks{
					AuthorName: "Ab",
				},
				QueryString: "/api/v1/books?author_name=Ab",
				ErrField:    "author_name",
				ErrMessage:  "author_name must be at least 3 characters",
			},
			{
				Name: "maximum 'author_name' length",
				InvalidQueryParams: web.QueryParamsGetBooks{
					AuthorName: "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
				},
				QueryString: "/api/v1/books?author_name=Di%20tengah%20derasnya%20arus%20teknologi%20modern%20kemampuan%20manusia%20untuk%20beradaptasi%20berpikir%20kritis%20dan%20berinovasi%20menjadi%20penentu%20utama%20dalam%20menghadapi%20tantangan%20global%20yang%20terus%20berkembang%20tanpa%20henti%20di%20segala%20bidang%20kehidupan%20manusia%20saat%20ini%20terutama%20dalam%20bidang%20teknologi.",
				ErrField:    "author_name",
				ErrMessage:  "author_name must be at most 255 characters",
			},
			{
				Name: "invalid 'author_name' format",
				InvalidQueryParams: web.QueryParamsGetBooks{
					AuthorName: "HelloAnonymous!@",
				},
				QueryString: "/api/v1/books?author_name=HelloAnonymous!@",
				ErrField:    "author_name",
				ErrMessage:  "author_name must not contain numbers or symbols",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				queries := c.InvalidQueryParams
				expectedServiceError := validate.Struct(queries)

				mockService := &MockBookService{
					MockError: expectedServiceError,
				}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodGet, c.QueryString, nil)
				res := httptest.NewRecorder()

				handler.GetAll(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				errorList, ok := actualResponseBody.Errors.([]interface{})
				if ok {
					val, ok := errorList[0].(map[string]interface{})
					if ok {
						if val["field"] != c.ErrField {
							t.Errorf("expected error field is %s but got %s", c.ErrField, val["field"])
						}

						if val["message"] != c.ErrMessage {
							t.Errorf("expected error message is %s but got %s", c.ErrMessage, val["message"])
						}
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
				}
			})
		}
	})
}

func TestBookGetByIdHandler(t *testing.T) {
	t.Run("get book by id", func(t *testing.T) {
		pathValue := web.PathParamsGetBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		expectedServiceResponse := web.GetBookResponse{
			Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoURL:      "http://127.0.0.1:9000/book-images/ac0a9b20-2e77-4905-a665-3006763d1935.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=myAccessKey%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T235228Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=2bd452b31372e87129987c9d8e24b4ae556bde8b983db61d3c6b7fe98dba02a7",
			Status:        "completed",
			CompletedDate: "2025-09-29",
		}

		mockService := &MockBookService{
			GetByIdMockResponse: expectedServiceResponse,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.GetById(res, req)

		// Check status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code of %d but got %d", http.StatusOK, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebSuccessResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body message
		if actualResponseBody.Message != "Success get book" {
			t.Errorf("expected %s as response message but got %s", "Success get book", actualResponseBody.Message)
		}

		// Check response body data
		val, ok := actualResponseBody.Data.(map[string]interface{})
		if ok {
			if val["id"] != expectedServiceResponse.Id {
				t.Errorf("expected %s as id but got %s", expectedServiceResponse.Id, val["id"])
			}

			if val["name"] != expectedServiceResponse.Name {
				t.Errorf("expected %s as name but got %s", expectedServiceResponse.Name, val["name"])
			}

			if int(val["total_page"].(float64)) != expectedServiceResponse.TotalPage {
				t.Errorf("expected %d as total_page but got %d", expectedServiceResponse.TotalPage, val["total_page"])
			}

			if val["author_id"] != expectedServiceResponse.AuthorId {
				t.Errorf("expected %s as author_id but got %s", expectedServiceResponse.AuthorId, val["author_id"])
			}

			if val["photo_url"] != expectedServiceResponse.PhotoURL {
				t.Errorf("expected %s as photo_url but got %s", expectedServiceResponse.PhotoURL, val["photo_url"])
			}

			if val["status"] != expectedServiceResponse.Status {
				t.Errorf("expected %s as status but got %s", expectedServiceResponse.Status, val["status"])
			}

			if val["completed_date"] != expectedServiceResponse.CompletedDate {
				t.Errorf("expected %s as completed_date but got %s", expectedServiceResponse.CompletedDate, val["completed_date"])
			}
		} else {
			t.Error("val should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.GetByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdMockPathValue)
		}
	})

	t.Run("get book by id with invalid uuid", func(t *testing.T) {
		invalidUUID := "InvalidUUID"

		pathValue := web.PathParamsGetBook{
			Id: invalidUUID,
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(pathValue)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/books/InvalidUUID", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", invalidUUID)

		handler.GetById(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "id" {
					t.Errorf("expected %s as field name but got %s", "id", val["field"])
				}

				if val["message"] != fmt.Sprintf("'%s' is not a valid UUID", invalidUUID) {
					t.Errorf("expected %s as message but got %s", fmt.Sprintf("'%s' is not a valid UUID", invalidUUID), val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.GetByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdMockPathValue)
		}
	})

	t.Run("get book with not found id", func(t *testing.T) {
		pathValue := web.PathParamsGetBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusNotFound,
			[]appError.ErrAggregate{
				{
					Field:   "id",
					Message: "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found",
				},
			},
			fmt.Errorf("author with id '%s' is not found", "43723811-c8e3-4cba-85cc-142954064ae4"),
		)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.GetById(res, req)

		// Check status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code of %d but got %d", http.StatusNotFound, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "id" {
					t.Errorf("expected %s as field name but got %s", "id", val["field"])
				}

				if val["message"] != "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found" {
					t.Errorf("expected %s as message but got %s", "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.GetByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdMockPathValue)
		}
	})
}

func TestBookUpdateByIdHandler(t *testing.T) {
	t.Run("update book by id with complete data", func(t *testing.T) {
		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		bookRequest := web.UpdateBookRequest{
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}
		expectedServiceResponse := web.UpdateBookResponse{
			Id:            "43723811-c8e3-4cba-85cc-142954064ae4",
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}

		mockService := &MockBookService{
			UpdateByIdMockResponse: expectedServiceResponse,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code of %d but got %d", http.StatusOK, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebSuccessResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body message
		if actualResponseBody.Message != "Book updated successfully" {
			t.Errorf("expected %s as response message but got %s", "Book updated successfully", actualResponseBody.Message)
		}

		// Check response body data
		val, ok := actualResponseBody.Data.(map[string]interface{})
		if ok {
			if val["id"] != expectedServiceResponse.Id {
				t.Errorf("expected %s as id but got %s", expectedServiceResponse.Id, val["id"])
			}

			if val["name"] != expectedServiceResponse.Name {
				t.Errorf("expected %s as name but got %s", expectedServiceResponse.Name, val["name"])
			}

			if val["author_id"] != expectedServiceResponse.AuthorId {
				t.Errorf("expected %s as author_id but got %s", expectedServiceResponse.AuthorId, val["author_id"])
			}

			if val["photo_key"] != expectedServiceResponse.PhotoKey {
				t.Errorf("expected %s as photo_key but got %s", expectedServiceResponse.PhotoKey, val["photo_key"])
			}

			if val["status"] != expectedServiceResponse.Status {
				t.Errorf("expected %s as status but got %s", expectedServiceResponse.Status, val["status"])
			}

			if val["completed_date"] != expectedServiceResponse.CompletedDate {
				t.Errorf("expected %s as completed_date but got %s", expectedServiceResponse.CompletedDate, val["completed_date"])
			}
		} else {
			t.Error("val should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockRequest, bookRequest) {
			t.Errorf("expected %+v as request body but got %+v", bookRequest, mockService.UpdateByIdMockRequest)
		}
	})

	t.Run("update book with empty request body", func(t *testing.T) {
		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(web.UpdateBookRequest{})

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(web.UpdateBookRequest{}))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "name" {
					t.Errorf("expected %s as field name but got %s", "name", val["field"])
				}

				if val["message"] != "name is required" {
					t.Errorf("expected %s as message but got %s", "name is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			val, ok = errorList[1].(map[string]interface{})
			if ok {
				if val["field"] != "total_page" {
					t.Errorf("expected %s as field name but got %s", "total_page", val["field"])
				}

				if val["message"] != "total_page is required" {
					t.Errorf("expected %s as message but got %s", "total_page is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			val, ok = errorList[2].(map[string]interface{})
			if ok {
				if val["field"] != "author_id" {
					t.Errorf("expected %s as field name but got %s", "author_id", val["field"])
				}

				if val["message"] != "author_id is required" {
					t.Errorf("expected %s as message but got %s", "author_id is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			val, ok = errorList[3].(map[string]interface{})
			if ok {
				if val["field"] != "photo_key" {
					t.Errorf("expected %s as field name but got %s", "photo_key", val["field"])
				}

				if val["message"] != "photo_key is required" {
					t.Errorf("expected %s as message but got %s", "photo_key is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			val, ok = errorList[4].(map[string]interface{})
			if ok {
				if val["field"] != "status" {
					t.Errorf("expected %s as field name but got %s", "status", val["field"])
				}

				if val["message"] != "status is required" {
					t.Errorf("expected %s as message but got %s", "status is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}
	})

	t.Run("update book with invalid JSON payload", func(t *testing.T) {
		cases := []struct {
			Name               string
			InvalidJSONPayload string
			ErrMessage         string
		}{
			{
				Name:               "invalid type name field",
				InvalidJSONPayload: `{"name": 1, "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: name",
			},
			{
				Name:               "invalid type total_page field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": "379", "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: total_page",
			},
			{
				Name:               "invalid type author_id field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": 123, "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: author_id",
			},
			{
				Name:               "invalid type photo_key field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": 123, "status": "completed", "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: photo_key",
			},
			{
				Name:               "invalid type status field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": 1, "completed_date": "2025-10-29"}`,
				ErrMessage:         "Invalid JSON type for field: status",
			},
			{
				Name:               "invalid type completed_date field",
				InvalidJSONPayload: `{"name": "Laut Bercerita", "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": 2025}`,
				ErrMessage:         "Invalid JSON type for field: completed_date",
			},
			{
				Name:               "invalid JSON payload",
				InvalidJSONPayload: `{"name": "Laut Bercerita" "total_page": 379, "author_id": "c512ae16-5f33-4a3c-a1e1-977bd5a20af3", "photo_key": "ac0a9b20-2e77-4905-a665-3006763d1935.jpg", "status": "completed", "completed_date": 2025-10-29}`,
				ErrMessage:         "Invalid JSON payload",
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				mockService := &MockBookService{}

				handler := NewBookHandler(mockService)

				req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", strings.NewReader(c.InvalidJSONPayload))
				res := httptest.NewRecorder()

				req.Header.Set("Content-Type", "application/json")
				// Path value must be set since httptest.NewRequest never goes through http.ServeMux
				req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

				handler.UpdateById(res, req)

				// Check status code
				if res.Code != http.StatusBadRequest {
					t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
				}

				// Get the actual response
				var actualResponseBody web.WebFailedResponse
				err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
				if err != nil {
					t.Fatalf("error when parsing res body: %v", err)
				}

				val, ok := actualResponseBody.Errors.(string)
				if ok {
					if val != c.ErrMessage {
						t.Errorf("expected %s but got %s", c.ErrMessage, val)
					}
				} else {
					t.Error("val should be true but got false")
				}
			})
		}
	})

	t.Run("update book with not found id", func(t *testing.T) {
		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		bookRequest := web.UpdateBookRequest{
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusNotFound,
			[]appError.ErrAggregate{
				{
					Field:   "id",
					Message: "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found",
				},
			},
			fmt.Errorf("book with id '%s' is not found", "43723811-c8e3-4cba-85cc-142954064ae4"),
		)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code of %d but got %d", http.StatusNotFound, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "id" {
					t.Errorf("expected %s as field name but got %s", "id", val["field"])
				}

				if val["message"] != "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found" {
					t.Errorf("expected %s as message but got %s", "book with id '43723811-c8e3-4cba-85cc-142954064ae4' is not found", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockRequest, bookRequest) {
			t.Errorf("expected %+v as request body but got %+v", bookRequest, mockService.UpdateByIdMockRequest)
		}
	})

	t.Run("update book with invalid id", func(t *testing.T) {
		invalidUUID := "InvalidUUID"

		pathValue := web.PathParamsUpdateBook{
			Id: invalidUUID,
		}
		bookRequest := web.UpdateBookRequest{
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(pathValue)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/InvalidUUID", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", invalidUUID)

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "id" {
					t.Errorf("expected %s as field name but got %s", "id", val["field"])
				}

				if val["message"] != fmt.Sprintf("'%s' is not a valid UUID", invalidUUID) {
					t.Errorf("expected %s as message but got %s", fmt.Sprintf("'%s' is not a valid UUID", invalidUUID), val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}
	})

	t.Run("update book with not found author_id", func(t *testing.T) {
		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		bookRequest := web.UpdateBookRequest{
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusNotFound,
			[]appError.ErrAggregate{
				{
					Field:   "author_id",
					Message: "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found",
				},
			},
			fmt.Errorf("author with id '%s' is not found", "c512ae16-5f33-4a3c-a1e1-977bd5a20af3"),
		)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code of %d but got %d", http.StatusNotFound, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "author_id" {
					t.Errorf("expected %s as field name but got %s", "author_id", val["field"])
				}

				if val["message"] != "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found" {
					t.Errorf("expected %s as message but got %s", "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockRequest, bookRequest) {
			t.Errorf("expected %+v as request body but got %+v", bookRequest, mockService.UpdateByIdMockRequest)
		}
	})

	t.Run("update book with invalid author id", func(t *testing.T) {
		invalidAuthorId := "InvalidUUID"

		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		bookRequest := web.UpdateBookRequest{
			Name:          "New Book Name",
			TotalPage:     100,
			AuthorId:      invalidAuthorId,
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1935.jpg",
			Status:        "plan_to_read",
			CompletedDate: "0000-00-00",
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(bookRequest)

		mockService := &MockBookService{
			MockError: expectedServiceError,
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "43723811-c8e3-4cba-85cc-142954064ae4")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		// Check response body data
		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "author_id" {
					t.Errorf("expected %s as field name but got %s", "author_id", val["field"])
				}

				if val["message"] != fmt.Sprintf("'%s' is not a valid UUID", invalidAuthorId) {
					t.Errorf("expected %s as message but got %s", fmt.Sprintf("'%s' is not a valid UUID", invalidAuthorId), val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}
	})

	t.Run("update book with existing name with the same author", func(t *testing.T) {
		pathValue := web.PathParamsUpdateBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}
		bookRequest := web.CreateBookRequest{
			Name:          "Laut Bercerita",
			TotalPage:     379,
			AuthorId:      "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			PhotoKey:      "ac0a9b20-2e77-4905-a665-3006763d1934.jpg",
			Status:        "reading",
			CompletedDate: "0000-00-00",
		}
		expectedServiceError := []appError.ErrAggregate{
			appError.ErrAggregate{
				Field:   "name",
				Message: "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists",
			},
		}

		mockService := &MockBookService{
			MockError: appError.NewAppError(
				http.StatusBadRequest,
				expectedServiceError,
				nil,
			),
		}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", ToJSON(bookRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", pathValue.Id)

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code of %d but got %d", http.StatusBadRequest, res.Code)
		}

		// Get the actual response
		var actualResponseBody web.WebFailedResponse
		err := json.NewDecoder(res.Body).Decode(&actualResponseBody)
		if err != nil {
			t.Fatalf("error when parsing res body: %v", err)
		}

		errorList, ok := actualResponseBody.Errors.([]interface{})
		if ok {
			val, ok := errorList[0].(map[string]interface{})
			if ok {
				if val["field"] != "name" {
					t.Errorf("expected error field is %s but got %s", "name", val["field"])
				}

				if val["message"] != "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists" {
					t.Errorf("expected error message is %s but got %s", "Laut Bercerita with author id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is already exists", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdMockPathValue)
		}
	})
}

func TestBookDeleteByIdHandler(t *testing.T) {
	t.Run("delete book with id", func(t *testing.T) {
		pathValue := web.PathParamsDeleteBook{
			Id: "43723811-c8e3-4cba-85cc-142954064ae4",
		}

		mockService := &MockBookService{}

		handler := NewBookHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/books/43723811-c8e3-4cba-85cc-142954064ae4", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", pathValue.Id)

		handler.DeleteById(res, req)

		// Check status code
		if res.Code != http.StatusNoContent {
			t.Errorf("expected status code %d but got %d", http.StatusNoContent, res.Code)
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.DeleteByIdMockPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.DeleteByIdMockPathValue)
		}
	})
}
