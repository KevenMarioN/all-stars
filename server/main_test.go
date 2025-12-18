package server_test

import (
	"testing"

	"github.com/rs/zerolog"
)
func TestMain(m *testing.M){
	zerolog.SetGlobalLevel(zerolog.Disabled)
	m.Run()
}
