package spentcalories

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                            = 0.65  // средняя длина шага.
	mInKm                              = 1000  // количество метров в километре.
	minInH                             = 60    // количество минут в часе.
	kmhInMsec                          = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM                              = 100   // количество сантиметров в метре.
	runningCaloriesMeanSpeedMultiplier = 18.0  // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0  // среднее количество сжигаемых калорий при беге.
	walkingCaloriesWeightMultiplier    = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier       = 0.029 // множитель роста.
)

// parseTraining парсит строку с данными о тренировке.
// Формат строки: "3456,Ходьба,3h00m", где 3456 — количество шагов, "Ходьба" — вид активности, "3h00m" — продолжительность.
// Возвращает количество шагов, вид активности, продолжительность и ошибку (если есть).
func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	// Парсим количество шагов
	steps, err := parseInt(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	// Получаем вид активности (убираем лишние пробелы)
	activity := strings.TrimSpace(parts[1])

	// Парсим продолжительность активности
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	// Возвращаем количество шагов, вид активности, продолжительность и nil (ошибки нет)
	return steps, activity, duration, nil
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

// distance возвращает дистанцию (в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// Дистанция = количество шагов * длина шага / 1000 (перевод в километры)
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	// Если продолжительность равна 0, возвращаем 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps)

	// Переводим продолжительность в часы
	hours := duration.Hours()

	// Средняя скорость = дистанция / продолжительность в часах
	return dist / hours
}

// RunningSpentCalories возвращает количество потраченных калорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, duration)

	// Формула для расчета калорий: (множитель * скорость - сдвиг) * вес
	return (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight
}

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, duration)

	// Переводим продолжительность в часы
	hours := duration.Hours()

	// Формула для расчета калорий: (множитель веса * вес + (скорость^2 / рост) * множитель роста) * продолжительность в часах * 60
	return (walkingCaloriesWeightMultiplier*weight + (speed*speed/height)*walkingSpeedHeightMultiplier) * hours * minInH
}

// TrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// Парсим данные о тренировке
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "неизвестный тип тренировки"
	}

	var dist, speed, calories float64

	// В зависимости от вида активности вычисляем дистанцию, скорость и калории
	switch activity {
	case "Бег":
		dist = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = RunningSpentCalories(steps, weight, duration)
	case "Ходьба":
		dist = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "неизвестный тип тренировки"
	}

	// Формируем и возвращаем строку с результатами
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, duration.Hours(), dist, speed, calories)
}
