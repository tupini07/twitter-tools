package print_utils

import (
	"log"
	"os"
	"time"

	pb "github.com/schollz/progressbar/v3"
)

func WaitWithBar(amount time.Duration, description string) {
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

	if amount <= 0 {
		log.Fatalf("Cannot wait 0 or less amount of time! Tried to wait %s", amount.String())
	}

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
