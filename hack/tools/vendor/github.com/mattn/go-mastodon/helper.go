package mastodon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type APIError struct {
	prefix     string
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	errMsg := fmt.Sprintf("%s: %d %s", e.prefix, e.StatusCode, http.StatusText(e.StatusCode))
	if e.Message == "" {
		return errMsg
	}

	return fmt.Sprintf("%s: %s", errMsg, e.Message)
}

// Base64EncodeFileName returns the base64 data URI format string of the file with the file name.
func Base64EncodeFileName(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Base64Encode(file)
}

// Base64Encode returns the base64 data URI format string of the file.
func Base64Encode(file *os.File) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	d := make([]byte, fi.Size())
	_, err = file.Read(d)
	if err != nil {
		return "", err
	}

	return "data:" + http.DetectContentType(d) +
		";base64," + base64.StdEncoding.EncodeToString(d), nil
}

// String is a helper function to get the pointer value of a string.
func String(v string) *string { return &v }

func parseAPIError(prefix string, resp *http.Response) error {
	res := APIError{
		prefix:     prefix,
		StatusCode: resp.StatusCode,
	}
	var e struct {
		Error string `json:"error"`
	}

	json.NewDecoder(resp.Body).Decode(&e)
	if e.Error != "" {
		res.Message = e.Error
	}

	return &res
}
