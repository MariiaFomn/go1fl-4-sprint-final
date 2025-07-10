package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	spentCalories "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разделяет и конвертирует данные строки в шаги и продолжительность прогулки.
func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected format 'steps, duration', got %d", len(parts))
	}

	stepsStr, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, 0, fmt.Errorf("error converting string to integer, %w", err)
	}

	if stepsStr <= 0 {
		return 0, 0, errors.New("steps must be a positive integer")
	}

	durationStr, err := time.ParseDuration(parts[1])

	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if durationStr <= 0 {
		return 0, 0, errors.New("durationStr must be a positive integer")
	}

	return stepsStr, durationStr, nil
}

// DayActionInfo парсит строку с данными с помощью parsePackage(), вычисляет дистанцию и количество потраченных каллорий.
func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Printf("invalid parsing operation: %v", err)
		return ""
	}

	if steps <= 0 {
		err := fmt.Errorf("steps must be a positive integer")
		log.Println(err)
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm

	walkingSpentCalories, err := spentCalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Printf("invalid calculations walking spent calories: %v", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, walkingSpentCalories)
}
