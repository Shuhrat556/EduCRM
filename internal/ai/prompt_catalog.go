package ai

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// PromptCatalog loads system prompts and user prompt templates from disk with built-in fallbacks.
type PromptCatalog struct {
	dir string
}

// NewPromptCatalog uses dir when non-empty; missing files fall back to defaults.
func NewPromptCatalog(dir string) *PromptCatalog {
	return &PromptCatalog{dir: strings.TrimSpace(dir)}
}

// Render returns system text and rendered user prompt for a named scenario (e.g. "debtors_summary").
func (c *PromptCatalog) Render(scenario string, dataJSON string) (system string, user string, err error) {
	system, err = c.loadSystem(scenario)
	if err != nil {
		return "", "", err
	}
	user, err = c.renderUser(scenario, dataJSON)
	if err != nil {
		return "", "", err
	}
	return system, user, nil
}

func (c *PromptCatalog) loadSystem(scenario string) (string, error) {
	name := scenario + ".system.txt"
	if c.dir != "" {
		if b, err := os.ReadFile(filepath.Join(c.dir, name)); err == nil {
			return strings.TrimSpace(string(b)), nil
		}
	}
	s, ok := defaultSystemPrompts[scenario]
	if !ok {
		return "", fmt.Errorf("unknown ai prompt scenario %q", scenario)
	}
	return s, nil
}

func (c *PromptCatalog) renderUser(scenario string, dataJSON string) (string, error) {
	tmplName := scenario + ".user.tmpl"
	var tmplStr string
	if c.dir != "" {
		if b, err := os.ReadFile(filepath.Join(c.dir, tmplName)); err == nil {
			tmplStr = string(b)
		}
	}
	if tmplStr == "" {
		var ok bool
		tmplStr, ok = defaultUserTemplates[scenario]
		if !ok {
			return "", fmt.Errorf("unknown ai user template for scenario %q", scenario)
		}
	}
	t, err := template.New(tmplName).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parse user template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, struct{ DataJSON string }{DataJSON: dataJSON}); err != nil {
		return "", fmt.Errorf("execute user template: %w", err)
	}
	return strings.TrimSpace(buf.String()), nil
}
