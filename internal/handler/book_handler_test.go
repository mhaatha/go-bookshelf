package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mhaatha/go-bookshelf/internal/config"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockBookService struct {
	// CreateNewBook
	CreateMockRequest  web.CreateBookRequest
	CreateMockResponse web.CreateBookResponse

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
	return nil, nil
}

func (m *MockBookService) GetBookById(ctx context.Context, pathValues web.PathParamsGetBook) (web.GetBookResponse, error) {
	return web.GetBookResponse{}, nil
}

func (m *MockBookService) UpdateBookById(ctx context.Context, pathValues web.PathParamsUpdateBook, request web.UpdateBookRequest) (web.UpdateBookResponse, error) {
	return web.UpdateBookResponse{}, nil
}

func (m *MockBookService) DeleteBookById(ctx context.Context, pathValues web.PathParamsDeleteBook) error {
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
}
