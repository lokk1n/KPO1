package items

type ElectricBicycle struct {
	ID           int64
	DailyEnergy  float64
	BeginnerEase int64
}

func (e ElectricBicycle) EnergyPerDay() float64 {
	return e.DailyEnergy
}

func (e ElectricBicycle) DifficultyLevel() int64 {
	return e.BeginnerEase
}

func (e ElectricBicycle) GetID() int64 {
	return e.ID
}

func (e ElectricBicycle) GetName() string {
	return "Electric Bicycle"
}
