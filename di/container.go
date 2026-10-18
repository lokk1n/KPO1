package di

import "fmt"

type Factory func(c *Container) (any, error)

type Container struct {
	factories map[string]Factory
}

func NewContainer() *Container {
	return &Container{
		factories: make(map[string]Factory),
	}
}

func (c *Container) Register(name string, factory Factory) error {
	if factory == nil {
		return fmt.Errorf("фабрика %q не задана", name)
	}

	if _, exists := c.factories[name]; exists {
		return fmt.Errorf("Зависимость %q уже зарегестрирована", name)
	}

	c.factories[name] = factory
	return nil
}

func (c *Container) Resolve(name string) (any, error) {
	factory, exists := c.factories[name]

	if !exists {
		return nil, fmt.Errorf("зависимость %q не зарегестрирована", name)
	}

	return factory(c)
}
