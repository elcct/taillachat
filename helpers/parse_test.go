package helpers

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestParse(t *testing.T) {
	tpl, err := template.New("test").Parse("<h1>{{ .Test }}</h1>")
	assert.Nil(t, err)

	result := Parse(tpl, "test", map[string]string{"Test": "test"})

	expected := "<h1>test</h1>"
	assert.Equal(t, expected, result)
}
