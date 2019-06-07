package testhelper

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)



func init() {
	// Ensure we seed the RNG so generated names aren't deterministic
	rand.Seed(time.Now().UTC().UnixNano())

}

var (
	Mux    *http.ServeMux
	Server *httptest.Server
)

func SetupHTTP() {
	Mux = http.NewServeMux()
	Server = httptest.NewServer(Mux)
}

func TeardownHTTP() {
	Server.Close()
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func AssertError(t *testing.T, err error, expected string) {
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error to contain '%s', but was '%s'", expected, err)
	}
}

func RandomString(prefix string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Intn(26) + 97)
	}

	return fmt.Sprintf("%s%s", prefix, string(b))
}
