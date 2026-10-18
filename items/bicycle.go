package items

type Bicycle struct {
	ID           int64
	BeginnerEase int64
}

func (b Bicycle) DifficultyLevel() int64 {
	return b.BeginnerEase
}

func (b Bicycle) GetID() int64 {
	return b.ID
}

func (b Bicycle) GetName() string {
	return "Bicycle"
}
