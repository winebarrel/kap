package kap

import (
	"net/url"
)

type Options struct {
	Backend *url.URL `short:"b" required:"" env:"KAP_BACKEND" help:"Backend URL."`
	Port    uint     `short:"p" required:"" env:"KAP_PORT" help:"Listening port."`
	Key     string   `short:"k" required:"" env:"KAP_KEY" help:"Auth key name."`
	Secret  Secret   `short:"s" required:"" env:"KAP_SECRET" help:"Auth secret value."`
}
