package kap

import (
	"net/url"
)

type Options struct {
	Backend *url.URL `short:"b" required:"" env:"GAP_BACKEND" help:"Backend URL."`
	Port    uint     `short:"p" required:"" env:"GAP_PORT" help:"Listening port."`
	Key     string   `short:"k" required:"" env:"GAP_KEY" help:"Auth key name."`
	Secret  string   `short:"s" required:"" env:"GAP_SECRET" help:"Auth secret value."`
}
