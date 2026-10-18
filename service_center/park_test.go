package servicecenter_test

import (
	"math"
	"testing"

	"kpo1/contracts"
	"kpo1/items"
	"kpo1/service_center"
)

type StubInspector struct {
	Result bool
	Calls  int
}

func (s *StubInspector) Check(item contracts.Inventory) bool {
	s.Calls++
	return s.Result
}

func TestParkAdmission(t *testing.T) {
	tests := []struct {
		name      string
		accept    bool
		wantCount int
	}{
		{"accept", true, 1},
		{"reject", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inspector := &StubInspector{Result: tt.accept}
			park := servicecenter.NewPark(inspector)

			scooter := items.Scooter{
				ID:           1,
				DailyEnergy:  1.2,
				BeginnerEase: 7,
			}

			err := park.Add(scooter)

			if (err == nil) != tt.accept {
				t.Fatalf("решение инспектора: %v, ошибка: %v",
					tt.accept, err)
			}

			if park.Count() != tt.wantCount {
				t.Fatalf("ожидали %d предметов, получили %d",
					tt.wantCount, park.Count())
			}

			if inspector.Calls != 1 {
				t.Fatalf("ожидали один техосмотр, получили %d",
					inspector.Calls)
			}
		})
	}
}

func TestDuplicateIDDoesNotReplaceItem(t *testing.T) {
	inspector := &StubInspector{Result: true}
	park := servicecenter.NewPark(inspector)

	if err := park.Add(items.Helmet{ID: 1}); err != nil {
		t.Fatal(err)
	}

	err := park.Add(items.Bicycle{
		ID:           1,
		BeginnerEase: 7,
	})

	if err == nil {
		t.Fatal("ожидали ошибку повторного ID")
	}

	list := park.List()

	if len(list) != 1 {
		t.Fatalf("ожидали один предмет, получили %d", len(list))
	}

	if list[0].GetName() != "Helmet" {
		t.Fatal("исходный шлем был заменён")
	}

	if inspector.Calls != 1 {
		t.Fatal("повторный ID должен отклоняться до техосмотра")
	}
}

func TestInvalidItemDoesNotReachInspector(t *testing.T) {
	inspector := &StubInspector{Result: true}
	park := servicecenter.NewPark(inspector)

	err := park.Add(items.Bicycle{
		ID:           1,
		BeginnerEase: 0,
	})

	if err == nil {
		t.Fatal("ожидали ошибку некорректной простоты")
	}

	if park.Count() != 0 {
		t.Fatal("некорректный предмет попал в парк")
	}

	if inspector.Calls != 0 {
		t.Fatal("некорректный предмет дошёл до техосмотра")
	}
}

func TestMixedInventoryReports(t *testing.T) {
	park := servicecenter.NewPark(&StubInspector{Result: true})

	inventory := []contracts.Inventory{
		items.Scooter{
			ID:           1,
			DailyEnergy:  0.7,
			BeginnerEase: 5,
		},
		items.Bicycle{
			ID:           2,
			BeginnerEase: 6,
		},
		items.ElectricBicycle{
			ID:           3,
			DailyEnergy:  1.1,
			BeginnerEase: 10,
		},
		items.Helmet{ID: 4},
		items.ChargingStation{
			ID:          5,
			DailyEnergy: 0.4,
		},
	}

	for _, item := range inventory {
		if err := park.Add(item); err != nil {
			t.Fatal(err)
		}
	}

	if park.Count() != 5 {
		t.Fatalf("ожидали 5 предметов, получили %d", park.Count())
	}

	if got := park.AllEnergyPerDay(); math.Abs(got-2.2) > 1e-9 {
		t.Fatalf("ожидали энергию 2.2, получили %g", got)
	}

	listed := make(map[int64]string)

	for _, item := range park.List() {
		listed[item.GetID()] = item.GetName()
	}

	for _, item := range inventory {
		if listed[item.GetID()] != item.GetName() {
			t.Fatalf("в отчёте отсутствует предмет %d", item.GetID())
		}
	}

	beginners := park.ForBeginners()

	if len(beginners) != 2 {
		t.Fatalf("ожидали 2 предмета для новичков, получили %d",
			len(beginners))
	}

	found := make(map[int64]bool)

	for _, item := range beginners {
		found[item.GetID()] = true
	}

	if !found[2] || !found[3] {
		t.Fatal("для новичков должны подходить велосипед 2 и электровелосипед 3")
	}
}

func TestListDoesNotExposeInternalCollection(t *testing.T) {
	park := servicecenter.NewPark(&StubInspector{Result: true})

	if err := park.Add(items.Helmet{ID: 1}); err != nil {
		t.Fatal(err)
	}

	list := park.List()
	list[0] = items.Helmet{ID: 99}

	if park.List()[0].GetID() != 1 {
		t.Fatal("изменение возвращённого среза изменило состав парка")
	}
}

func TestEmptyParkReports(t *testing.T) {
	park := servicecenter.NewPark(&StubInspector{Result: true})

	if park.Count() != 0 {
		t.Fatal("пустой парк содержит предметы")
	}

	if park.AllEnergyPerDay() != 0 {
		t.Fatal("энергия пустого парка должна быть нулевой")
	}

	if len(park.List()) != 0 || len(park.ForBeginners()) != 0 {
		t.Fatal("списки пустого парка должны быть пустыми")
	}
}
