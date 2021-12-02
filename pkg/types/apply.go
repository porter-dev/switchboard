package types

import "github.com/rs/zerolog"

type ApplyOpts struct {
	BasePath       string
	Logger         *zerolog.Logger
	ResourceLogger *zerolog.Logger
}
