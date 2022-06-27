package twitter_api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/print_utils"
)

var clientInstance *twitter.Client

func getApiClient() *twitter.Client {
	if clientInstance == nil {
		appConf := app_config.GetConfig()

		config := oauth1.NewConfig(
			appConf.Auth.APIKey,
			appConf.Auth.APISecretKey)
		token := oauth1.NewToken(
			appConf.Auth.AccessToken,
			appConf.Auth.AccessTokenSecret)

		httpClient := config.Client(oauth1.NoContext, token)

		clientInstance = twitter.NewClient(httpClient)
	}

	return clientInstance
}

const minimumDelayBetweenRequests = 10 * time.Second

var lastRequestTime time.Time = time.Now().Add(-(time.Hour * 24)) // first time is always far away

// makes a request and ensures that subsequent requests all respect the "delay",
// which corresponds to the rate limit in Twitter's documentation
func makeTimeoutHandledRequest(delay time.Duration, reqFunc func() (interface{}, *http.Response, error)) interface{} {
	numErrors := 0
	const maxErrors = 7

	for {
		// if not enough time has elapsed since last request then ensure we wait
		effectiveDelay := minimumDelayBetweenRequests
		if delay != -1 {
			effectiveDelay = delay
		}

		waitingDifference := effectiveDelay - time.Since(lastRequestTime)
		if waitingDifference > 0 {
			print_utils.WaitWithBar(waitingDifference, "Avoid overloading Twitter API")
		}

		data, resp, err := reqFunc()
		lastRequestTime = time.Now()

		respLogger := log.WithField("err", err)
		if resp != nil {
			respLogger.WithField("status", resp.Status)
		}
		respLogger.Debug("Got response from server")

		if resp != nil {
			if resp.StatusCode == 429 ||
				resp.StatusCode == 88 ||
				resp.StatusCode == 420 {
				print_utils.WaitWithBar(15*time.Minute, "[yellow]Waiting timeout from Twitter[reset]")
				continue
			}

			if resp.StatusCode == 161 {
				logger_str, _ := respLogger.String()
				print_utils.Fatalf("Getting a warning from Twitter that we can no longer follow users! Stopping execution since proceeding may cause "+
					"our account to be blocked. Please update the 'max_total_following' value so that we no longer try to follow more users than what Twitter allows. A good "+
					"idea is to set this value to your current number of followers or less.", logger_str)
			}

			if resp.StatusCode == 326 {
				logger_str, _ := respLogger.String()
				print_utils.Fatalf("Our account has been blocked! Stopping execution. Please log in to https://twitter.com to unlock your account. "+
					"Wait some days before running this tool again to avoid angering the 'twitter gods' even more. Data: %s\n", logger_str)
			}
		}

		if err != nil || resp == nil {
			// TODO for some reason 161 errors have an empty resp so the handler above is not working.. An issue with the upstream dep maybe?
			// so we'll handle it in a ugly way here.
			if strings.Contains(err.Error(), "161 You are unable to follow more people at this time") {
				//! NOTE remove this! This is duplicate code, same as the `if resp.StatusCode == 161` handler above
				logger_str, _ := respLogger.String()
				print_utils.Fatalf("[FROM DUPLICATED SECTION] Getting a warning from Twitter that we can no longer follow users! Stopping execution since proceeding may cause "+
					"our account to be blocked. Please update the 'max_total_following' value so that we no longer try to follow more users than what Twitter allows. A good "+
					"idea is to set this value to your current number of followers or less.", logger_str)
			}

			numErrors += 1

			print_utils.Errorf("Retrying since there was an error: %s", err)
			print_utils.Errorf("Total amount of errors encountered: %d", numErrors)

			if numErrors < maxErrors {
				print_utils.WaitWithBar(time.Minute, "Waiting because of error")
				continue
			} else {
				print_utils.Fatal("Stopping because we had too many errors")
			}
		}

		return data
	}
}

var (
	cyan       = color.New(color.FgCyan).SprintFunc()
	green      = color.New(color.FgGreen, color.Underline).SprintFunc()
	red        = color.New(color.FgRed).SprintFunc()
	strike     = color.New(color.CrossedOut).SprintFunc()
	underline  = color.New(color.Underline).SprintFunc()
	white      = color.New(color.FgWhite).SprintFunc()
	white_bold = color.New(color.FgWhite, color.Bold).SprintFunc()
	yellow     = color.New(color.FgYellow).SprintFunc()
)

func printTitle(title string) {
	decorated := fmt.Sprintf("%s %s %s\n",
		strike(yellow("###")),
		underline(white_bold(title)),
		strike(yellow("###")))

	decoratedLen := 0
	for _, c := range decorated {
		if unicode.Is(unicode.Latin, c) {
			decoratedLen++
		}
	}

	wrapperLine := strike(yellow(strings.Repeat("#", decoratedLen))) + "\n"

	// #################
	// ### some text ###
	// #################

	print_utils.Printf("\n%s%s%s\n", wrapperLine, decorated, wrapperLine)
}

func printAction(action string) {
	print_utils.Printf("\t%s%s\n",
		yellow("•"),
		cyan(action))
}

func printStepAction(currentStep, totalSteps int, action string) {
	printAction("" +
		red("[") + white(currentStep) + red("/") +
		white(totalSteps) + red("] ") +
		white(action))
}

func printActionLog(action string) {
	print_utils.Printf("\t%s\n", yellow(action))
}
