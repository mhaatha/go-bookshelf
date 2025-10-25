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

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockAuthorService struct {
	CreateCalledWithRequest web.CreateAuthorRequest
	MockCreateResponse      web.CreateAuthorResponse
	MockError               error
}

func (m *MockAuthorService) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	m.CreateCalledWithRequest = request

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
}

func toJSON(data interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(jsonBytes)
}
