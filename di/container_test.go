package di_test

import (
	"errors"
	"testing"

	"kpo1/di"
)

func TestResolveDependencyChain(t *testing.T) {
	container := di.NewContainer()
	calls := 0

	err := container.Register("number",
		func(c *di.Container) (any, error) {
			calls++
			return 7, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = container.Register("double",
		func(c *di.Container) (any, error) {
			value, err := c.Resolve("number")
			if err != nil {
				return nil, err
			}

			number, ok := value.(int)
			if !ok {
				return nil, errors.New("неверный тип числа")
			}

			return number * 2, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if calls != 0 {
		t.Fatal("регистрация не должна вызывать фабрику")
	}

	for i := 0; i < 2; i++ {
		value, err := container.Resolve("double")
		if err != nil {
			t.Fatal(err)
		}

		if value != 14 {
			t.Fatalf("ожидали 14, получили %v", value)
		}
	}

	if calls != 2 {
		t.Fatalf("ожидали 2 вызова фабрики, получили %d", calls)
	}
}

func TestContainerRejectsNilFactory(t *testing.T) {
	container := di.NewContainer()

	if err := container.Register("broken", nil); err == nil {
		t.Fatal("ожидали ошибку для nil-фабрики")
	}
}

func TestContainerUnknownDependency(t *testing.T) {
	container := di.NewContainer()

	if _, err := container.Resolve("missing"); err == nil {
		t.Fatal("ожидали ошибку неизвестной зависимости")
	}
}

func TestDuplicateRegistrationPreservesFactory(t *testing.T) {
	container := di.NewContainer()

	err := container.Register("number",
		func(c *di.Container) (any, error) {
			return 1, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = container.Register("number",
		func(c *di.Container) (any, error) {
			return 2, nil
		},
	)
	if err == nil {
		t.Fatal("ожидали ошибку повторной регистрации")
	}

	value, err := container.Resolve("number")
	if err != nil {
		t.Fatal(err)
	}

	if value != 1 {
		t.Fatal("повторная регистрация заменила исходную фабрику")
	}
}

func TestContainerReturnsFactoryError(t *testing.T) {
	container := di.NewContainer()
	expected := errors.New("не удалось создать объект")

	err := container.Register("broken",
		func(c *di.Container) (any, error) {
			return nil, expected
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = container.Resolve("broken")

	if !errors.Is(err, expected) {
		t.Fatalf("ожидали ошибку фабрики, получили: %v", err)
	}
}
