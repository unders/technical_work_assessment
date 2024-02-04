// nolint: forbidigo,funlen
package test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
)

// RoundTripperRecorder is used for recording RoundTripper calls.
type RoundTripperRecorder struct {
	counter          int
	prefix           string
	update           bool
	printRequest     bool
	realRoundTripper http.RoundTripper
}

// TransportRecorder returns a RoundTripper interface that records RoundTrip calls.
func TransportRecorder(updateFile, printRequest bool, prefix string) http.RoundTripper {
	rt := http.DefaultTransport

	return &RoundTripperRecorder{
		prefix:           prefix,
		realRoundTripper: rt,
		update:           updateFile,
		printRequest:     printRequest,
	}
}

// RoundTrip intercepts RoundTrip calls and saves the response to a file when
// `update` is true; when `update`is false it does not make any calls, instead
// it reads from a file to construct the response.
func (r *RoundTripperRecorder) RoundTrip(req *http.Request) (*http.Response, error) {
	r.counter++

	if r.printRequest {
		printRequest(req)
	}

	path := filePath(r.counter, r.prefix, req)

	//
	// Do the real request and update the recording with the response
	//
	if r.update {
		var (
			dump []byte
			err  error
		)

		var resp *http.Response

		if resp, err = r.realRoundTripper.RoundTrip(req); err != nil {
			return nil, err
		}

		var respBody string

		resp.Body, respBody, err = cloneResponseBody(resp)
		if err != nil {
			return nil, err
		}

		if dump, err = httputil.DumpResponse(resp, true); err != nil {
			return nil, err
		}

		err = os.WriteFile(path, dump, 0o644) // #nosec

		if r.printRequest {
			printResponse(resp, respBody)
		}

		return resp, err
	}

	//
	// read response from file
	//
	resp, err := readResponse(path, req)
	if err != nil {
		return nil, err
	}

	var (
		respBody string
		body     io.ReadCloser
	)

	body, respBody, err = cloneResponseBody(resp)
	if err != nil {
		return nil, err
	}

	resp.Body = body

	if r.printRequest {
		printResponse(resp, respBody)
	}

	return resp, err
}

//
// private
//

func printRequest(req *http.Request) {
	fmt.Println("--------REQUEST-----------")
	fmt.Printf("%s %s://%s%s\n", req.Method, req.URL.Scheme, req.Host, req.URL.RequestURI())
	fmt.Printf("\nHeader{\n")

	for k, v := range req.Header {
		fmt.Printf("    %s: %s\n", k, strings.Join(v, ","))
	}

	fmt.Printf("}\n")

	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodDelete:
	default:
		reader, err := req.GetBody()
		if err != nil {
			fmt.Printf("error getting body")

			break
		}

		defer func() {
			if err := reader.Close(); err != nil {
				fmt.Printf("reader.Close() failed; %s", err)
			}
		}()

		b, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("error reading body")

			break
		}

		fmt.Printf("\nBody: %s\n", b)
	}
}

func cloneResponseBody(resp *http.Response) (io.ReadCloser, string, error) {
	if resp == nil {
		return nil, "", errors.New("resp cannot be nil")
	}

	if resp.ContentLength == 0 {
		return resp.Body, "", nil
	}

	var b bytes.Buffer

	_, err := b.ReadFrom(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return io.NopCloser(&b), b.String(), nil
}

func printResponse(resp *http.Response, body string) {
	if resp != nil {
		fmt.Println("\n--------RESPONSE-----------")
		fmt.Printf("\nStatus: %s\n", resp.Status)
		fmt.Printf("\nHeaders{\n")

		for k, v := range resp.Header {
			fmt.Printf("    %s: %s\n", k, strings.Join(v, ","))
		}

		fmt.Printf("}\n")
		fmt.Printf("Body:\n %s\n", body)
	}
}

// readResponse returns a response from a file.
func readResponse(path string, req *http.Request) (*http.Response, error) {
	dump, err := os.ReadFile(path) // #nosec
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(bytes.NewReader(dump)), req)
}

func filePath(counter int, prefix string, req *http.Request) string {
	method := strings.ToLower(req.Method)
	host := "__" + strings.ReplaceAll(req.Host, "/", "__")
	path := strings.ReplaceAll(req.URL.EscapedPath(), "/", "__") + ".response"

	if prefix == "" {
		prefix = "request"
	}

	if counter > 1 {
		prefix = fmt.Sprintf("%s__call__%d__", prefix, counter)
	}

	return filepath.Join("testdata", fmt.Sprintf("%s_%s%s%s", prefix, method, host, path))
}
