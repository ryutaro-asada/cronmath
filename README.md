# CronMath

[![Go Reference](https://pkg.go.dev/badge/github.com/ryutaro-asada/cronmath.svg)](https://pkg.go.dev/github.com/ryutaro-asada/cronmath)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryutaro-asada/cronmath)](https://goreportcard.com/report/github.com/ryutaro-asada/cronmath)
[![CI](https://github.com/ryutaro-asada/cronmath/actions/workflows/ci.yml/badge.svg)](https://github.com/ryutaro-asada/cronmath/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/ryutaro-asada/cronmath/badge.svg?branch=main)](https://coveralls.io/github/ryutaro-asada/cronmath?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ryutaro-asada/cronmath)](https://go.dev/)
[![Release](https://img.shields.io/github/release/ryutaro-asada/cronmath.svg)](https://github.com/ryutaro-asada/cronmath/releases/latest)

A Go library for performing time arithmetic operations on cron expressions. Add or subtract minutes and hours from cron expressions with a simple, fluent API.

## ✨ Features

- 🕒 **Time Arithmetic** - Add or subtract minutes and hours from cron expressions
- 🔄 **Day Boundary Handling** - Automatically handles transitions across midnight
- 🔗 **Fluent Interface** - Chain multiple operations for complex adjustments
- ⚡ **Zero Dependencies** - Pure Go implementation with no external dependencies
- 🛡️ **Type Safe** - Compile-time type checking with Go's type system
- 📊 **Well Tested** - Comprehensive test coverage with edge cases
- 🚀 **Performance** - Efficient operations with minimal allocations

## 📦 Installation

```bash
go get github.com/ryutaro-asada/cronmath
```

Requirements:
- Go 1.21 or higher

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/ryutaro-asada/cronmath"
)

func main() {
    // Simple operation - subtract 5 minutes from 9:05 AM
    result := cronmath.New("5 9 * * *").Sub(cronmath.Minutes(5))
    fmt.Println(result.String()) // Output: 0 9 * * *
    
    // Chain multiple operations
    result = cronmath.New("30 10 * * *").
        Add(cronmath.Hours(2)).
        Sub(cronmath.Minutes(15))
    fmt.Println(result.String()) // Output: 15 12 * * *
    
    // Error handling
    result, err := cronmath.New("0 14 * * *").
        Sub(cronmath.Hours(3)).
        Result()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // Output: 0 11 * * *
}
```

## 📖 Documentation

### Core Types

#### `CronTime`
Represents a parsed cron expression that can be manipulated.

```go
cron, err := cronmath.ParseCron("0 9 * * *")
if err != nil {
    log.Fatal(err)
}
cron.Add(cronmath.Minutes(30))
fmt.Println(cron.String()) // "30 9 * * *"
```

#### `CronMath`
Provides a fluent interface for cron arithmetic operations.

```go
cm := cronmath.New("0 9 * * *")
result := cm.Add(cronmath.Hours(1)).Sub(cronmath.Minutes(15))
```

### Functions

| Function | Description | Example |
|----------|-------------|---------|
| `New(cronStr string) *CronMath` | Creates a new CronMath instance | `cronmath.New("0 9 * * *")` |
| `ParseCron(cronStr string) (*CronTime, error)` | Parses a cron expression | `cronmath.ParseCron("0 9 * * *")` |
| `Minutes(n int) Duration` | Creates a duration of n minutes | `cronmath.Minutes(30)` |
| `Hours(n int) Duration` | Creates a duration of n hours | `cronmath.Hours(2)` |

### Methods

| Method | Description | Returns |
|--------|-------------|---------|
| `Add(d Duration) *CronMath` | Adds duration to the cron expression | `*CronMath` for chaining |
| `Sub(d Duration) *CronMath` | Subtracts duration from the cron expression | `*CronMath` for chaining |
| `String() string` | Returns the cron expression as string | Current cron expression |
| `Error() error` | Returns any error that occurred | Error or nil |
| `Result() (string, error)` | Returns the cron expression and error | Expression and error |

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

### Integration with Terraform Provider

Example of using cronmath in a Terraform provider:

```go
package provider

import (
    "github.com/ryutaro-asada/cronmath"
    "github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *CronScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan CronScheduleModel
    
    // Get the base cron expression
    baseCron := plan.BaseCron.ValueString()
    
    // Apply adjustments
    cm := cronmath.New(baseCron)
    for _, adjustment := range plan.Adjustments {
        minutes := adjustment.Minutes.ValueInt64()
        if adjustment.Type.ValueString() == "add" {
            cm = cm.Add(cronmath.Minutes(int(minutes)))
        } else {
            cm = cm.Sub(cronmath.Minutes(int(minutes)))
        }
    }
    
    // Get the final result
    finalCron, err := cm.Result()
    if err != nil {
        resp.Diagnostics.AddError("Cron Calculation Error", err.Error())
        return
    }
    
    plan.FinalCron = types.StringValue(finalCron)
}
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

Run tests with coverage:
```bash
go test -race -cover ./...
```

Run tests with detailed coverage report:
```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Write clear, idiomatic Go code
- Add tests for new functionality
- Maintain backward compatibility
- Update documentation as needed
- Follow Go naming conventions
- Keep functions small and focused

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for simple cron expression manipulation in Go
- Designed for use with Terraform providers and CI/CD pipelines
- Thanks to all contributors and users of this library

## 📮 Support

- 📧 For questions and support, please [open an issue](https://github.com/ryutaro-asada/cronmath/issues)
- 🐛 For bug reports, please include a minimal reproducible example
- 💡 For feature requests, please describe your use case

## 🔗 Links

- [Documentation](https://pkg.go.dev/github.com/ryutaro-asada/cronmath)
- [Examples](https://github.com/ryutaro-asada/cronmath/tree/main/examples)
- [Releases](https://github.com/ryutaro-asada/cronmath/releases)
- [Contributing Guide](CONTRIBUTING.md)

---

Made with ❤️ by [ryutaro-asada](https://github.com/ryutaro-asada)
