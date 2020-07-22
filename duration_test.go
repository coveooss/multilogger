package multilogger

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func ExampleGetDurationFunc() {
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

	fmt.Fprintln(w, "Native\tClassic\tClassic Rounded\tPrecise\tPrecise Rounded\tNative long\tRounded long\t.")
	fmt.Fprintln(w, "------\t-------\t---------------\t-------\t---------------\t-----------\t------------\t.")

	print := func(duration time.Duration) {
		fmt.Fprintf(w, "%v\t", duration)
		fmt.Fprintf(w, "%s\t", GetDurationFunc(ClassicFormat, false, false)(duration))
		fmt.Fprintf(w, "%s\t", GetDurationFunc(ClassicFormat, true, false)(duration))
		fmt.Fprintf(w, "%s\t", GetDurationFunc(PreciseFormat, false, false)(duration))
		fmt.Fprintf(w, "%s\t", GetDurationFunc(PreciseFormat, true, false)(duration))
		fmt.Fprintf(w, "%s\t", GetDurationFunc(NativeFormat, false, true)(duration))
		fmt.Fprintf(w, "%s\t", GetDurationFunc(PreciseFormat, true, true)(duration))
		fmt.Fprintln(w, ".")
	}

	var total time.Duration
	for _, duration := range durations {
		total += 5 * duration
		print(duration)
		if duration != 0 {
			print(10 * duration)
			if total != duration {
				print(total)
			}
		}
	}
	w.Flush()

	// Output:
	// Native                 Classic                       Classic Rounded     Precise                          Precise Rounded     Native long                                 Rounded long                      .
	// ------                 -------                       ---------------     -------                          ---------------     -----------                                 ------------                      .
	// 0s                     0s                            0s                  0s                               0s                  0 second                                    0 second                          .
	// 1ns                    1ns                           1ns                 1ns                              1ns                 1 nanosecond                                1 nanosecond                      .
	// 10ns                   10ns                          10ns                10ns                             10ns                10 nanoseconds                              10 nanoseconds                    .
	// 5ns                    5ns                           5ns                 5ns                              5ns                 5 nanoseconds                               5 nanoseconds                     .
	// 1µs                    1µs                           1µs                 1µs                              1µs                 1 microsecond                               1 microsecond                     .
	// 10µs                   10µs                          10µs                10µs                             10µs                10 microseconds                             10 microseconds                   .
	// 5.005µs                5.005µs                       5.01µs              5µs5ns                           5µs10ns             5.005 microseconds                          5 microseconds 10 nanoseconds     .
	// 1ms                    1ms                           1ms                 1ms                              1ms                 1 millisecond                               1 millisecond                     .
	// 10ms                   10ms                          10ms                10ms                             10ms                10 milliseconds                             10 milliseconds                   .
	// 5.005005ms             5.005005ms                    5.01ms              5ms5µs5ns                        5ms10µs             5.005005 milliseconds                       5 milliseconds 10 microseconds    .
	// 1s                     1s                            1s                  1s                               1s                  1 second                                    1 second                          .
	// 10s                    10s                           10s                 10s                              10s                 10 seconds                                  10 seconds                        .
	// 5.005005005s           5.005005005s                  5.01s               5s5ms5µs5ns                      5s10ms              5.005005005 seconds                         5 seconds 10 milliseconds         .
	// 1m0s                   1m                            1m                  1m                               1m                  1 minute                                    1 minute                          .
	// 10m0s                  10m                           10m                 10m                              10m                 10 minutes                                  10 minutes                        .
	// 5m5.005005005s         5m5.005005005s                5m5s                5m5s5ms5µs5ns                    5m5s                5 minutes 5.005005005 seconds               5 minutes 5 seconds               .
	// 1h0m0s                 1h                            1h                  1h                               1h                  1 hour                                      1 hour                            .
	// 10h0m0s                10h                           10h                 10h                              10h                 10 hours                                    10 hours                          .
	// 5h5m5.005005005s       5h5m5.005005005s              5h5m                5h5m5s5ms5µs5ns                  5h5m                5 hours 5 minutes 5.005005005 seconds       5 hours 5 minutes                 .
	// 24h0m0s                1d                            1d                  1d                               1d                  24 hours                                    1 day                             .
	// 240h0m0s               1w3d                          1w3d                1w3d                             1w3d                240 hours                                   1 week 3 days                     .
	// 125h5m5.005005005s     5d5h5m5.005005005s            5d5h5m              5d5h5m5s5ms5µs5ns                5d5h5m              125 hours 5 minutes 5.005005005 seconds     5 days 5 hours 5 minutes          .
	// 168h0m0s               7d                            7d                  7d                               7d                  168 hours                                   7 days                            .
	// 1680h0m0s              2mo1w3d                       2mo1w3d             2mo1w3d                          2mo1w3d             1680 hours                                  2 months 1 week 3 days            .
	// 965h5m5.005005005s     5w5d5h5m5.005005005s          5w5d                5w5d5h5m5s5ms5µs5ns              5w5d                965 hours 5 minutes 5.005005005 seconds     5 weeks 5 days                    .
	// 720h0m0s               4w2d                          4w2d                4w2d                             4w2d                720 hours                                   4 weeks 2 days                    .
	// 7200h0m0s              10mo                          10mo                10mo                             10mo                7200 hours                                  10 months                         .
	// 4565h5m5.005005005s    6mo1w3d5h5m5.005005005s       6mo1w3d             6mo1w3d5h5m5s5ms5µs5ns           6mo1w3d             4565 hours 5 minutes 5.005005005 seconds    6 months 1 week 3 days            .
	// 8640h0m0s              1y                            1y                  1y                               1y                  8640 hours                                  1 year                            .
	// 86400h0m0s             10y                           10y                 10y                              10y                 86400 hours                                 10 years                          .
	// 47765h5m5.005005005s   5y6mo1w3d5h5m5.005005005s     5y6mo1w3d           5y6mo1w3d5h5m5s5ms5µs5ns         5y6mo1w3d           47765 hours 5 minutes 5.005005005 seconds   5 years 6 months 1 week 3 days    .
	// 83694h5m4.003002001s   9y8mo7d6h5m4.003002001s       9y8mo7d             9y8mo7d6h5m4s3ms2µs1ns           9y8mo7d             83694 hours 5 minutes 4.003002001 seconds   9 years 8 months 7 days           .
	// 836940h50m40.03002001s 96y10mo1w5d12h50m40.03002001s 96y10mo1w6d         96y10mo1w5d12h50m40s30ms20µs10ns 96y10mo1w6d         836940 hours 50 minutes 40.03002001 seconds 96 years 10 months 1 week 6 days  .
	// 466235h30m25.02001501s 53y11mo2w2d11h30m25.02001501s 53y11mo2w2d         53y11mo2w2d11h30m25s20ms15µs10ns 53y11mo2w2d         466235 hours 30 minutes 25.02001501 seconds 53 years 11 months 2 weeks 2 days .
}
