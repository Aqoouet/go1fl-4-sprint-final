package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Aqouet/go1fl-4-sprint-final/internal/common"
	"github.com/Aqouet/go1fl-4-sprint-final/internal/spentcalories"
)

// parsePackage parses a string like "678,0h50m" to number of steps and duration.
// Returns errors in case of parsing errors or if expected limits of parameters are exceeded.
func parsePackage(data string) (int, time.Duration, error) {

	dataSl := strings.Split(data, ",")
	if len(dataSl) != 2 {
		return 0, time.Duration(0), fmt.Errorf("slice after string %q split = %v with length %d, expected length = 2: %w", data, dataSl, len(dataSl), common.ErrSliceLen)
	}

	steps, err := strconv.Atoi(dataSl[0])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("string representing steps = %q: %w", dataSl[0], common.ErrParseInt)
	}

	if steps <= 0 {
		return 0, time.Duration(0), fmt.Errorf("steps number = %d,  expected non-negative and non-zero steps number: %w", steps, common.ErrParamLimitExceeded)
	}

	duration, err := time.ParseDuration(dataSl[1])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("string representing duration = %q: %w", dataSl[1], common.ErrParseDuration)
	}

	if duration <= 0 {
		return 0, time.Duration(0), fmt.Errorf("duration = %v, expected non-negative and non-zero duration: %w", duration, common.ErrParamLimitExceeded)
	}

	return steps, duration, nil
}

// DayActionInfo parse string and using information about user weight and height returns informative message about:
//   - number od steps
//   - length of a trip
//   - calories
func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceM := float64(steps) * common.StepLength
	distanceKm := distanceM / common.MinKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)

}
