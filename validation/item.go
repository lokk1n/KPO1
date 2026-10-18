package validation

import (
	"fmt"
	"math"

	"kpo1/contracts"
)

func Validate(item contracts.Inventory) error {
	if item == nil {
		return fmt.Errorf("предмет не задан")
	}

	if item.GetID() <= 0 {
		return fmt.Errorf(
			"инвентарный номер должен быть положительным: %d",
			item.GetID(),
		)
	}

	if transport, ok := item.(contracts.Transport); ok {
		level := transport.DifficultyLevel()

		if level < 1 || level > 10 {
			return fmt.Errorf(
				"простота для новичка должна быть от 1 до 10: %d",
				level,
			)
		}
	}

	if consumer, ok := item.(contracts.Energy); ok {
		energy := consumer.EnergyPerDay()

		if math.IsNaN(energy) || math.IsInf(energy, 0) {
			return fmt.Errorf(
				"расход энергии должен быть конечным числом",
			)
		}

		if energy < 0 {
			return fmt.Errorf(
				"расход энергии не может быть отрицательным: %g",
				energy,
			)
		}
	}

	return nil
}
