package greeting

import (
	"errors"
	"fmt"
)

type Greeting struct {
	id      int
	name    string
	message string
}

func New(name string) (*Greeting, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	return &Greeting{
		name:    name,
		message: fmt.Sprintf("Hello, %s!", name),
	}, nil
}

func Reconstruct(id int, name, message string) *Greeting {
	return &Greeting{id: id, name: name, message: message}
}

func (g *Greeting) GetID() int      { return g.id }
func (g *Greeting) GetName() string { return g.name }
func (g *Greeting) GetMessage() string { return g.message }

func (g *Greeting) SetID(id int) {
	g.id = id
}
