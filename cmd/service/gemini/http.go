package gemini

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jeraldyik/crypto_dca_go/cmd/util"
)

type Api struct {
	url    string
	key    string
	secret string
}

func New(key, secret string, live bool) *Api {
	var url string
	if live {
		url = baseURL
	} else {
		url = sandboxURL
	}
	return &Api{url: url, key: key, secret: secret}
}

// buildHeader handles the conversion of post parameters into headers formatted
// according to Gemini specification. Resulting headers include the API key,
// the payload and the signature.
func (api *Api) buildHeader(req map[string]any) http.Header {

	reqStr, _ := json.Marshal(req)
	payload := base64.StdEncoding.EncodeToString([]byte(reqStr))

	mac := hmac.New(sha512.New384, []byte(api.secret))
	if _, err := mac.Write([]byte(payload)); err != nil {
		panic(err)
	}

	signature := hex.EncodeToString(mac.Sum(nil))

	header := http.Header{}
	header.Set("Content-Length", "0")
	header.Set("Content-Type", "text/plain")
	header.Set("Cache-Control", "no-cache")
	header.Set("X-GEMINI-APIKEY", api.key)
	header.Set("X-GEMINI-PAYLOAD", payload)
	header.Set("X-GEMINI-SIGNATURE", signature)

	return header
}

// request makes the HTTP request to Gemini and handles any returned errors
func (api *Api) request(verb, path string, params map[string]any) ([]byte, error) {
	url := api.url + path

	req, err := http.NewRequest(verb, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	if params != nil {
		if verb == http.MethodGet {
			q := req.URL.Query()
			for key, val := range params {
				q.Add(key, val.(string))
			}
			req.URL.RawQuery = q.Encode()
		} else {
			req.Header = api.buildHeader(params)
		}
	}

	log.Printf("[gemini.request] request verb:%s, url:%s, params:%+v\n", verb, url, params)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[gemini.request] response: %v\n", util.SafeJsonDump(resp))

	if resp.StatusCode != 200 {
		statusCode := fmt.Sprintf("HTTP Status Code: %d", resp.StatusCode)
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "API entry point has moved, see Location: header. Most likely an http: to https: redirect.")
		} else if resp.StatusCode == 400 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "Auction not open or paused, ineligible timing, market not open, or the request was malformed; in the case of a private API request, missing or malformed Gemini private API authentication headers")
		} else if resp.StatusCode == 403 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "The API key is missing the role necessary to access this private API endpoint")
		} else if resp.StatusCode == 404 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "Unknown API entry point or Order not found")
		} else if resp.StatusCode == 406 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "Insufficient Funds")
		} else if resp.StatusCode == 429 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "Rate Limiting was applied")
		} else if resp.StatusCode == 500 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "The server encountered an error")
		} else if resp.StatusCode == 502 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "Technical issues are preventing the request from being satisfied")
		} else if resp.StatusCode == 503 {
			return nil, fmt.Errorf("%s\n%s", statusCode, "The exchange is down for maintenance")
		}
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[gemini.request] response.body: %v\n", string(body))

	return body, nil
}
