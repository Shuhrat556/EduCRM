package config

import "strings"

// parseCORSOrigins splits a comma-separated list of origins; empty entries are dropped.
func (c *Config) parseCORSOrigins() []string {
	if c == nil {
		return nil
	}
	raw := strings.TrimSpace(c.CORS.AllowedOrigins)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// EffectiveCORSOrigins returns origins for middleware: non-production defaults to * when unset.
func (c *Config) EffectiveCORSOrigins() []string {
	origins := c.parseCORSOrigins()
	if len(origins) > 0 {
		return origins
	}
	if !c.IsProduction() {
		return []string{"*"}
	}
	return nil
}
