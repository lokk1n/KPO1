package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"kpo1/contracts"
	"kpo1/di"
	"kpo1/items"
	"kpo1/service_center"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}

func run() error {
	container := di.NewContainer()

	err := container.Register("inspector",
		func(c *di.Container) (any, error) {
			return servicecenter.NewServiceCenter(), nil
		},
	)
	if err != nil {
		return err
	}

	err = container.Register("park",
		func(c *di.Container) (any, error) {
			value, err := c.Resolve("inspector")
			if err != nil {
				return nil, err
			}

			inspector, ok := value.(servicecenter.Inspector)
			if !ok {
				return nil, fmt.Errorf("неверный тип инспектора")
			}

			return servicecenter.NewPark(inspector), nil
		},
	)
	if err != nil {
		return err
	}

	value, err := container.Resolve("park")
	if err != nil {
		return err
	}

	park, ok := value.(*servicecenter.Park)
	if !ok {
		return fmt.Errorf("неверный тип парка")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 — Добавить предмет")
		fmt.Println("2 — Всё имущество")
		fmt.Println("3 — Статистика")
		fmt.Println("4 — Транспорт для новичков")
		fmt.Println("0 — Выход")

		action, err := readInt(scanner, "Действие: ", 0, 4)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch action {
		case 0:
			fmt.Println("До свидания!")
			return nil

		case 1:
			item, err := readItem(scanner)
			if err == io.EOF {
				fmt.Println("\nВвод завершён, предмет не добавлен.")
				return nil
			}
			if err != nil {
				return err
			}

			if err := park.Add(item); err != nil {
				fmt.Println("Не удалось добавить:", err)
			} else {
				fmt.Printf("Добавлен %s с ID %d\n",
					item.GetName(), item.GetID())
			}

		case 2:
			list := park.List()

			if len(list) == 0 {
				fmt.Println("Имущества пока нет.")
			}

			for _, item := range list {
				fmt.Printf("ID: %d, название: %s\n",
					item.GetID(), item.GetName())
			}

		case 3:
			fmt.Printf("Всего предметов: %d\n", park.Count())
			fmt.Printf("Суточная энергия: %.2f кВт·ч\n",
				park.AllEnergyPerDay())

		case 4:
			list := park.ForBeginners()

			if len(list) == 0 {
				fmt.Println("Подходящего транспорта нет.")
			}

			for _, item := range list {
				fmt.Printf("ID: %d, название: %s\n",
					item.GetID(), item.GetName())
			}
		}
	}
}

func readInt(
	scanner *bufio.Scanner,
	prompt string,
	min, max int64,
) (int64, error) {
	for {
		fmt.Print(prompt)

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		text := strings.TrimSpace(scanner.Text())
		value, err := strconv.ParseInt(text, 10, 64)

		if err != nil || value < min || value > max {
			fmt.Printf("Введите целое число от %d до %d.\n", min, max)
			continue
		}

		return value, nil
	}
}

func readEnergy(scanner *bufio.Scanner) (float64, error) {
	for {
		fmt.Print("Суточный расход энергии, кВт·ч: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		text := strings.TrimSpace(scanner.Text())
		text = strings.ReplaceAll(text, ",", ".")

		value, err := strconv.ParseFloat(text, 64)

		if err != nil ||
			math.IsNaN(value) ||
			math.IsInf(value, 0) ||
			value < 0 {
			fmt.Println("Введите конечное неотрицательное число.")
			continue
		}

		return value, nil
	}
}

func readItem(scanner *bufio.Scanner) (contracts.Inventory, error) {
	fmt.Println("\nТип предмета:")
	fmt.Println("1 — Самокат")
	fmt.Println("2 — Велосипед")
	fmt.Println("3 — Электровелосипед")
	fmt.Println("4 — Шлем")
	fmt.Println("5 — Зарядная станция")

	kind, err := readInt(scanner, "Тип: ", 1, 5)
	if err != nil {
		return nil, err
	}

	id, err := readInt(
		scanner,
		"Инвентарный номер: ",
		1,
		math.MaxInt64,
	)
	if err != nil {
		return nil, err
	}

	var ease int64

	if kind <= 3 {
		ease, err = readInt(
			scanner,
			"Простота для новичка: ",
			1,
			10,
		)
		if err != nil {
			return nil, err
		}
	}

	var energy float64

	if kind == 1 || kind == 3 || kind == 5 {
		energy, err = readEnergy(scanner)
		if err != nil {
			return nil, err
		}
	}

	switch kind {
	case 1:
		return items.Scooter{
			ID:           id,
			DailyEnergy:  energy,
			BeginnerEase: ease,
		}, nil

	case 2:
		return items.Bicycle{
			ID:           id,
			BeginnerEase: ease,
		}, nil

	case 3:
		return items.ElectricBicycle{
			ID:           id,
			DailyEnergy:  energy,
			BeginnerEase: ease,
		}, nil

	case 4:
		return items.Helmet{ID: id}, nil

	default:
		return items.ChargingStation{
			ID:          id,
			DailyEnergy: energy,
		}, nil
	}
}