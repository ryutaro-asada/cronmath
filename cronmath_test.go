package cronmath

import (
	"testing"
	"time"
)

func TestParseCron(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid cron", "5 9 * * *", false},
		{"invalid fields", "5 9 *", true},
		{"all wildcards", "* * * * *", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCron(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCron() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCronTime_Sub(t *testing.T) {
	tests := []struct {
		name     string
		cronStr  string
		duration time.Duration
		want     string
	}{
		{
			name:     "subtract 5 minutes from 9:05",
			cronStr:  "5 9 * * *",
			duration: Minutes(5),
			want:     "0 9 * * *",
		},
		{
			name:     "subtract 30 minutes from 10:15",
			cronStr:  "15 10 * * *",
			duration: Minutes(30),
			want:     "45 9 * * *",
		},
		{
			name:     "subtract with hour rollback",
			cronStr:  "0 10 * * *",
			duration: Minutes(30),
			want:     "30 9 * * *",
		},
		{
			name:     "subtract to previous day",
			cronStr:  "30 0 * * *",
			duration: Minutes(60),
			want:     "30 23 * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, err := ParseCron(tt.cronStr)
			if err != nil {
				t.Fatalf("ParseCron() error = %v", err)
			}

			err = cron.Sub(tt.duration)
			if err != nil {
				t.Fatalf("Sub() error = %v", err)
			}

			if got := cron.String(); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronTime_Add(t *testing.T) {
	tests := []struct {
		name     string
		cronStr  string
		duration time.Duration
		want     string
	}{
		{
			name:     "add 5 minutes to 9:00",
			cronStr:  "0 9 * * *",
			duration: Minutes(5),
			want:     "5 9 * * *",
		},
		{
			name:     "add 45 minutes to 10:30",
			cronStr:  "30 10 * * *",
			duration: Minutes(45),
			want:     "15 11 * * *",
		},
		{
			name:     "add with hour rollover",
			cronStr:  "45 9 * * *",
			duration: Minutes(30),
			want:     "15 10 * * *",
		},
		{
			name:     "add to next day",
			cronStr:  "30 23 * * *",
			duration: Minutes(60),
			want:     "30 0 * * *",
		},
		{
			name:     "add 2 hours",
			cronStr:  "15 10 * * *",
			duration: Hours(2),
			want:     "15 12 * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, err := ParseCron(tt.cronStr)
			if err != nil {
				t.Fatalf("ParseCron() error = %v", err)
			}

			err = cron.Add(tt.duration)
			if err != nil {
				t.Fatalf("Add() error = %v", err)
			}

			if got := cron.String(); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronMath_FluentInterface(t *testing.T) {
	tests := []struct {
		name       string
		cronStr    string
		operations func(*CronMath) *CronMath
		want       string
	}{
		{
			name:    "subtract 5 minutes",
			cronStr: "5 9 * * *",
			operations: func(cm *CronMath) *CronMath {
				return cm.Sub(Minutes(5))
			},
			want: "0 9 * * *",
		},
		{
			name:    "add then subtract",
			cronStr: "30 10 * * *",
			operations: func(cm *CronMath) *CronMath {
				return cm.Add(Minutes(30)).Sub(Minutes(15))
			},
			want: "45 10 * * *",
		},
		{
			name:    "complex operations",
			cronStr: "0 12 * * *",
			operations: func(cm *CronMath) *CronMath {
				return cm.Add(Hours(2)).Add(Minutes(30)).Sub(Minutes(15))
			},
			want: "15 14 * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := New(tt.cronStr)
			result := tt.operations(cm)

			if result.Error() != nil {
				t.Fatalf("operations error = %v", result.Error())
			}

			if got := result.String(); got != tt.want {
				t.Errorf("operations result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronTime_Wildcards(t *testing.T) {
	cron, _ := ParseCron("* 9 * * *")
	err := cron.Sub(Minutes(5))

	if err == nil {
		t.Errorf("Expected error when adjusting wildcards, got nil")
	}
}
