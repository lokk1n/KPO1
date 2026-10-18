package servicecenter

import "kpo1/contracts"

type ServiceCenter struct {
}

func NewServiceCenter() *ServiceCenter {
	return &ServiceCenter{}
}

func (s ServiceCenter) Check(i contracts.Inventory) bool {
	consumer, ok := i.(contracts.Energy)

	if !ok {
		return true
	}

	energy := consumer.EnergyPerDay()
	return energy >= 0.1 && energy <= 2.0
}
