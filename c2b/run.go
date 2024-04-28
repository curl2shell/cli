package c2b

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

var defaultRemoteURL string = "http://0.0.0.0:3000/scripts"

const idempotencyKeyHeader = "Idempotency-Key"

type requestPayload struct {
	URL    string
	Script string

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

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Run(args []string) error {
	out, err := runCurl(args)
	if err != nil {
		return err
	}

	// upload payload
	payload := requestPayload{
		URL:            "url",
		Script:         string(out),
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
