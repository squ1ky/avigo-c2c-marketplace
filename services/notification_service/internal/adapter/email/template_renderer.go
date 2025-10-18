package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/config"
	render "github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/usecase"
)

type HTMLTemplateRenderer struct {
	templatesPath string
	templates     *template.Template
}

func NewHTMLTemplateRenderer(cfg *config.EmailConfig) (*HTMLTemplateRenderer, error) {
	templatesPath := filepath.Join(cfg.TemplatesPath, "*.html")
	tmpl, err := template.ParseGlob(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	return &HTMLTemplateRenderer{
		templatesPath: cfg.TemplatesPath,
		templates:     tmpl,
	}, nil
}

func (r *HTMLTemplateRenderer) Render(
	ctx context.Context,
	templateName string,
	data render.TemplateData,
) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("template rendering cancelled: %w", ctx.Err())
	default:
	}

	if filepath.Ext(templateName) == "" {
		templateName += ".html"
	}

	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return buf.String(), nil
}
