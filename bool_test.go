package multilogger

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"True", true},
		{"  True  	\n", true},
		{"T", true},
		{"t", true},
		{"TrUe", true},
		{"On", true},
		{"Y", true},
		{"Yes", true},
		{"Whatever", true},
		{"False", false},
		{"FaLsE", false},
		{"f", false},
		{"off", false},
		{"No", false},
		{"n", false},
		{"NO", false},
		{"	 no  		", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBool(tt.name); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func ExampleFormatDurationPrecise() {
	durations := []time.Duration{
		0,
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		time.Minute,
		time.Hour,
		day,
		week,
		month,
		year,
		9*year + 8*month + 7*day + 6*time.Hour + 5*time.Minute + 4*time.Second + 3*time.Millisecond + 2*time.Microsecond + 1*time.Nanosecond,
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Native\tClassic\tPrecise\tRounded\tNative long\tRounded long\t")
	fmt.Fprintln(w, "------\t-------\t-------\t-------\t-----------\t------------\t")

	print := func(duration time.Duration) {
		fmt.Fprintf(w, "%v\t", duration)
		fmt.Fprintf(w, "%s\t", FormatDurationClassic(duration))
		fmt.Fprintf(w, "%s\t", FormatDurationPrecise(duration))
		fmt.Fprintf(w, "%s\t", FormatDurationRounded(duration))
		fmt.Fprintf(w, "%s\t", FormatDurationNativeLong(duration))
		fmt.Fprintf(w, "%s\t", FormatDurationRoundedLong(duration))
		fmt.Fprintln(w, "|")
	}

	var total time.Duration
	for _, duration := range durations {
		total += duration
		print(duration)
		if duration != 0 {
			print(10 * duration)
			if total != duration {
				print(total)
			}
		}
	}
	w.Flush()

	// // Output:
	// Native                 Classic                       Precise                          Rounded             Native long                                 Rounded long
	// ------                 -------                       -------                          -------             -----------                                 ------------
	// 0s                     0s                            0s                               0s                  0 second                                    0 second                         |
	// 1ns                    1ns                           1ns                              1ns                 1 nanosecond                                1 nanosecond                     |
	// 10ns                   10ns                          10ns                             10ns                10 nanoseconds                              10 nanoseconds                   |
	// 1µs                    1µs                           1µs                              1µs                 1 microsecond                               1 microsecond                    |
	// 10µs                   10µs                          10µs                             10µs                10 microseconds                             10 microseconds                  |
	// 1.001µs                1.001µs                       1µs1ns                           1µs                 1.001 microsecond                           1 microsecond                    |
	// 1ms                    1ms                           1ms                              1ms                 1 millisecond                               1 millisecond                    |
	// 10ms                   10ms                          10ms                             10ms                10 milliseconds                             10 milliseconds                  |
	// 1.001001ms             1.001001ms                    1ms1µs1ns                        1ms                 1.001001 millisecond                        1 millisecond                    |
	// 1s                     1s                            1s                               1s                  1 second                                    1 second                         |
	// 10s                    10s                           10s                              10s                 10 seconds                                  10 seconds                       |
	// 1.001001001s           1.001001001s                  1s1ms1µs1ns                      1s                  1.001001001 second                          1 second                         |
	// 1m0s                   1m                            1m                               1m                  1 minute                                    1 minute                         |
	// 10m0s                  10m                           10m                              10m                 10 minutes                                  10 minutes                       |
	// 1m1.001001001s         1m1.001001001s                1m1s1ms1µs1ns                    1m1s                1 minute 1.001001001 second                 1 minute 1 second                |
	// 1h0m0s                 1h                            1h                               1h                  1 hour                                      1 hour                           |
	// 10h0m0s                10h                           10h                              10h                 10 hours                                    10 hours                         |
	// 1h1m1.001001001s       1h1m1.001001001s              1h1m1s1ms1µs1ns                  1h1m                1 hour 1 minute 1.001001001 second          1 hour 1 minute                  |
	// 24h0m0s                1d                            1d                               1d                  24 hours                                    1 day                            |
	// 240h0m0s               1w3d                          1w3d                             1w3d                240 hours                                   1 week 3 days                    |
	// 25h1m1.001001001s      1d1h1m1.001001001s            1d1h1m1s1ms1µs1ns                1d1h1m              25 hours 1 minute 1.001001001 second        1 day 1 hour 1 minute            |
	// 168h0m0s               7d                            7d                               7d                  168 hours                                   7 days                           |
	// 1680h0m0s              2mo1w3d                       2mo1w3d                          2mo1w3d             1680 hours                                  2 months 1 week 3 days           |
	// 193h1m1.001001001s     8d1h1m1.001001001s            8d1h1m1s1ms1µs1ns                8d1h1m              193 hours 1 minute 1.001001001 second       8 days 1 hour 1 minute           |
	// 720h0m0s               4w2d                          4w2d                             4w2d                720 hours                                   4 weeks 2 days                   |
	// 7200h0m0s              10mo                          10mo                             10mo                7200 hours                                  10 months                        |
	// 913h1m1.001001001s     5w3d1h1m1.001001001s          5w3d1h1m1s1ms1µs1ns              5w3d                913 hours 1 minute 1.001001001 second       5 weeks 3 days                   |
	// 8640h0m0s              1y                            1y                               1y                  8640 hours                                  1 year                           |
	// 86400h0m0s             10y                           10y                              10y                 86400 hours                                 10 years                         |
	// 9553h1m1.001001001s    1y5w3d1h1m1.001001001s        1y5w3d1h1m1s1ms1µs1ns            1y5w3d              9553 hours 1 minute 1.001001001 second      1 year 5 weeks 3 days            |
	// 83694h5m4.003002001s   9y8mo7d6h5m4.003002001s       9y8mo7d6h5m4s3ms2µs1ns           9y8mo7d             83694 hours 5 minutes 4.003002001 seconds   9 years 8 months 7 days          |
	// 836940h50m40.03002001s 96y10mo1w5d12h50m40.03002001s 96y10mo1w5d12h50m40s30ms20µs10ns 96y10mo1w6d         836940 hours 50 minutes 40.03002001 seconds 96 years 10 months 1 week 6 days |
	// 93247h6m5.004003002s   10y9mo2w1d7h6m5.004003002s    10y9mo2w1d7h6m5s4ms3µs2ns        10y9mo2w1d          93247 hours 6 minutes 5.004003002 seconds   10 years 9 months 2 weeks 1 day  |
}
