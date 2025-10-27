package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mhaatha/go-bookshelf/internal/config"
	appError "github.com/mhaatha/go-bookshelf/internal/errors"
	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockAuthorService struct {
	// CreateNewAuthor
	CreateCalledWithRequest web.CreateAuthorRequest
	MockCreateResponse      web.CreateAuthorResponse

	// GetAllAuthors
	GetAllCalledWithQuery web.QueryParamsGetAuthors
	MockGetAllResponse    []web.GetAuthorResponse

	// GetAuthorById
	GetByIdCalledWithPathValue web.PathParamsGetAuthor
	MockGetByIdResponse        web.GetAuthorResponse

	// UpdateAuthorById
	UpdateByIdCalledWithRequest   web.UpdateAuthorRequest
	UpdateByIdCalledWithPathValue web.PathParamsUpdateAuthor
	MockUpdateByIdResponse        web.UpdateAuthorResponse

	MockError error
}

func (m *MockAuthorService) CreateNewAuthor(ctx context.Context, request web.CreateAuthorRequest) (web.CreateAuthorResponse, error) {
	m.CreateCalledWithRequest = request

	if m.MockError != nil {
		return m.MockCreateResponse, m.MockError
	}

	return m.MockCreateResponse, nil
}

func (m *MockAuthorService) GetAllAuthors(ctx context.Context, queris web.QueryParamsGetAuthors) ([]web.GetAuthorResponse, error) {
	m.GetAllCalledWithQuery = queris

	if m.MockError != nil {
		return m.MockGetAllResponse, m.MockError
	}

	return m.MockGetAllResponse, nil
}

func (m *MockAuthorService) GetAuthorById(ctx context.Context, pathValues web.PathParamsGetAuthor) (web.GetAuthorResponse, error) {
	m.GetByIdCalledWithPathValue = pathValues

	if m.MockError != nil {
		return m.MockGetByIdResponse, m.MockError
	}

	return m.MockGetByIdResponse, nil
}

func (m *MockAuthorService) UpdateAuthorById(ctx context.Context, pathValues web.PathParamsUpdateAuthor, request web.UpdateAuthorRequest) (web.UpdateAuthorResponse, error) {
	m.UpdateByIdCalledWithPathValue = pathValues
	m.UpdateByIdCalledWithRequest = request

	if m.MockError != nil {
		return m.MockUpdateByIdResponse, m.MockError
	}

	return m.MockUpdateByIdResponse, nil
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
		}

		mockService := &MockAuthorService{
			MockCreateResponse: expectedServiceResponse,
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
		} else {
			t.Error("val should be true but got false")
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
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
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
					} else {
						t.Error("val should be true but got false")
					}
				} else {
					t.Error("errorList should be true but got false")
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
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}
	})

	t.Run("create author with invalid JSON payload", func(t *testing.T) {
		invalidJSONPayload := `{"full_name":}`
		mockService := &MockAuthorService{}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/authors", strings.NewReader(invalidJSONPayload))
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
			if val != "Invalid JSON payload" {
				t.Errorf("expected %s but got %s", "Invalid JSON payload", val)
			}
		} else {
			t.Error("val should be true but got false")
		}
	})
}

