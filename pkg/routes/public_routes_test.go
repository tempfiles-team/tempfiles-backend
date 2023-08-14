package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

func SetServer() *fiber.App {

	if err := godotenv.Load("../../.env.test"); err != nil {
		panic(err)
	}

	utils.ReadyComponent()
	app := fiber.New()
	PublicRoutes(app)

	return app
}

// TestUploadTest tests the file upload functionality
func TestUploadTest(t *testing.T) {

	defer os.RemoveAll("./tmp")

	app := SetServer()

	tests := []struct {
		description   string
		expectedError bool
		expectedCode  int
		header        map[string]string
	}{
		{
			description:   "just file upload test",
			expectedError: false,
			expectedCode:  200,
			header:        map[string]string{"": ""},
		},
		{
			description:   "file upload with X-Download-Limit header",
			expectedError: false,
			expectedCode:  200,
			header:        map[string]string{"X-Download-Limit": "100"},
		},
		{
			description:   "file upload with X-Time-Limit header",
			expectedError: false,
			expectedCode:  200,
			header:        map[string]string{"X-Time-Limit": "100"},
		},
	}

	for _, test := range tests {

		// Write some test data to the temporary file
		testData := []byte("Test file upload content")
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = io.Copy(part, bytes.NewReader(testData))
		if err != nil {
			t.Fatalf("Failed to copy file content to form file: %v", err)
		}
		writer.Close()

		// Create a test request
		req := httptest.NewRequest("POST", "/upload", body)

		for key, value := range test.header {
			req.Header.Set(key, value)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Perform the request and record the response
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform test request: %v", err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		fmt.Println(string(respBody))

		var jsonData utils.Response

		err = json.Unmarshal(respBody, &jsonData)

		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		// Check if the file was actually uploaded and saved
		uploadDir := "./tmp/" + jsonData.Data.(map[string]interface{})["id"].(string)

		uploadedFilePath := fmt.Sprintf("%s/%s", uploadDir, "test.txt")
		if _, err := os.Stat(uploadedFilePath); os.IsNotExist(err) {
			t.Errorf("Uploaded file not found at path %q", uploadedFilePath)
		}
	}
}
