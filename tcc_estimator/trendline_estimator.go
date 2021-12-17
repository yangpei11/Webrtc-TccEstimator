package tccEstimator

import (
	"math"
)

const(
	kDeltaCounterMax = 1000
	smoothing_coef = 0.9
	window_size = 20
	kMinNumDeltas = 60
	threshold_gain = 4.0
	kOverUsingTimeThreshold = 10.0
	kMaxAdaptOffsetMs = 15.0
	k_down = 0.039
	k_up = 0.0087
)

type PacketTiming struct{
	arrival_time_ms float64
	smoothed_delay_ms float64
	raw_delay_ms float64
}

type TrendlineEstimator struct {
	num_of_deltas_ int
	first_arrival_time_ms_ int64
	accumulated_delay_ float64
	smooth_delay_ float64
	delay_hist_ []PacketTiming
	prev_trend_ float64
	hypothesis_ BandwidthUsage
	threshold_ float64
	time_over_using_ float64
	overuse_counter_ int
	last_update_ms_ int64
	record_trend float64
}

func NewTrendlineEstimator()TrendlineEstimator{
	return TrendlineEstimator{
		num_of_deltas_: 0,
		first_arrival_time_ms_: -1,
		accumulated_delay_: 0,
		smooth_delay_: 0,
		delay_hist_: make([]PacketTiming, 0),
		prev_trend_: 0,
		hypothesis_: kBwNormal,
		threshold_: 12.5,
		time_over_using_: -1,
		overuse_counter_: 0,
		last_update_ms_: -1,
	}
}

func (trendline_estimator *TrendlineEstimator)State()BandwidthUsage{
	return trendline_estimator.hypothesis_
}

func LinearFitSlope(packets []PacketTiming) float64{
	sum_x := float64(0)
	sum_y := float64(0)
	for _, value := range packets{
		sum_x += value.arrival_time_ms
		sum_y += value.smoothed_delay_ms
	}
	x_avg := sum_x/float64(len(packets))
	y_avg := sum_y/float64(len(packets))

	numerator := float64(0)
	denominator := float64(0)
	for _, value := range packets{
		x := value.arrival_time_ms
		y := value.smoothed_delay_ms
		numerator += (x-x_avg)*(y-y_avg)
		denominator += (x-x_avg)*(x-x_avg)
	}
	if(denominator == 0){
		return 0.0
	}

	return numerator/denominator

}

func (trendline_estimator *TrendlineEstimator)UpdateThreshold(modified_trend float64, now_ms int64){
	if(trendline_estimator.last_update_ms_ == -1){
		trendline_estimator.last_update_ms_ = now_ms
	}
	if( math.Abs(modified_trend) > trendline_estimator.threshold_ + kMaxAdaptOffsetMs){
		trendline_estimator.last_update_ms_ = now_ms
		return
	}
	//fmt.Println("modify_trend", modified_trend)

	var k float64
	if(math.Abs(modified_trend) < trendline_estimator.threshold_){
		k = k_down
	}else{
		k = k_up
	}
	kMaxTimeDeltaMs := int64(100)
	time_delta_ms := math.Min(float64(now_ms-trendline_estimator.last_update_ms_), float64(kMaxTimeDeltaMs))
	trendline_estimator.threshold_ += k * (math.Abs(modified_trend)-trendline_estimator.threshold_) * time_delta_ms
	if(trendline_estimator.threshold_ < 6.0){
		trendline_estimator.threshold_ = 6.0
	}else if(trendline_estimator.threshold_ > 600){
		trendline_estimator.threshold_ = 600
	}
	trendline_estimator.last_update_ms_ = now_ms
}

func (trendline_estimator* TrendlineEstimator)Detect(trend float64, ts_delta float64, now_ms int64){
	if(trendline_estimator.num_of_deltas_ < 2){
		trendline_estimator.hypothesis_ = kBwNormal
		return
	}

	modified_trend := math.Min(float64(trendline_estimator.num_of_deltas_), kMinNumDeltas) * trend * threshold_gain

	if(modified_trend > trendline_estimator.threshold_){

		if(trendline_estimator.time_over_using_ == -1){
			trendline_estimator.time_over_using_ = ts_delta/2
		}else{
			trendline_estimator.time_over_using_ += ts_delta
		}
		trendline_estimator.overuse_counter_++

		if(trendline_estimator.time_over_using_ > kOverUsingTimeThreshold && trendline_estimator.overuse_counter_ > 1){
			if(trend >= trendline_estimator.prev_trend_){
				trendline_estimator.time_over_using_ = 0
				trendline_estimator.overuse_counter_ = 0
				trendline_estimator.hypothesis_ = kBwOverusing
			}
		}
	}else if(modified_trend < -trendline_estimator.threshold_){
		trendline_estimator.time_over_using_ = -1
		trendline_estimator.overuse_counter_ = 0
		trendline_estimator.hypothesis_ = kBwUnderusing
	}else{
		trendline_estimator.time_over_using_ = -1
		trendline_estimator.overuse_counter_ = 0
		trendline_estimator.hypothesis_ = kBwNormal
	}
	trendline_estimator.prev_trend_ = trend
	trendline_estimator.UpdateThreshold(modified_trend, now_ms)
}

func (trendline_estimator *TrendlineEstimator)UpdateTrendline(recv_delta_ms float64, send_delta_ms float64, send_time_ms int64, arrival_time_ms int64, packet_size int64){
	delta_ms := recv_delta_ms - send_delta_ms
	trendline_estimator.num_of_deltas_++
	if(trendline_estimator.num_of_deltas_ > kDeltaCounterMax){
		trendline_estimator.num_of_deltas_ = kDeltaCounterMax
	}
	if(trendline_estimator.first_arrival_time_ms_ == -1){
		trendline_estimator.first_arrival_time_ms_ = arrival_time_ms
	}

	trendline_estimator.accumulated_delay_ += delta_ms
	trendline_estimator.smooth_delay_ = smoothing_coef * trendline_estimator.smooth_delay_ + (1-smoothing_coef)*trendline_estimator.accumulated_delay_

	trendline_estimator.delay_hist_ = append(trendline_estimator.delay_hist_, PacketTiming{arrival_time_ms:float64(arrival_time_ms-trendline_estimator.first_arrival_time_ms_), smoothed_delay_ms: trendline_estimator.smooth_delay_, raw_delay_ms: trendline_estimator.accumulated_delay_})
	if(len(trendline_estimator.delay_hist_) > window_size){
		trendline_estimator.delay_hist_ = trendline_estimator.delay_hist_[1:]
	}
	trend := trendline_estimator.prev_trend_
	if(len(trendline_estimator.delay_hist_) == window_size){
		trend = LinearFitSlope(trendline_estimator.delay_hist_)
	}
	trendline_estimator.record_trend = trend
	//fmt.Println("trend:", trend)

	trendline_estimator.Detect(trend, send_delta_ms, arrival_time_ms)
}

func (trendline_estimator *TrendlineEstimator)Update(recv_delta_ms float64, send_delta_ms float64, send_time_ms int64, arrival_time_ms int64, packet_size int64, calculated_deltas bool){
	if(calculated_deltas){
		trendline_estimator.UpdateTrendline(recv_delta_ms, send_delta_ms, send_time_ms, arrival_time_ms, packet_size)
	}
}