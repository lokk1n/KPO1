package servicecenter_test

import (
	"testing"

	"kpo1/contracts"
	"kpo1/items"
	"kpo1/service_center"
)

func TestServiceCenterCheck(t *testing.T) {
	center := servicecenter.NewServiceCenter()

	tests := []struct {
		name string
		item contracts.Inventory
		want bool
	}{
		{
			"helmet_without_energy",
			items.Helmet{ID: 1},
			true,
		},
		{
			"bicycle_without_energy",
			items.Bicycle{ID: 1, BeginnerEase: 6},
			true,
		},
		{
			"below_minimum",
			items.Scooter{DailyEnergy: 0.09},
			false,
		},
		{
			"minimum",
			items.Scooter{DailyEnergy: 0.1},
			true,
		},
		{
			"maximum",
			items.Scooter{DailyEnergy: 2.0},
			true,
		},
		{
			"above_maximum",
			items.Scooter{DailyEnergy: 2.01},
			false,
		},
		{
			"charging_station_accepted",
			items.ChargingStation{DailyEnergy: 1.5},
			true,
		},
		{
			"charging_station_rejected",
			items.ChargingStation{DailyEnergy: 3},
			false,
		},
		{
			"electric_bicycle",
			items.ElectricBicycle{DailyEnergy: 1},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := center.Check(tt.item); got != tt.want {
				t.Fatalf("ожидали %v, получили %v", tt.want, got)
			}
		})
	}
}
