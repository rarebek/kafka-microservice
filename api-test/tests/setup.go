package tests

import (
	kv2 "4microservice/api_gateway/api-test/storage/kv"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

const testHost = "http://localhost"

func SetupMinimumInstance(path string) error {

	_ = path
	cnf := viper.New()
	cnf.Set("mode", "test")

	//rClient := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//})

	kv2.Init(kv2.NewInMemory())
	return nil
}

func Serve(handler func(c *gin.Context), method, uri string, body []byte) (*httptest.ResponseRecorder, error) {
	router := gin.New()

	gin.SetMode(gin.TestMode)

	switch method {
	case http.MethodPost:
		router.POST(uri, handler)
	case http.MethodGet:
		router.GET(uri, handler)
	case http.MethodDelete:
		router.DELETE(uri, handler)
	case http.MethodPatch:
		router.PATCH(uri, handler)

	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w, nil
}

func NewResponse() *http.Response {
	return &http.Response{}
}

func NewRequest(method string, uri string, body []byte) *http.Request {
	req, err := http.NewRequest(method, testHost+uri, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Host", "localhost")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("X-Forwarded-For", "79.104.42.249")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return req
}

func OpenFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}
