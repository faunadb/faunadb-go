package faunadb

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReturnBadRequestOn400(t *testing.T) {
	err := checkForResponseErrors(&http.Response{StatusCode: 400})
	require.Equal(t, BadRequest{}, err)
}
