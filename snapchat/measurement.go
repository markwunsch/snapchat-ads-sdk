package snapchat

import (
	"context"
	"fmt"
	"strings"
)

// MeasurementService provides functions for getting snapchat measurement metrics
type MeasurementService service

// GetMeasurementsResponse is the response object for get measurement requests
type GetMeasurementsResponse struct {
	RequestStatus string        `json:"request_status"` // the status of the request
	RequestId     string        `json:"request_id"`     // the id of the request
	TotalStats    []*TotalStats `json:"total_stats"`    // the object containing total stats for this entity
}

// TotalStats is a wrapper object for the stats response
type TotalStats struct {
	SubRequestStatus string    `json:"sub_request_status"`
	TotalStat        TotalStat `json:"total_stat"`
}

// TotalStat is an object that contains the total metrics for the requested entity
type TotalStat struct {
	Id          string           `json:"id"`          // the id of the entity
	Type        string           `json:"type"`        // the type of the entity
	Granularity string           `json:"granularity"` // the level of granulatiry for stats reporting
	Stats       MeasurementStats `json:"stats"`       // the object containing actual metrics
}

// MeasurementStats is the object actually containing the measurement metrics
type MeasurementStats struct {
	Impressions      int   `json:"impressions"`        // number of impressions
	Swipes           int   `json:"swipes"`             // number of swipe-ups
	Spend            int64 `json:"spend"`              // amount spent (micro-currency)
	FirstQuartile    int   `json:"quartile_1"`         // number of video views to 25%
	SecondQuartile   int   `json:"quartile_2"`         // number of video views to 50%
	ThirdQuartile    int   `json:"quartile_3"`         // number of video views to 75%
	ScreenTimeMillis int64 `json:"screen_time_millis"` // total time spent on top snap ad (milliseconds)
	ViewCompletion   int   `json:"view_completion"`    // the number of views to completion
	VideoViews       int   `json:"video_views"`        // the total number of impressions that meet the qualifying video view criteria of at least 2 seconds of consecutive watch time or a swipe up action on the Top Snap
}

// GetStatsForAdSquad returns the measurement metrics for the given ad squad
func (measurement *MeasurementService) GetStatsForAdSquad(ctx context.Context, adSquadId string) (*GetMeasurementsResponse, error) {
	path := fmt.Sprintf(`adsquads/%s/stats`, adSquadId)
	req, err := measurement.client.createRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	m := new(GetMeasurementsResponse)
	err = measurement.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(m.RequestStatus) == "success" {
		return m, nil
	}
	return nil, fmt.Errorf(`non-success status returned from snapchat api (get stats for ad squad with id %s): %s`, adSquadId, m.RequestStatus)
}
