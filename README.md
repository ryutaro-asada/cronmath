# CronMath

[![Go Reference](https://pkg.go.dev/badge/github.com/ryutaro-asada/cronmath.svg)](https://pkg.go.dev/github.com/ryutaro-asada/cronmath)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryutaro-asada/cronmath)](https://goreportcard.com/report/github.com/ryutaro-asada/cronmath)
[![CI](https://github.com/ryutaro-asada/cronmath/actions/workflows/ci.yml/badge.svg)](https://github.com/ryutaro-asada/cronmath/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/ryutaro-asada/cronmath/badge.svg?branch=main)](https://coveralls.io/github/ryutaro-asada/cronmath?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ryutaro-asada/cronmath)](https://go.dev/)
[![Release](https://img.shields.io/github/release/ryutaro-asada/cronmath.svg)](https://github.com/ryutaro-asada/cronmath/releases/latest)

A Go library for performing time arithmetic operations on cron expressions. Add or subtract minutes and hours from cron expressions with a simple, fluent API.

## 📦 Installation

```bash
go get github.com/ryutaro-asada/cronmath
```

## 📚 Examples

### Timezone Adjustment

Convert schedules between timezones:

```go
// Convert UTC to JST (UTC+9)
jstSchedule := cronmath.New("0 3 * * *").Add(cronmath.Hours(9))
fmt.Println(jstSchedule.String()) // "0 12 * * *"

// Convert EST to UTC (EST is UTC-5)
utcSchedule := cronmath.New("0 10 * * *").Add(cronmath.Hours(5))
fmt.Println(utcSchedule.String()) // "0 15 * * *"

// Convert PST to EST (PST is UTC-8, EST is UTC-5)
estSchedule := cronmath.New("0 9 * * *").Add(cronmath.Hours(3))
fmt.Println(estSchedule.String()) // "0 12 * * *"
```

### Staggered Job Scheduling

Create staggered schedules to avoid resource contention:

```go
baseTime := "0 2 * * *" // 2:00 AM

schedules := []struct {
    name   string
    offset cronmath.Duration
}{
    {"primary", 0},
    {"secondary", cronmath.Minutes(15)},
    {"tertiary", cronmath.Minutes(30)},
    {"quaternary", cronmath.Hours(1)},
}

for _, s := range schedules {
    schedule := cronmath.New(baseTime).Add(s.offset)
    fmt.Printf("%s: %s\n", s.name, schedule.String())
}
// Output:
// primary: 0 2 * * *
// secondary: 15 2 * * *
// tertiary: 30 2 * * *
// quaternary: 0 3 * * *
```

### Day Boundary Handling

The library automatically handles transitions across midnight:

```go
// Subtract time crossing midnight backward
result := cronmath.New("30 0 * * *").Sub(cronmath.Hours(1))
fmt.Println(result.String()) // "30 23 * * *"

// Add time crossing midnight forward
result = cronmath.New("30 23 * * *").Add(cronmath.Hours(2))
fmt.Println(result.String()) // "30 1 * * *"

// Complex boundary case
result = cronmath.New("45 23 * * *").
    Add(cronmath.Hours(3)).
    Add(cronmath.Minutes(30))
fmt.Println(result.String()) // "15 3 * * *"
```

### Error Handling Patterns

Different approaches for error handling:

```go
// Method 1: Check error at the end
cm := cronmath.New("invalid cron")
result := cm.Add(cronmath.Hours(1))
if err := result.Error(); err != nil {
    log.Printf("Error: %v", err)
}

// Method 2: Use Result() for combined return
result, err := cronmath.New("0 9 * * *").
    Add(cronmath.Hours(2)).
    Result()
if err != nil {
    return fmt.Errorf("failed to calculate cron: %w", err)
}

// Method 3: Direct CronTime manipulation
cron, err := cronmath.ParseCron("0 9 * * *")
if err != nil {
    return err
}
if err := cron.Add(cronmath.Hours(1)); err != nil {
    return err
}
```

## 🎯 Cron Expression Format

This library supports standard 5-field cron expressions:

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

### Valid Examples
- `0 9 * * *` - Every day at 9:00 AM
- `30 14 * * 1-5` - Weekdays at 2:30 PM
- `0 0 1 * *` - First day of every month at midnight
- `15 10 * * 6,0` - Weekends at 10:15 AM

## ⚠️ Limitations

- **Only minute and hour fields can be modified** - Day, month, and day of week fields remain unchanged
- **Wildcards in time fields cannot be adjusted** - Expressions like `* 9 * * *` will return an error
- **No support for complex expressions** - Ranges (`0-30`), lists (`15,30,45`), and steps (`*/5`) are not supported in minute/hour fields
- **No validation of day/month combinations** - The library doesn't validate if the resulting date is valid

## 🧪 Testing

Run tests:
```bash
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📮 Support

- 📧 For questions and support, please [open an issue](https://github.com/ryutaro-asada/cronmath/issues)
- 🐛 For bug reports, please include a minimal reproducible example
- 💡 For feature requests, please describe your use case

## 🔗 Links

- [Documentation](https://pkg.go.dev/github.com/ryutaro-asada/cronmath)
- [Examples](https://github.com/ryutaro-asada/cronmath/tree/main/examples)
- [Releases](https://github.com/ryutaro-asada/cronmath/releases)

