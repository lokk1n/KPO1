package servicecenter

import (
	"kpo1/contracts"
)

func (p *Park) Count() int {
	return len(p.items)
}

func (p *Park) AllEnergyPerDay() float64 {
	var total float64

	for _, item := range p.items {
		if consumer, ok := item.(contracts.Energy); ok {
			total += consumer.EnergyPerDay()
		}
	}

	return total
}

func (p *Park) ForBeginners() []contracts.Inventory {
	result := make([]contracts.Inventory, 0)

	for _, item := range p.items {
		transport, ok := item.(contracts.Transport)

		if ok && transport.DifficultyLevel() >= 6 {
			result = append(result, item)
		}
	}

	return result
}

func (p *Park) List() []contracts.Inventory {
	result := make([]contracts.Inventory, 0, len(p.items))

	for _, item := range p.items {
		result = append(result, item)
	}

	return result
}
