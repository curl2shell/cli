package c2b

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
)

var defaultRemoteURL string = "http://0.0.0.0:3000/scripts"

const idempotencyKeyHeader = "Idempotency-Key"

type requestPayload struct {
	URL     string `json:"url"`
	Content string `json:"content"`

	idempotencyKey string
	uploadToken    string
}

func runCurl(args []string) ([]byte, error) {
	// execute underlying curl
	cmd := exec.Command("curl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func uploadResults(payload requestPayload) error {
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := defaultRemoteURL

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("unable to create upload request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(idempotencyKeyHeader, payload.idempotencyKey)

	req.SetBasicAuth("", payload.uploadToken)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed: %s", string(body))
	}

	return nil
}

// get the last URL in the arguments list
// only https URLs are considered for security reasons
func findURL(args []string) (string, bool) {
	for i := len(args) - 1; i >= 0; i-- {
		u, err := url.ParseRequestURI(args[i])
		if err != nil {
			continue
		}
		if u.Scheme == "https" || u.Hostname() == "localhost" {
			return args[i], true
		}
	}

	return "", false
}

func Run(args []string) error {
	out, err := runCurl(args)
	if err != nil {
		return err
	}

	url, foundURL := findURL(args)
	if !foundURL {
		return nil
	}

	// upload payload
	payload := requestPayload{
		URL:            url,
		Content:        string(out),
		idempotencyKey: "idempotencyKey",
		uploadToken:    "uploadToken",
	}

	err = uploadResults(payload)
	if err != nil {
		return err
	}

	fmt.Print(string(out))
	return nil
}
