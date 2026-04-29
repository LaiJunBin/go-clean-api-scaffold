package greeting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-clean-api-scaffold/internal/app/hello_domain/hello/domain/entity/greeting"
)

func TestNew(t *testing.T) {
	g, err := greeting.New("Alice")

	assert.NoError(t, err)
	assert.Equal(t, 0, g.GetID())
	assert.Equal(t, "Alice", g.GetName())
	assert.Equal(t, "Hello, Alice!", g.GetMessage())
}

func TestNew_EmptyName(t *testing.T) {
	g, err := greeting.New("")

	assert.Nil(t, g)
	assert.EqualError(t, err, "name cannot be empty")
}

func TestReconstruct(t *testing.T) {
	g := greeting.Reconstruct(7, "Bob", "Hello, Bob!")

	assert.Equal(t, 7, g.GetID())
	assert.Equal(t, "Bob", g.GetName())
	assert.Equal(t, "Hello, Bob!", g.GetMessage())
}
