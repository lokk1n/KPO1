package servicecenter

import "kpo1/contracts"

type Inspector interface {
	Check(item contracts.Inventory) bool
}
