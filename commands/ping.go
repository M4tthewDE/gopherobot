package commands

import (
	"fmt"
	"time"

	"de.com.fdm/gopherobot/config"
	"github.com/hako/durafmt"
)

func Ping(startTime time.Time, latency time.Duration, config *config.Config) string {
	uptime := time.Since(startTime)

	result := "Uptime: " + durafmt.Parse(uptime).LimitFirstN(2).String() + " |"
	result += " Commit: " + config.Git.Commit + " |"
	result += " Branch: " + config.Git.Branch + " |"
	result += " Latency to tmi: " + fmt.Sprint(latency)

	return result
}
