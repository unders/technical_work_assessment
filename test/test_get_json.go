package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func GetJSON(t *testing.T, url string, wantStatusCode int, respBody any) {
	t.Helper()

	resp, err := http.Get(url) //#nosec
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if wantStatusCode != resp.StatusCode {
		t.Errorf("\nContext: resp.StatusCode\n   Want: %d\n    Got: %d\n", wantStatusCode, resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		t.Fatal(err)
	}
}
