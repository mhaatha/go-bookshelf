package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mhaatha/go-bookshelf/internal/model/web"
)

type MockUploadService struct {
	// GetBookPresignedURL
	GetPresignedURLMockResponse web.GetBookPresignedURLResponse

	MockError error
}

func (m *MockUploadService) GetBookPresignedURL(ctx context.Context) (web.GetBookPresignedURLResponse, error) {
	if m.MockError != nil {
		return web.GetBookPresignedURLResponse{}, m.MockError
	}

	return m.GetPresignedURLMockResponse, nil
}

func TestGetBookPresignedURL(t *testing.T) {
	t.Run("get book presigned url", func(t *testing.T) {
		expectedServiceResponse := web.GetBookPresignedURLResponse{
			URL:      "http://127.0.0.1:9000/book-images/",
			FormData: `{ "Content-Type": "image/jpeg", "bucket": "book-images", "key": "35e45eae-123c-4727-8b46-e6b2ea939e12.jpg", "policy": "eyJleHBpcmF0aW9uIjoiMjAyNS0xMS0wOFQwMzoxODoyNC43NTdaIiwiY29uZGl0aW9ucyI6W1siZXEiLCIkYnVja2V0IiwiYm9vay1pbWFnZXMiXSxbImVxIiwiJGtleSIsIjM1ZTQ1ZWFlLTEyM2MtNDcyNy04YjQ2LWU2YjJlYTkzOWUxMi5qcGciXSxbImVxIiwiJENvbnRlbnQtVHlwZSIsImltYWdlL2pwZWciXSxbImVxIiwiJHgtYW16LWRhdGUiLCIyMDI1MTEwOFQwMzEzMjRaIl0sWyJlcSIsIiR4LWFtei1hbGdvcml0aG0iLCJBV1M0LUhNQUMtU0hBMjU2Il0sWyJlcSIsIiR4LWFtei1jcmVkZW50aWFsIiwibXlBY2Nlc3NLZXkvMjAyNTExMDgvdXMtZWFzdC0xL3MzL2F3czRfcmVxdWVzdCJdLFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLCAxMDI0LCA1MjQyODgwXV19", "x-amz-algorithm": "AWS4-HMAC-SHA256", "x-amz-credential": "myAccessKey/20251108/us-east-1/s3/aws4_request", "x-amz-date": "20251108T031324Z", "x-amz-signature": "3025c0c4b29f7ffac992e04f9c1d81e205f2f33f6b4dc44d1d7cc96e5c071453" }`,
		}

		mockService := &MockUploadService{
			GetPresignedURLMockResponse: expectedServiceResponse,
		}

		handler := NewUploadHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/upload/books/presigned-url", nil)
		res := httptest.NewRecorder()

		handler.GetBookPresignedURL(res, req)

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
		if actualResponseBody.Message != "Success get presigned URL" {
			t.Errorf("expected '%s' as response message but got '%s'", "Success get presigned URL", actualResponseBody.Message)
		}

		// Check response body data
		val, ok := actualResponseBody.Data.(map[string]interface{})
		if ok {
			if val["url"] != expectedServiceResponse.URL {
				t.Errorf("expected url '%s' but got '%s'", expectedServiceResponse.URL, val["url"])
			}

			formDataStr, ok := val["form_data"].(string)
			if ok {
				var formData map[string]interface{}
				if err := json.Unmarshal([]byte(formDataStr), &formData); err != nil {
					t.Fatalf("failed to unmarshal form_data: %v", err)
				}

				if formData["Content-Type"] != "image/jpeg" {
					t.Errorf("expected '%s' as Content-Type but got '%s'", "image/jpeg", formData["Content-Type"])
				}

				if formData["bucket"] != "book-images" {
					t.Errorf("expected '%s' as bucket but got '%s'", "book-images", formData["bucket"])
				}

				if formData["key"] != "35e45eae-123c-4727-8b46-e6b2ea939e12.jpg" {
					t.Errorf("expected '%s' as key but got '%s'", "35e45eae-123c-4727-8b46-e6b2ea939e12.jpg", formData["key"])
				}

				if formData["policy"] != "eyJleHBpcmF0aW9uIjoiMjAyNS0xMS0wOFQwMzoxODoyNC43NTdaIiwiY29uZGl0aW9ucyI6W1siZXEiLCIkYnVja2V0IiwiYm9vay1pbWFnZXMiXSxbImVxIiwiJGtleSIsIjM1ZTQ1ZWFlLTEyM2MtNDcyNy04YjQ2LWU2YjJlYTkzOWUxMi5qcGciXSxbImVxIiwiJENvbnRlbnQtVHlwZSIsImltYWdlL2pwZWciXSxbImVxIiwiJHgtYW16LWRhdGUiLCIyMDI1MTEwOFQwMzEzMjRaIl0sWyJlcSIsIiR4LWFtei1hbGdvcml0aG0iLCJBV1M0LUhNQUMtU0hBMjU2Il0sWyJlcSIsIiR4LWFtei1jcmVkZW50aWFsIiwibXlBY2Nlc3NLZXkvMjAyNTExMDgvdXMtZWFzdC0xL3MzL2F3czRfcmVxdWVzdCJdLFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLCAxMDI0LCA1MjQyODgwXV19" {
					t.Errorf("expected '%s' as policy but got '%s'", "eyJleHBpcmF0aW9uIjoiMjAyNS0xMS0wOFQwMzoxODoyNC43NTdaIiwiY29uZGl0aW9ucyI6W1siZXEiLCIkYnVja2V0IiwiYm9vay1pbWFnZXMiXSxbImVxIiwiJGtleSIsIjM1ZTQ1ZWFlLTEyM2MtNDcyNy04YjQ2LWU2YjJlYTkzOWUxMi5qcGciXSxbImVxIiwiJENvbnRlbnQtVHlwZSIsImltYWdlL2pwZWciXSxbImVxIiwiJHgtYW16LWRhdGUiLCIyMDI1MTEwOFQwMzEzMjRaIl0sWyJlcSIsIiR4LWFtei1hbGdvcml0aG0iLCJBV1M0LUhNQUMtU0hBMjU2Il0sWyJlcSIsIiR4LWFtei1jcmVkZW50aWFsIiwibXlBY2Nlc3NLZXkvMjAyNTExMDgvdXMtZWFzdC0xL3MzL2F3czRfcmVxdWVzdCJdLFsiY29udGVudC1sZW5ndGgtcmFuZ2UiLCAxMDI0LCA1MjQyODgwXV19", formData["policy"])
				}

				if formData["x-amz-algorithm"] != "AWS4-HMAC-SHA256" {
					t.Errorf("expected '%s' as x-amz-algorithm but got '%s'", "AWS4-HMAC-SHA256", formData["x-amz-algorithm"])
				}

				if formData["x-amz-credential"] != "myAccessKey/20251108/us-east-1/s3/aws4_request" {
					t.Errorf("expected '%s' as x-amz-credential but got '%s'", "myAccessKey/20251108/us-east-1/s3/aws4_request", formData["x-amz-credential"])
				}

				if formData["x-amz-date"] != "20251108T031324Z" {
					t.Errorf("expected '%s' as x-amz-date but got '%s'", "20251108T031324Z", formData["x-amz-date"])
				}

				if formData["x-amz-signature"] != "3025c0c4b29f7ffac992e04f9c1d81e205f2f33f6b4dc44d1d7cc96e5c071453" {
					t.Errorf("expected '%s' as x-amz-signature but got '%s'", "3025c0c4b29f7ffac992e04f9c1d81e205f2f33f6b4dc44d1d7cc96e5c071453", formData["x-amz-signature"])
				}
			} else {
				t.Error("formDataStr should be true but got false")
			}
		} else {
			t.Error("val should be true but got false")
		}
	})
}
