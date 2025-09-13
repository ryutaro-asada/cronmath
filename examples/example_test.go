package cronmath_test

import (
	"fmt"

	"github.com/ryutaro-asada/cronmath"
)

func ExampleNew() {
	cm := cronmath.New("5 9 * * *")
	fmt.Println(cm.String())
	// Output: 5 9 * * *
}

func ExampleCronMath_Sub() {
	result := cronmath.New("5 9 * * *").Sub(cronmath.Minutes(5))
	fmt.Println(result.String())
	// Output: 0 9 * * *
}

func ExampleCronMath_Add() {
	result := cronmath.New("30 10 * * *").Add(cronmath.Hours(2))
	fmt.Println(result.String())
	// Output: 30 12 * * *
}

func ExampleCronMath_chainedOperations() {
	result := cronmath.New("0 10 * * *").
		Add(cronmath.Hours(2)).
		Add(cronmath.Minutes(30)).
		Sub(cronmath.Minutes(15))
	fmt.Println(result.String())
	// Output: 15 12 * * *
}

func ExampleCronMath_dayBoundary() {
	// Crossing midnight backward
	result := cronmath.New("30 0 * * *").Sub(cronmath.Hours(1))
	fmt.Println(result.String())
	// Output: 30 23 * * *
}

func ExampleCronMath_Error() {
	cm := cronmath.New("invalid cron expression")
	if err := cm.Error(); err != nil {
		fmt.Println("Error:", err)
	}
	// Output: Error: invalid cron expression: expected 5 fields, got 3
}

func ExampleMinutes() {
	duration := cronmath.Minutes(30)
	result := cronmath.New("0 9 * * *").Add(duration)
	fmt.Println(result.String())
	// Output: 30 9 * * *
}

func ExampleHours() {
	duration := cronmath.Hours(3)
	result := cronmath.New("0 9 * * *").Add(duration)
	fmt.Println(result.String())
	// Output: 0 12 * * *
}
