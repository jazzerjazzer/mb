package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	ts := httptest.NewServer(Method(http.MethodPost, fakeHandler))
	defer ts.Close()

	var u bytes.Buffer
	u.WriteString(string(ts.URL))
	u.WriteString("/")

	res, err := http.Get(u.String())
	require.NoError(t, err)
	if res != nil {
		defer res.Body.Close()
	}
	require.Equal(t, res.StatusCode, http.StatusMethodNotAllowed)

	res, err = http.Post(u.String(), "", nil)
	require.NoError(t, err)
	if res != nil {
		defer res.Body.Close()
	}
	require.Equal(t, res.StatusCode, http.StatusOK)
}
