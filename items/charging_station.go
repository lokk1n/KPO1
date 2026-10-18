package items

type ChargingStation struct {
	ID          int64
	DailyEnergy float64
}

func (c ChargingStation) GetID() int64 {
	return c.ID
}

func (c ChargingStation) GetName() string {
	return "Charging Station"
}

func (c ChargingStation) EnergyPerDay() float64 {
	return c.DailyEnergy
}
