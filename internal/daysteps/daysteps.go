package daysteps

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Парсим количество шагов
	steps, err := parseInt(parts[0])
	if err != nil {
		return 0, 0, err
	}

	// Парсим продолжительность активности
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	// Возвращаем количество шагов, продолжительность и nil (ошибки нет)
	return steps, duration, nil
}

// parseInt преобразует строку в целое число (количество шагов).
func parseInt(s string) (int, error) {
	var steps int
	_, err := fmt.Sscanf(s, "%d", &steps)
	if err != nil {
		return 0, err
	}
	return steps, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	/// Парсим данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	// Проверяем чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в километрах
	distance := float64(steps) * StepLength / 1000

	// Вычисляем количество потраченных калорий
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	// Формируем и возвращаем строку с результатами
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distance, calories)

}
