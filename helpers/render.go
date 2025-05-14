package helpers

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var TemplateEngine *html.Engine // Set this from main.go

func RenderPartial(templateName string, data fiber.Map) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := TemplateEngine.Render(buf, templateName, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
