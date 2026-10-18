package validation_test

import (
	"math"
	"testing"

	"kpo1/contracts"
	"kpo1/items"
	"kpo1/validation"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		item  contracts.Inventory
		valid bool
	}{
		{
			"nil_item",
			nil,
			false,
		},
		{
			"zero_id",
			items.Helmet{ID: 0},
			false,
		},
		{
			"negative_id",
			items.Helmet{ID: -1},
			false,
		},
		{
			"valid_helmet",
			items.Helmet{ID: 1},
			true,
		},
		{
			"ease_below_minimum",
			items.Bicycle{ID: 1, BeginnerEase: 0},
			false,
		},
		{
			"ease_above_maximum",
			items.Bicycle{ID: 1, BeginnerEase: 11},
			false,
		},
		{
			"ease_minimum",
			items.Bicycle{ID: 1, BeginnerEase: 1},
			true,
		},
		{
			"ease_maximum",
			items.Bicycle{ID: 1, BeginnerEase: 10},
			true,
		},
		{
			"negative_energy",
			items.ChargingStation{ID: 1, DailyEnergy: -0.1},
			false,
		},
		{
			"zero_energy",
			items.ChargingStation{ID: 1, DailyEnergy: 0},
			true,
		},
		{
			"nan_energy",
			items.ChargingStation{ID: 1, DailyEnergy: math.NaN()},
			false,
		},
		{
			"positive_infinity",
			items.ChargingStation{ID: 1, DailyEnergy: math.Inf(1)},
			false,
		},
		{
			"negative_infinity",
			items.ChargingStation{ID: 1, DailyEnergy: math.Inf(-1)},
			false,
		},
		{
			"valid_scooter",
			items.Scooter{ID: 1, DailyEnergy: 1.2, BeginnerEase: 6},
			true,
		},
		{
			"valid_electric_bicycle",
			items.ElectricBicycle{ID: 1, DailyEnergy: 0.8, BeginnerEase: 7},
			true,
		},
		{
			"high_but_valid_energy",
			items.ChargingStation{ID: 1, DailyEnergy: 5},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.item)

			if (err == nil) != tt.valid {
				t.Fatalf("ожидали valid=%v, получили ошибку: %v",
					tt.valid, err)
			}
		})
	}
}
