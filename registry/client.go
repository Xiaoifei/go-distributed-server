package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Responded with code %v", res.StatusCode)
	}

	return nil
}
func ShutdownService(url string) error {
	buf := bytes.NewBufferString(url)
	req, err := http.NewRequest("DELETE", ServicesURL, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove service. Responded with code %v", res.StatusCode)
	}

	return nil
}
