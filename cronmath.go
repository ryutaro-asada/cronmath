package cronmath

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronTime represents a cron expression that can be manipulated
type CronTime struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// ParseCron parses a cron expression string into a CronTime struct
func ParseCron(cronStr string) (*CronTime, error) {
	parts := strings.Fields(cronStr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}

	return &CronTime{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
	}, nil
}

// String returns the cron expression as a string
func (c *CronTime) String() string {
	return fmt.Sprintf("%s %s %s %s %s", c.Minute, c.Hour, c.DayOfMonth, c.Month, c.DayOfWeek)
}

// Add adds a duration to the cron expression
func (c *CronTime) Add(d time.Duration) error {
	return c.adjustTime(d)
}

// Sub subtracts a duration from the cron expression
func (c *CronTime) Sub(d time.Duration) error {
	return c.adjustTime(-d)
}

// adjustTime adjusts the cron time by the given duration
func (c *CronTime) adjustTime(d time.Duration) error {
	// Only handle minute and hour adjustments for now
	// More complex adjustments (days, months) would require more sophisticated logic

	totalMinutes := int(d.Minutes())

	// Parse current minute and hour
	currentMinute, err := c.parseField(c.Minute, 0, 59)
	if err != nil {
		return fmt.Errorf("error parsing minute: %v", err)
	}

	currentHour, err := c.parseField(c.Hour, 0, 23)
	if err != nil {
		return fmt.Errorf("error parsing hour: %v", err)
	}

	// Skip if wildcards
	if currentMinute == -1 || currentHour == -1 {
		return fmt.Errorf("cannot adjust wildcards")
	}

	// Calculate new time
	totalCurrentMinutes := currentHour*60 + currentMinute
	newTotalMinutes := totalCurrentMinutes + totalMinutes

	// Handle overflow/underflow for daily schedule
	if newTotalMinutes < 0 {
		// Go to previous day
		newTotalMinutes += 24 * 60
	} else if newTotalMinutes >= 24*60 {
		// Go to next day
		newTotalMinutes -= 24 * 60
	}

	newHour := newTotalMinutes / 60
	newMinute := newTotalMinutes % 60

	c.Minute = strconv.Itoa(newMinute)
	c.Hour = strconv.Itoa(newHour)

	return nil
}

// parseField parses a cron field value
func (c *CronTime) parseField(field string, min, max int) (int, error) {
	if field == "*" {
		return -1, nil // Wildcard
	}

	// Handle simple numeric values
	val, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("unsupported field format: %s", field)
	}

	if val < min || val > max {
		return 0, fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
	}

	return val, nil
}

// Duration represents a time duration for cron operations
type Duration = time.Duration

// Minutes creates a duration of n minutes
func Minutes(n int) Duration {
	return time.Duration(n) * time.Minute
}

// Hours creates a duration of n hours
func Hours(n int) Duration {
	return time.Duration(n) * time.Hour
}

// CronMath provides a fluent interface for cron arithmetic
type CronMath struct {
	cron *CronTime
	err  error
}

// New creates a new CronMath instance from a cron string
func New(cronStr string) *CronMath {
	c, err := ParseCron(cronStr)
	return &CronMath{cron: c, err: err}
}

// Add adds duration to the cron expression
func (cm *CronMath) Add(d Duration) *CronMath {
	if cm.err != nil {
		return cm
	}
	cm.err = cm.cron.Add(d)
	return cm
}

// Sub subtracts duration from the cron expression
func (cm *CronMath) Sub(d Duration) *CronMath {
	if cm.err != nil {
		return cm
	}
	cm.err = cm.cron.Sub(d)
	return cm
}

// String returns the resulting cron expression
func (cm *CronMath) String() string {
	if cm.err != nil {
		return fmt.Sprintf("error: %v", cm.err)
	}
	return cm.cron.String()
}

// Error returns any error that occurred during operations
func (cm *CronMath) Error() error {
	return cm.err
}
