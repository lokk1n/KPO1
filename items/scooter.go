package items

type Scooter struct {
	ID           int64
	DailyEnergy  float64
	BeginnerEase int64
}

func (s Scooter) EnergyPerDay() float64 {
	return s.DailyEnergy
}

func (s Scooter) DifficultyLevel() int64 {
	return s.BeginnerEase
}

func (s Scooter) GetID() int64 {
	return s.ID
}

func (s Scooter) GetName() string {
	return "Scooter"
}
