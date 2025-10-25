package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/mhaatha/go-bookshelf/internal/config"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockAuthorService struct {
	CreateCalledWithRequest web.CreateAuthorRequest
	MockCreateResponse      web.CreateAuthorResponse
	MockError               error
}

func (m *MockAuthorService) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	m.CreateCalledWithRequest = request

	if m.MockError != nil {
		return m.MockCreateResponse, m.MockError
	}

	return m.MockCreateResponse, nil
}

func (m *MockAuthorService) GetAllAuthors(ctx context.Context, queris web.QueryParamsGetAuthors) ([]web.GetAuthorResponse, error) {
	return nil, nil
}

func (m *MockAuthorService) GetAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor) (web.GetAuthorResponse, error) {
	return web.GetAuthorResponse{}, nil
}

func (m *MockAuthorService) UpdateAuthorById(ctx context.Context, pathValues web.PathParamsUpdateAuthor, request web.UpdateAuthorRequest) (web.UpdateAuthorResponse, error) {
	return web.UpdateAuthorResponse{}, nil
}

func (m *MockAuthorService) DeleteAuthorById(ctx context.Context, pathValues web.PathParamsDeleteAuthor) error {
	return nil
}

func TestAuthorCreateHandler(t *testing.T) {
	t.Run("create author with complete data", func(t *testing.T) {
		authorRequest := web.CreateAuthorRequest{
			FullName:    "Leila S. Chudori",
			Nationality: "Indonesia",
		}
		expectedServiceResponse := web.CreateAuthorResponse{
			Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			FullName:    "Leila S. Chudori",
			Nationality: "Indonesia",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockService := &MockAuthorService{
			MockCreateResponse: expectedServiceResponse,
			MockError:          nil,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/authors", toJSON(authorRequest))
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
		if actualResponseBody.Message != "Author created successfully" {
			t.Errorf("expected %s as response message but got %s", "Author created successfully", actualResponseBody.Message)
		}

		// Check response body data
		val, ok := actualResponseBody.Data.(map[string]interface{})
		if ok {
			if val["id"] != expectedServiceResponse.Id {
				t.Errorf("expected %s as id but got %s", expectedServiceResponse.Id, val["id"])
			}

			if val["full_name"] != expectedServiceResponse.FullName {
				t.Errorf("expected %s as full_name but got %s", expectedServiceResponse.FullName, val["full_name"])
			}

			if val["nationality"] != expectedServiceResponse.Nationality {
				t.Errorf("expected %s as nationality but got %s", expectedServiceResponse.Nationality, val["nationality"])
			}
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.CreateCalledWithRequest, authorRequest) {
			t.Errorf("expected %+v as request body but got %+v", authorRequest, mockService.CreateCalledWithRequest)
		}
	})

	t.Run("create author with invalid full_name", func(t *testing.T) {
		cases := []struct {
			Name          string
			AuthorRequest web.CreateAuthorRequest
			ErrField      string
			ErrMessage    string
		}{
			{
				Name: "minimum length",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Hi",
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at least 3 characters",
			},
			{
				Name: "maximum length",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at most 255 characters",
			},
			{
				Name: "required",
				AuthorRequest: web.CreateAuthorRequest{
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name is required",
			},
			{
				Name: "valid full_name",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Invalid Full Name #123",
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must not contain numbers or symbols",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				authorRequest := c.AuthorRequest
				expectedServiceError := validate.Struct(authorRequest)

				mockService := &MockAuthorService{
					MockError: expectedServiceError,
				}

				handler := NewAuthorHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/authors", toJSON(authorRequest))
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
					}
				}
			})
		}
	})

	t.Run("create author with invalid nationality", func(t *testing.T) {
		cases := []struct {
			Name          string
			AuthorRequest web.CreateAuthorRequest
			ErrField      string
			ErrMessage    string
		}{
			{
				Name: "minimum length",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Leila S. Chudori",
					Nationality: "In",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must be at least 3 characters",
			},
			{
				Name: "maximum length",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Leila S. Chudori",
					Nationality: "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must be at most 255 characters",
			},
			{
				Name: "required",
				AuthorRequest: web.CreateAuthorRequest{
					FullName: "Leila S. Chudori",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality is required",
			},
			{
				Name: "alpha only",
				AuthorRequest: web.CreateAuthorRequest{
					FullName:    "Leila S. Chudori",
					Nationality: "Invalid Nationality Name #123",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must not contain numbers or symbols",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				authorRequest := c.AuthorRequest
				expectedServiceError := validate.Struct(authorRequest)

				mockService := &MockAuthorService{
					MockError: expectedServiceError,
				}

				handler := NewAuthorHandler(mockService)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/authors", toJSON(authorRequest))
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
					}
				}
			})
		}
	})

	t.Run("create author with existing full_name", func(t *testing.T) {
		authorRequest := web.CreateAuthorRequest{
			FullName:    "Leila S. Chudori",
			Nationality: "Indonesia",
		}
		expectedServiceError := []appError.ErrAggregate{
			{
				Field:   "full_name",
				Message: "author Leila S. Chudori is already exists",
			},
		}

		mockService := &MockAuthorService{
			MockError: appError.NewAppError(
				http.StatusBadRequest,
				expectedServiceError,
				nil,
			),
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/authors", toJSON(authorRequest))
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
				if val["field"] != "full_name" {
					t.Errorf("expected error field is %s but got %s", "full_name", val["field"])
				}

				if val["message"] != "author Leila S. Chudori is already exists" {
					t.Errorf("expected error message is %s but got %s", "author Leila S. Chudori is already exists", val["message"])
				}
			}
		}
	})
}

func toJSON(data interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(jsonBytes)
}
