package timewebdns

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/timeweb"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *timeweb.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.timeweb",
		New: func() caddy.Module { return &Provider{new(timeweb.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.ApiToken = repl.ReplaceAll(p.Provider.ApiToken, "")
	p.Provider.ApiURL = repl.ReplaceAll(p.Provider.ApiURL, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	timeweb [<api_token>] {
//	    api_token <api_token>
//	    api_url <timeweb_domain>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.ApiToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.ApiToken != "" {
					return d.Err("API token already set")
				}
				p.Provider.ApiToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_url":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.ApiURL != "" {
					return d.Err("API url already set")
				}
				p.Provider.ApiURL = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.ApiToken == "" {
		return d.Err("missing API token")
	}

	if p.Provider.ApiURL == "" {
		p.Provider.ApiURL = "https://api.timeweb.cloud/"
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
