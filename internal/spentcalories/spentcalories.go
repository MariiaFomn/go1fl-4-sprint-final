package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбивает строку с данными тренировки. Возвращаются шаги, вид и продолжительность активности.
func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("expected format 'steps, type of training, duration', got %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, "", 0, fmt.Errorf("error converting string to integer, %v", err)
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("invalid value: steps must be a positive integer, got %d", steps)
	}

	duration, err := time.ParseDuration(parts[2])

	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %v", err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("invalid duration: must be a positive value, got %v", duration)
	}
	return steps, parts[1], duration, nil
}

// distance рассчитывет пройденную дистанцию в зависимости от роста пользователя.
func distance(steps int, height float64) float64 {

	if height <= 0 {
		return 0
	}

	stepLength := stepLengthCoefficient * height

	distanceVal := (float64(steps) * stepLength) / mInKm

	return distanceVal
}

// meanSpeed рассчитывает среднюю скорость в зависимости от пройденной дистанции, шагов и роста пользователя.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	meanSpeed := distance / duration.Hours()

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)

	if err != nil {
		log.Printf("Failed to parse training datas: %v", err)
		return "", fmt.Errorf("Failed to parse training datas: %v", err)
	}

	var result string
	var distanceVal, speed, spentCalories float64

	switch activity {
	case "Бег":
		distanceVal = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("invalid calculation running spent calories: %w", err)
		}
	case "Ходьба":
		distanceVal = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)

		if err != nil {
			return "", fmt.Errorf("invalid calculation walking spent calories: %w", err)
		}

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceVal, speed, spentCalories)

	return result, nil
}

// RunningSpentCalories рассчитывает потраченные каллориии в зависимoсти от количества шагов, веса и роста пользователя
// и продолжительности тренировки.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("invalid value: steps must be a positive integer, got %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: must be a positive value, got %2.f", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("invalid height: must be a positive value, got %2.f", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: must be a positive value, got %v", duration)
	}

	meanSpeed := meanSpeed(steps, height, duration)

	runningSpentCalories := (weight * meanSpeed * duration.Minutes()) / minInH

	if runningSpentCalories < 0 {
		return 0, fmt.Errorf("calculation error: negative spentCalories %.2f", runningSpentCalories)
	}

	return runningSpentCalories, nil
}

// WalkingSpentCalories рассчитывает потраченные каллориии в зависимoсти от количества шагов, веса и роста пользователя
// и продолжительности тренировки.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("invalid value: steps must be a positive integer, got %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: must be a positive value, got %2.f", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("invalid height: must be a positive value, got %2.f", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: must be a positive value, got %v", duration)
	}

	meanSpeed := meanSpeed(steps, height, duration)

	spentCalories := (weight * meanSpeed * duration.Minutes()) / minInH

	if spentCalories < 0 {
		return 0, fmt.Errorf("calculation error: negative spentCalories %.2f", spentCalories)
	}

	walkingSpentCalories := spentCalories * walkingCaloriesCoefficient

	return walkingSpentCalories, nil
}
