package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Aqouet/go1fl-4-sprint-final/internal/common"
)

// parseTraining parses a string like "3456,Ходьба,3h00m" to number of steps, activity type and duration.
// Returns errors in case of parsing errors or if expected limits of parameters are exceeded.
func parseTraining(data string) (int, string, time.Duration, error) {
	dataSl := strings.Split(data, ",")
	if len(dataSl) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("slice after string %q split = %v with length %d, expected length = 3: %w", data, dataSl, len(dataSl), common.ErrSliceLen)
	}

	steps, err := strconv.Atoi(dataSl[0])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("string representing steps = %q: %w", dataSl[0], common.ErrParseInt)
	}

	if steps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("steps number = %d,  expected non-negative and non-zero steps number: %w", steps, common.ErrParamLimitExceeded)
	}

	activity := strings.TrimSpace(dataSl[1])

	if len(activity) == 0 {
		return 0, "", time.Duration(0), fmt.Errorf("activity must be set: %w", common.ErrEmptyString)
	}

	duration, err := time.ParseDuration(dataSl[2])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("string representing duration = %q: %w", dataSl[2], common.ErrParseDuration)
	}

	if duration <= time.Duration(0) {
		return 0, "", time.Duration(0), fmt.Errorf("duration = %v, expected non-negative and non-zero duration: %w", duration, common.ErrParamLimitExceeded)
	}

	return steps, activity, duration, nil
}

// distance calculates distance using number of steps and height
func distance(steps int, height float64) float64 {
	stepLength := common.StepLengthCoefficient * height
	distanceM := stepLength * float64(steps)
	distanceKm := distanceM / common.MinKm
	return distanceKm
}

// meanSpeed takes number of steps, height and activity duration.
// Returns mean speed.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 {
		return 0.
	}
	d := distance(steps, height)
	hours := duration.Hours()

	if hours <= 0 {
		return 0.0
	}

	return d / hours
}

// TrainingInfo generates informative message about training using:
//   - data string ("3456,Ходьба,3h00m")
//   - user's weight and height
//
// Returns formatted string or error if parsing or calculation failed.
func TrainingInfo(data string, weight, height float64) (string, error) {

	var err error

	steps, activity, duration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	var calories float64

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: activity = %q, expected activity = ['Ходьба', 'Бег']: %w", activity, common.ErrParamLimitExceeded)
	}

	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	mSpeed := meanSpeed(steps, height, duration)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Тип тренировки: %s\n", activity))
	b.WriteString(fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours()))
	b.WriteString(fmt.Sprintf("Дистанция: %.2f км.\n", dist))
	b.WriteString(fmt.Sprintf("Скорость: %.2f км/ч\n", mSpeed))
	b.WriteString(fmt.Sprintf("Сожгли калорий: %.2f\n", calories))

	return b.String(), nil

}

// RunningSpentCalories calculates spent calories number for running activity.
// Returns an error if parameters exceeds allowable limits.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps number = %d,  expected non-negative and non-zero steps number: %w", steps, common.ErrParamLimitExceeded)
	}

	if weight <= 0. {
		return 0, fmt.Errorf("weight = %.2f,  expected non-negative and non-zero weight: %w", weight, common.ErrParamLimitExceeded)
	}

	if height <= 0 {
		return 0, fmt.Errorf("height = %.2f,  expected non-negative and non-zero height: %w", height, common.ErrParamLimitExceeded)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("duration = %v,  expected non-negative and non-zero duration: %w", duration, common.ErrParamLimitExceeded)
	}

	mSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * mSpeed * durationInMinutes) / common.MinInH, nil
}

// WalkingSpentCalories calculates spent calories number for walking activity.
// Returns an error if parameters exceeds allowable limits.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	c, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return c * common.WalkingCaloriesCoefficient, nil

}
