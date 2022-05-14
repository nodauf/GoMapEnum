package ntp

import (
	"time"

	"github.com/beevik/ntp"
)

// GetTime returns the time of the NTP server
func (options *Options) GetTime() {
	response, err := ntp.Query(options.Target)
	if err != nil {
		options.Log.Error("Error reading time: %v", err)
		return
	}
	options.Log.Debug("Offset: %s RTT: %s", response.ClockOffset, response.RTT)
	if options.UTC {
		options.Log.Success("Time on the target: %s", response.Time.Format("02/01/06 15:04:05 -0700"))
	} else {
		// Print the time in the same timezone as the client. To do so, apply the ClockOffset on time.Now()
		now := time.Now()
		options.Log.Success("Time on the target: %s", now.Add(response.ClockOffset).Format("02/01/06 15:04:05 -0700"))
	}
}
