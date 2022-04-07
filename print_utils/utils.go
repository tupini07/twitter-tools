package print_utils

import (
	"os"
	"time"

	pb "github.com/schollz/progressbar/v3"
)

func WaitUntilDay(day time.Time) {
	const timeFormat = "2006-01-02"
	targetDayStr := day.Format(timeFormat)

	for time.Now().Format(timeFormat) != targetDayStr {
		Printf("\rWaiting for %s, currently %s", targetDayStr, time.Now())
		time.Sleep(30 * time.Minute)
	}

	Fprint(os.Stdout, "\r \r")
}

func WaitWithBar(amount time.Duration, description string) {
	if amount <= 0 {
		Fatalf("Cannot wait 0 or less amount of time! Tried to wait %s", amount.String())
	}

	if isLoggingToFile {
		// if we're logging to file then don't show waiting bar since it would be very
		// messy. Just show that we're waiting and how long
		Printf("Waiting because '%s' for amount '%s'\n", description, amount)
		return
	}

	desc := "[cyan]Waiting ...[reset]"
	if description != "" {
		desc = description
	}

	bar := pb.NewOptions(
		int(amount.Seconds()),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionEnableColorCodes(true),
		pb.OptionSetWidth(15),
		pb.OptionSetDescription(desc),
		pb.OptionThrottle(100*time.Millisecond),
		pb.OptionSetTheme(pb.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[light_red][[reset]",
			BarEnd:        "[light_red]][reset]",
		}),
		pb.OptionSpinnerType(11),
		pb.OptionClearOnFinish(),
	)
	bar.RenderBlank()

	// We use target time instead of just counting seconds because the latter
	// won't be correct if PC running this goes to sleep or similar thing
	// happens. During normal operation this should behave the same as waiting 1
	// second at a time.
	targetTime := time.Now().Add(amount)

	for i := 0; i < int(amount.Seconds()); i++ {
		if time.Now().After(targetTime) {
			break
		}

		time.Sleep(time.Second)
		bar.Add(1)
	}

	bar.Finish()
}
