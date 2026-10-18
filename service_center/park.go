package servicecenter

import (
	"fmt"
	"kpo1/contracts"
	"kpo1/validation"
)

type Park struct {
	inspector Inspector
	items     map[int64]contracts.Inventory
}

func NewPark(inspector Inspector) *Park {
	return &Park{
		inspector: inspector,
		items:     make(map[int64]contracts.Inventory),
	}
}

func (p *Park) Add(item contracts.Inventory) error {
	if err := validation.Validate(item); err != nil {
		return err
	}

	id := item.GetID()

	if _, exists := p.items[id]; exists {
		return fmt.Errorf("инвентарный номер %d уже занят", id)
	}

	if !p.inspector.Check(item) {
		return fmt.Errorf("предмет %d не прошёл техосмотр", id)
	}

	p.items[id] = item
	return nil
}
