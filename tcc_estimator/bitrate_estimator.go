package tccEstimator

import (
	"math"
)
const(
	small_sample_threshold = 0
	noninitial_window_ms = 150
	initial_window_ms = 500
	uncertainty_scale = 10.0
	small_sample_uncertainty_scale = 10.0
	uncertainty_scale_in_alr = 10.0

)
type BitrateEstimator struct {
	current_window_ms_ int64
	prev_time_ms_ int64
	bitrate_estimate_kbps_ float32
	bitrate_estimate_var_ float32
	sum_ int
}

func NewBitrateEstimator() BitrateEstimator{
	return BitrateEstimator{current_window_ms_: 0, prev_time_ms_: -1, bitrate_estimate_var_: 50.0, bitrate_estimate_kbps_: -1.0, sum_: 0}
}

func(bitrate_estimator *BitrateEstimator) UpdateWindow(now_ms int64, bytes int, rate_window_ms int, is_samll_sample *bool) float32{
	if(now_ms < bitrate_estimator.prev_time_ms_){
		bitrate_estimator.prev_time_ms_ = -1
		bitrate_estimator.sum_ = 0
		bitrate_estimator.current_window_ms_ = 0
	}

	if(bitrate_estimator.prev_time_ms_ >= 0){
		bitrate_estimator.current_window_ms_ += now_ms - bitrate_estimator.prev_time_ms_
		if(now_ms- bitrate_estimator.prev_time_ms_ > int64(rate_window_ms)){
			bitrate_estimator.sum_ = 0
			bitrate_estimator.current_window_ms_ %= int64(rate_window_ms)

		}
	}

	bitrate_estimator.prev_time_ms_ = now_ms
	bitrate_sample := float32(-1.0)
	if(bitrate_estimator.current_window_ms_ >= int64(rate_window_ms)){
		*is_samll_sample = (bitrate_estimator.sum_ < small_sample_threshold)
		bitrate_sample = 8.0 * float32(bitrate_estimator.sum_)/float32(rate_window_ms)
		bitrate_estimator.current_window_ms_ -= int64(rate_window_ms)
		bitrate_estimator.sum_ = 0
	}
	bitrate_estimator.sum_ += bytes
	return bitrate_sample
}

func(bitrate_estimator* BitrateEstimator) Update(at_time int64, amount int64, in_alr bool){
	rate_windows_ms := noninitial_window_ms
	if(bitrate_estimator.bitrate_estimate_kbps_ < float32(0.0)){
		rate_windows_ms = initial_window_ms
	}
	is_small_sample := false
	bitrate_sample_kbps := bitrate_estimator.UpdateWindow(at_time, int(amount), rate_windows_ms, &is_small_sample)
	if(bitrate_sample_kbps < float32(0.0)){
		return
	}

	if(bitrate_estimator.bitrate_estimate_kbps_ < 0.0){
		bitrate_estimator.bitrate_estimate_kbps_ = bitrate_sample_kbps
		return
	}

	var scale float32 = uncertainty_scale
	if(is_small_sample && bitrate_sample_kbps < bitrate_estimator.bitrate_estimate_kbps_){
		scale = small_sample_uncertainty_scale
	} else if(in_alr && bitrate_sample_kbps < bitrate_estimator.bitrate_estimate_kbps_){
		scale = uncertainty_scale_in_alr
	}

	sample_uncertainty := scale * float32(math.Abs(float64(bitrate_estimator.bitrate_estimate_kbps_-bitrate_sample_kbps)))/bitrate_estimator.bitrate_estimate_kbps_
	sample_var := sample_uncertainty * sample_uncertainty
	pred_bitrate_estimate_var := bitrate_estimator.bitrate_estimate_var_ + 5.0
	bitrate_estimator.bitrate_estimate_kbps_ = (sample_var*bitrate_estimator.bitrate_estimate_kbps_+pred_bitrate_estimate_var*bitrate_sample_kbps)/(sample_var+pred_bitrate_estimate_var)
	if(bitrate_estimator.bitrate_estimate_kbps_ < 0){
		bitrate_estimator.bitrate_estimate_kbps_ = 0.0
	}
	bitrate_estimator.bitrate_estimate_var_ = sample_var*pred_bitrate_estimate_var/(sample_var+pred_bitrate_estimate_var)
	//fmt.Println("bitrate_sample_kbps is ",bitrate_sample_kbps , "kpbs_ is ", bitrate_estimator.bitrate_estimate_kbps_, "var_ is ", bitrate_estimator.bitrate_estimate_var_)
}
