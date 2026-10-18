package items

type Helmet struct {
	ID int64
}

func (h Helmet) GetID() int64 {
	return h.ID
}

func (h Helmet) GetName() string {
	return "Helmet"
}