func TestAuthorGetAllHandler(t *testing.T) {
	t.Run("get all authors", func(t *testing.T) {
		expectedServiceResponse := []web.GetAuthorResponse{
			{
				Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
				FullName:    "Leila S. Chudori",
				Nationality: "Indonesia",
			},
			{
				Id:          "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
				FullName:    "Henry Manampiring",
				Nationality: "Indonesia",
			},
		}

		mockService := &MockAuthorService{
			MockGetAllResponse: expectedServiceResponse,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors", nil)
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
		if actualResponseBody.Message != "Success get all authors" {
			t.Errorf("expected %s as response message but got %s", "Success get all authors", actualResponseBody.Message)
		}

		// Check response body data
		dataList, ok := actualResponseBody.Data.([]interface{})
		if ok {
			// First data from the JSON array
			val, ok := dataList[0].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[0].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[0].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[0].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[0].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[0].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[0].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			// Second data from the JSON array
			val, ok = dataList[1].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[1].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[1].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[1].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[1].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[1].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[1].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("dataList should be true but got false")
		}
	})

	t.Run("get authors by full_name query parameter", func(t *testing.T) {
		expectedQueries := web.QueryParamsGetAuthors{
			FullName: "Leila",
		}
		expectedServiceResponse := []web.GetAuthorResponse{
			{
				Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
				FullName:    "Leila S. Chudori",
				Nationality: "Indonesia",
			},
		}

		mockService := &MockAuthorService{
			MockGetAllResponse:    expectedServiceResponse,
			GetAllCalledWithQuery: expectedQueries,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors?full_name=Leila", nil)
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
			// First data from the JSON array
			val, ok := dataList[0].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[0].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[0].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[0].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[0].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[0].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[0].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("dataList should be true but got false")
		}

		// Check actual queries params that has been parsed in service
		if !reflect.DeepEqual(mockService.GetAllCalledWithQuery, expectedQueries) {
			t.Errorf("expected %+v as query params but got %+v", expectedQueries, mockService.GetAllCalledWithQuery)
		}
	})

	t.Run("get authors by nationality query parameter", func(t *testing.T) {
		expectedQueries := web.QueryParamsGetAuthors{
			Nationality: "Indonesia",
		}
		expectedServiceResponse := []web.GetAuthorResponse{
			{
				Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
				FullName:    "Leila S. Chudori",
				Nationality: "Indonesia",
			},
			{
				Id:          "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
				FullName:    "Henry Manampiring",
				Nationality: "Indonesia",
			},
		}

		mockService := &MockAuthorService{
			MockGetAllResponse:    expectedServiceResponse,
			GetAllCalledWithQuery: expectedQueries,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors?nationality=Indonesia", nil)
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
			// First data from the JSON array
			val, ok := dataList[0].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[0].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[0].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[0].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[0].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[0].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[0].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			// Second data from the JSON array
			val, ok = dataList[1].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[1].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[1].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[1].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[1].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[1].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[1].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("dataList should be true but got false")
		}

		// Check actual queries params that has been parsed in service
		if !reflect.DeepEqual(mockService.GetAllCalledWithQuery, expectedQueries) {
			t.Errorf("expected %+v as query params but got %+v", expectedQueries, mockService.GetAllCalledWithQuery)
		}
	})

	t.Run("get authors by all queries parameter", func(t *testing.T) {
		expectedQueries := web.QueryParamsGetAuthors{
			FullName:    "Leila",
			Nationality: "Indonesia",
		}
		expectedServiceResponse := []web.GetAuthorResponse{
			{
				Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
				FullName:    "Leila S. Chudori",
				Nationality: "Indonesia",
			},
		}

		mockService := &MockAuthorService{
			MockGetAllResponse:    expectedServiceResponse,
			GetAllCalledWithQuery: expectedQueries,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors?full_name=Leila&nationality=Indonesia", nil)
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
			// First data from the JSON array
			val, ok := dataList[0].(map[string]interface{})
			if ok {
				if val["id"] != expectedServiceResponse[0].Id {
					t.Errorf("expected %s as id but got %s", expectedServiceResponse[0].Id, val["id"])
				}

				if val["full_name"] != expectedServiceResponse[0].FullName {
					t.Errorf("expected %s as full_name but got %s", expectedServiceResponse[0].FullName, val["full_name"])
				}

				if val["nationality"] != expectedServiceResponse[0].Nationality {
					t.Errorf("expected %s as nationality but got %s", expectedServiceResponse[0].Nationality, val["nationality"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("dataList should be true but got false")
		}

		// Check actual queries params that has been parsed in service
		if !reflect.DeepEqual(mockService.GetAllCalledWithQuery, expectedQueries) {
			t.Errorf("expected %+v as query params but got %+v", expectedQueries, mockService.GetAllCalledWithQuery)
		}
	})

	t.Run("get authors with invalid full_name query parameter", func(t *testing.T) {
		cases := []struct {
			Name       string
			Query      web.QueryParamsGetAuthors
			ErrField   string
			ErrMessage string
		}{
			{
				Name: "minimum length",
				Query: web.QueryParamsGetAuthors{
					FullName: "Hi",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at least 3 characters",
			},
			{
				Name: "maximum length",
				Query: web.QueryParamsGetAuthors{
					FullName: "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at most 255 characters",
			},
			{
				Name: "valid full_name",
				Query: web.QueryParamsGetAuthors{
					FullName: "Invalid FullName #123",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must not contain numbers or symbols",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				queries := web.QueryParamsGetAuthors{
					FullName: c.Query.FullName,
				}
				expectedServiceError := validate.Struct(queries)

				mockService := &MockAuthorService{
					MockError:             expectedServiceError,
					GetAllCalledWithQuery: queries,
				}

				handler := NewAuthorHandler(mockService)

				req := httptest.NewRequest(http.MethodGet, "/api/v1/authors", nil)
				q := req.URL.Query()
				q.Add("full_name", c.Query.FullName)
				req.URL.RawQuery = q.Encode()

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

	t.Run("get authors with invalid nationality query parameter", func(t *testing.T) {
		cases := []struct {
			Name       string
			Query      web.QueryParamsGetAuthors
			ErrField   string
			ErrMessage string
		}{
			{
				Name: "minimum length",
				Query: web.QueryParamsGetAuthors{
					Nationality: "Hi",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must be at least 3 characters",
			},
			{
				Name: "maximum length",
				Query: web.QueryParamsGetAuthors{
					Nationality: "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must be at most 255 characters",
			},
			{
				Name: "alpha only",
				Query: web.QueryParamsGetAuthors{
					Nationality: "Invalid FullName #123",
				},
				ErrField:   "nationality",
				ErrMessage: "nationality must not contain numbers or symbols",
			},
		}

		validate := config.ValidatorInit()
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				queries := web.QueryParamsGetAuthors{
					Nationality: c.Query.Nationality,
				}
				expectedServiceError := validate.Struct(queries)

				mockService := &MockAuthorService{
					MockError:             expectedServiceError,
					GetAllCalledWithQuery: queries,
				}

				handler := NewAuthorHandler(mockService)

				req := httptest.NewRequest(http.MethodGet, "/api/v1/authors", nil)
				q := req.URL.Query()
				q.Add("full_name", c.Query.Nationality)
				req.URL.RawQuery = q.Encode()

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

func TestAuthorGetByIdHandler(t *testing.T) {
	t.Run("get author by id", func(t *testing.T) {
		pathValue := web.PathParamsGetAuthor{
			Id: "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
		}
		expectedServiceResponse := web.GetAuthorResponse{
			Id:          "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
			FullName:    "Leila S. Chudori",
			Nationality: "Indonesia",
		}

		mockService := &MockAuthorService{
			GetByIdCalledWithPathValue: pathValue,
			MockGetByIdResponse:        expectedServiceResponse,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/c512ae16-5f33-4a3c-a1e1-977bd5a20af3", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "c512ae16-5f33-4a3c-a1e1-977bd5a20af3")

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
		if actualResponseBody.Message != "Success get author" {
			t.Errorf("expected %s as response message but got %s", "Success get author", actualResponseBody.Message)
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
		} else {
			t.Error("val should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.GetByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdCalledWithPathValue)
		}
	})

	t.Run("get author by id with invalid UUID", func(t *testing.T) {
		invalidUUID := "InvalidUUID"

		pathValue := web.PathParamsGetAuthor{
			Id: invalidUUID,
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(pathValue)

		mockService := &MockAuthorService{
			MockError: expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authors/%s", invalidUUID), nil)
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
		if !reflect.DeepEqual(mockService.GetByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdCalledWithPathValue)
		}
	})

	t.Run("get author with not found id", func(t *testing.T) {
		pathValue := web.PathParamsGetAuthor{
			Id: "c512ae16-5f33-4a3c-a1e1-977bd5a20af3",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusNotFound,
			[]appError.ErrAggregate{
				{
					Field:   "id",
					Message: "author with id 'c512ae16-5f33-4a3c-a1e1-977bd5a20af3' is not found",
				},
			},
			fmt.Errorf("author with id '%s' is not found", "c512ae16-5f33-4a3c-a1e1-977bd5a20af3"),
		)

		mockService := &MockAuthorService{
			GetByIdCalledWithPathValue: pathValue,
			MockError:                  expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/c512ae16-5f33-4a3c-a1e1-977bd5a20af3", nil)
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "c512ae16-5f33-4a3c-a1e1-977bd5a20af3")

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
		if !reflect.DeepEqual(mockService.GetByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.GetByIdCalledWithPathValue)
		}
	})
}

func TestAuthorUpdateHandler(t *testing.T) {
	t.Run("update author with complete data", func(t *testing.T) {
		pathValue := web.PathParamsUpdateAuthor{
			Id: "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
		}
		authorRequest := web.UpdateAuthorRequest{
			FullName:    "Henry Manampiring",
			Nationality: "Indonesia",
		}
		expectedServiceResponse := web.UpdateAuthorResponse{
			Id:          "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
			FullName:    "Henry Manampiring",
			Nationality: "Indonesia",
			UpdatedAt:   time.Date(2025, 10, 27, 8, 21, 0, 0, time.UTC),
		}

		mockService := &MockAuthorService{
			UpdateByIdCalledWithPathValue: pathValue,
			UpdateByIdCalledWithRequest:   authorRequest,
			MockUpdateByIdResponse:        expectedServiceResponse,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", toJSON(authorRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

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
		if actualResponseBody.Message != "Author updated successfully" {
			t.Errorf("expected %s as response message but got %s", "Author updated successfully", actualResponseBody.Message)
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

			if val["updated_at"] != expectedServiceResponse.UpdatedAt.Format(time.RFC3339) {
				t.Errorf("expected %s as updated_at but got %s", expectedServiceResponse.UpdatedAt.Format(time.RFC3339), val["updated_at"])
			}
		} else {
			t.Error("val should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdCalledWithPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithRequest, authorRequest) {
			t.Errorf("expected %+v as request body but got %+v", authorRequest, mockService.UpdateByIdCalledWithRequest)
		}
	})

	t.Run("update author with existing full_name", func(t *testing.T) {
		pathValue := web.PathParamsUpdateAuthor{
			Id: "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
		}
		authorRequest := web.UpdateAuthorRequest{
			FullName:    "Leila S. Chudori",
			Nationality: "Indonesia",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusConflict,
			[]appError.ErrAggregate{
				{
					Field:   "full_name",
					Message: "author with full_name 'Leila S. Chudori' already exists",
				},
			},
			fmt.Errorf("author with full_name '%s' already exists", "Leila S. Chudori"),
		)

		mockService := &MockAuthorService{
			UpdateByIdCalledWithPathValue: pathValue,
			UpdateByIdCalledWithRequest:   authorRequest,
			MockError:                     expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", toJSON(authorRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

		handler.UpdateById(res, req)

		// Check status code
		if res.Code != http.StatusConflict {
			t.Errorf("expected status code of %d but got %d", http.StatusConflict, res.Code)
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
				if val["field"] != "full_name" {
					t.Errorf("expected %s as field name but got %s", "full_name", val["field"])
				}

				if val["message"] != "author with full_name 'Leila S. Chudori' already exists" {
					t.Errorf("expected %s as message but got %s", "author with full_name 'Leila S. Chudori' already exists", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdCalledWithPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithRequest, authorRequest) {
			t.Errorf("expected %+v as request body but got %+v", authorRequest, mockService.UpdateByIdCalledWithRequest)
		}
	})

	t.Run("update author with empty request body", func(t *testing.T) {
		pathValue := web.PathParamsUpdateAuthor{
			Id: "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(web.UpdateAuthorRequest{})

		mockService := &MockAuthorService{
			UpdateByIdCalledWithPathValue: pathValue,
			MockError:                     expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", toJSON(web.UpdateAuthorRequest{}))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

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
				if val["field"] != "full_name" {
					t.Errorf("expected %s as field name but got %s", "full_name", val["field"])
				}

				if val["message"] != "full_name is required" {
					t.Errorf("expected %s as message but got %s", "full_name is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}

			val, ok = errorList[1].(map[string]interface{})
			if ok {
				if val["field"] != "nationality" {
					t.Errorf("expected %s as field name but got %s", "nationality", val["field"])
				}

				if val["message"] != "nationality is required" {
					t.Errorf("expected %s as message but got %s", "nationality is required", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdCalledWithPathValue)
		}
	})

	t.Run("update author with invalid JSON payload", func(t *testing.T) {
		invalidJSONPayload := `{"full_name":}`
		mockService := &MockAuthorService{}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", strings.NewReader(invalidJSONPayload))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

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
		val, ok := actualResponseBody.Errors.(string)
		if ok {
			if val != "Invalid JSON payload" {
				t.Errorf("expected %s as message but got %s", "Invalid JSON payload", val)
			}
		} else {
			t.Error("val should be true but got false")
		}
	})

	t.Run("update author with not found id", func(t *testing.T) {
		pathValue := web.PathParamsUpdateAuthor{
			Id: "84a069f3-2620-4da4-8bb5-5c39bbe7cda7",
		}
		authorRequest := web.UpdateAuthorRequest{
			FullName:    "Henry Manampiring",
			Nationality: "Indonesian",
		}
		expectedServiceError := appError.NewAppError(
			http.StatusNotFound,
			[]appError.ErrAggregate{
				{
					Field:   "id",
					Message: "author with id '84a069f3-2620-4da4-8bb5-5c39bbe7cda7' is not found",
				},
			},
			fmt.Errorf("author with id '%s' is not found", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7"),
		)

		mockService := &MockAuthorService{
			UpdateByIdCalledWithPathValue: pathValue,
			UpdateByIdCalledWithRequest:   authorRequest,
			MockError:                     expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", toJSON(authorRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

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

				if val["message"] != "author with id '84a069f3-2620-4da4-8bb5-5c39bbe7cda7' is not found" {
					t.Errorf("expected %s as message but got %s", "author with id '84a069f3-2620-4da4-8bb5-5c39bbe7cda7' is not found", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}

		// Check actual path values that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithPathValue, pathValue) {
			t.Errorf("expected %+v as path value but got %+v", pathValue, mockService.UpdateByIdCalledWithPathValue)
		}

		// Check actual request body that has been parsed in service
		if !reflect.DeepEqual(mockService.UpdateByIdCalledWithRequest, authorRequest) {
			t.Errorf("expected %+v as request body but got %+v", authorRequest, mockService.UpdateByIdCalledWithRequest)
		}
	})

	t.Run("update author with invalid UUID", func(t *testing.T) {
		pathValue := web.PathParamsUpdateAuthor{
			Id: "InvalidUUID",
		}
		authorRequest := web.UpdateAuthorRequest{
			FullName:    "Henry Manampiring",
			Nationality: "Indonesia",
		}
		validate := config.ValidatorInit()
		expectedServiceError := validate.Struct(pathValue)

		mockService := &MockAuthorService{
			MockError: expectedServiceError,
		}

		handler := NewAuthorHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/InvalidUUID", toJSON(authorRequest))
		res := httptest.NewRecorder()

		// Path value must be set since httptest.NewRequest never goes through http.ServeMux
		req.SetPathValue("id", "InvalidUUID")

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

				if val["message"] != "'InvalidUUID' is not a valid UUID" {
					t.Errorf("expected %s as message but got %s", "'InvalidUUID' is not a valid UUID", val["message"])
				}
			} else {
				t.Error("val should be true but got false")
			}
		} else {
			t.Error("errorList should be true but got false")
		}
	})

	t.Run("update author with invalid full_name", func(t *testing.T) {
		cases := []struct {
			Name          string
			AuthorRequest web.UpdateAuthorRequest
			ErrField      string
			ErrMessage    string
		}{
			{
				Name: "minimum length",
				AuthorRequest: web.UpdateAuthorRequest{
					FullName:    "Hi",
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at least 3 characters",
			},
			{
				Name: "maximum length",
				AuthorRequest: web.UpdateAuthorRequest{
					FullName:    "Di tengah derasnya arus teknologi modern kemampuan manusia untuk beradaptasi berpikir kritis dan berinovasi menjadi penentu utama dalam menghadapi tantangan global yang terus berkembang tanpa henti di segala bidang kehidupan manusia saat ini terutama dalam bidang teknologi.",
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name must be at most 255 characters",
			},
			{
				Name: "required",
				AuthorRequest: web.UpdateAuthorRequest{
					Nationality: "Indonesia",
				},
				ErrField:   "full_name",
				ErrMessage: "full_name is required",
			},
			{
				Name: "valid full_name",
				AuthorRequest: web.UpdateAuthorRequest{
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
				authorRequest := web.UpdateAuthorRequest{
					FullName:    c.AuthorRequest.FullName,
					Nationality: c.AuthorRequest.Nationality,
				}
				expectedServiceError := validate.Struct(authorRequest)

				mockService := &MockAuthorService{
					MockError: expectedServiceError,
				}

				handler := NewAuthorHandler(mockService)

				req := httptest.NewRequest(http.MethodPut, "/api/v1/authors/84a069f3-2620-4da4-8bb5-5c39bbe7cda7", toJSON(authorRequest))
				res := httptest.NewRecorder()

				// Path value must be set since httptest.NewRequest never goes through http.ServeMux
				req.SetPathValue("id", "84a069f3-2620-4da4-8bb5-5c39bbe7cda7")

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

// Helper functions
func toJSON(data interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(jsonBytes)
}
