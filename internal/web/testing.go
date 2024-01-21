package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeTestRouter(t *testing.T) *mux.Router {
	repos := core.MakeTestRepositories(t)
	r, err := newRouter(&Config{}, repos)
	require.Nil(t, err)
	return r
}
