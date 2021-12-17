package tccEstimator

import (
	"math"
)

const (
	kRcHold     = 100
	kRcIncrease = 101
	kRcDecrease = 102

	kDefaultRtt                  = 200 * 1000
	kMinIncreaseRateBpsPerSecond = 4000
	beta                         = float64(0.85)
)

type AimdRateControl struct {
	current_bitrate_             int64
	lasted_estimated_throughput_ int64
	bitrate_is_initialized_      bool
	rate_control_state_          int
	time_last_bitrate_change_    int64
	link_capacity_               LinkCapacityEstimator
	rtt_                         int64
	last_decrease_               int64
	time_last_bitrate_decrease_  int64
	min_configured_bitrate_      int64
}

func NewAimdRateControl() AimdRateControl {
	return AimdRateControl{
		current_bitrate_:             450000,
		lasted_estimated_throughput_: 300000,
		bitrate_is_initialized_:      true,
		rate_control_state_:          kRcHold,
		time_last_bitrate_change_:    -1,
		link_capacity_:               LinkCapacityEstimator{estimate_kbps_: -1, deviation_kbps_: 0.4},
		rtt_:                         kDefaultRtt,
		last_decrease_:               -1,
		time_last_bitrate_decrease_:  -1,
		min_configured_bitrate_:      250000,
	}
}

func (amid_rate_control *AimdRateControl) SetStartBitrate(start_bitrate int64) {
	amid_rate_control.current_bitrate_ = start_bitrate
	amid_rate_control.lasted_estimated_throughput_ = amid_rate_control.current_bitrate_
	amid_rate_control.bitrate_is_initialized_ = true
	amid_rate_control.min_configured_bitrate_ = 450000
}

func (amid_rate_control *AimdRateControl) ChangeState(input RateControlInput, at_time int64) {
	switch input.bw_state {
	//如果现在是hold，网络状况是normal，尝试增大带宽
	case kBwNormal:
		if amid_rate_control.rate_control_state_ == kRcHold {
			amid_rate_control.time_last_bitrate_change_ = at_time
			amid_rate_control.rate_control_state_ = kRcIncrease
		}
		break
	case kBwOverusing:
		if amid_rate_control.rate_control_state_ != kRcDecrease {
			amid_rate_control.rate_control_state_ = kRcDecrease
		}
		break
	case kBwUnderusing:
		amid_rate_control.rate_control_state_ = kRcHold
		break

	default:
		amid_rate_control.rate_control_state_ = kRcHold
	}
}

func (amid_rate_control *AimdRateControl) GetNearMaxIncreaseRateBpsSecond() float64 {
	//kFrameInterval := 1000000/30.0
	//frame_size := float64(amid_rate_control.current_bitrate_ ) * kFrameInterval
	kFrameInterval := int64(1000000 / 30)
	frame_size := (amid_rate_control.current_bitrate_*kFrameInterval + 4000000) / 8000000

	kPacketSize := 1200.0
	packets_per_frame := math.Ceil(float64(frame_size) / kPacketSize)
	avg_packet_size := math.Round(float64(frame_size) / packets_per_frame)

	response_time := amid_rate_control.rtt_ + 100*1000
	response_time = response_time * 2

	increase_rate_bps_per_second := int64(avg_packet_size) * 1000000 * 8 / response_time
	return math.Max(kMinIncreaseRateBpsPerSecond, float64(increase_rate_bps_per_second))
}

func (amid_rate_control *AimdRateControl) AdditiveRateIncrease(at_time int64, last_time int64) int64 {
	time_period_seconds := float64(at_time-last_time) / 1000000.0
	//time_period_seconds := (at_time - last_time)/int64(1000000)
	data_rate_increase_bps := amid_rate_control.GetNearMaxIncreaseRateBpsSecond() * time_period_seconds
	return int64(data_rate_increase_bps)

}

func (amid_rate_control *AimdRateControl) MultiplicativeRateIncrease(at_time int64, last_time int64, current_bitrate int64) int64 {
	alpha := 1.08
	multiplicative_increase := math.Max(float64(current_bitrate)*(alpha-1.0), 1000.0)
	return int64(multiplicative_increase)
}

func (amid_rate_control *AimdRateControl) ChangeBitrate(input RateControlInput, at_time int64) {
	var new_bitrate int64
	new_bitrate = -1
	var estimated_throughput int64
	if input.estimated_throughput != -1 {
		estimated_throughput = input.estimated_throughput
		amid_rate_control.lasted_estimated_throughput_ = input.estimated_throughput
	} else {
		estimated_throughput = amid_rate_control.lasted_estimated_throughput_
	}

	if !amid_rate_control.bitrate_is_initialized_ && input.bw_state != kBwOverusing {
		return
	}

	amid_rate_control.ChangeState(input, at_time)

	throughput_base_limit := int64(1.5*float64(estimated_throughput) + 10000)

	switch amid_rate_control.rate_control_state_ {
	case kRcHold:
		break
	case kRcIncrease:
		x := amid_rate_control.link_capacity_.UpperBound()
		if estimated_throughput > x {
			amid_rate_control.link_capacity_.Reset()
		}
		if amid_rate_control.current_bitrate_ < throughput_base_limit {
			var increased_bitrate int64
			if amid_rate_control.link_capacity_.has_estimate() {
				additiveRateIncrease := amid_rate_control.AdditiveRateIncrease(at_time, amid_rate_control.time_last_bitrate_change_)
				increased_bitrate = int64(amid_rate_control.current_bitrate_ + additiveRateIncrease)
			} else {
				multiplicative_increase := amid_rate_control.MultiplicativeRateIncrease(at_time, amid_rate_control.time_last_bitrate_change_, amid_rate_control.current_bitrate_)
				increased_bitrate = amid_rate_control.current_bitrate_ + multiplicative_increase
			}
			new_bitrate = int64(math.Min(float64(increased_bitrate), float64(throughput_base_limit)))
		}
		amid_rate_control.time_last_bitrate_change_ = at_time
	case kRcDecrease:
		decreased_bitrate := int64(math.Round(float64(estimated_throughput) * beta))
		if decreased_bitrate > amid_rate_control.current_bitrate_ {
			if amid_rate_control.link_capacity_.has_estimate() {
				decreased_bitrate = int64(beta * amid_rate_control.link_capacity_.estimate_kbps_ * 1000.0)
			}
		}

		if decreased_bitrate < amid_rate_control.current_bitrate_ {
			new_bitrate = decreased_bitrate
		}

		if estimated_throughput < amid_rate_control.current_bitrate_ {
			if new_bitrate == -1 {
				amid_rate_control.last_decrease_ = 0
			} else {
				amid_rate_control.last_decrease_ = amid_rate_control.current_bitrate_ - new_bitrate
			}
		}

		if estimated_throughput < amid_rate_control.link_capacity_.LowerBound() {
			amid_rate_control.link_capacity_.Reset()
		}

		amid_rate_control.bitrate_is_initialized_ = true
		amid_rate_control.link_capacity_.OnOveruseDetected(estimated_throughput)
		amid_rate_control.rate_control_state_ = kRcHold
		amid_rate_control.time_last_bitrate_change_ = at_time
		amid_rate_control.time_last_bitrate_decrease_ = at_time
	}
	if new_bitrate != -1 {
		amid_rate_control.current_bitrate_ = new_bitrate
		amid_rate_control.current_bitrate_ = int64(math.Max(float64(amid_rate_control.min_configured_bitrate_), float64(amid_rate_control.current_bitrate_)))
	}

}

func Clamped(value int64, MinValue int64, MaxValue int64) int64 {
	if value < MinValue {
		return MinValue
	} else if value > MaxValue {
		return MaxValue
	} else {
		return value
	}
}

func (amid_rate_control *AimdRateControl) ClampBitrate(new_bitrate int64) int64 {
	return int64(math.Max(float64(amid_rate_control.min_configured_bitrate_), float64(new_bitrate)))
}

func (amid_rate_control *AimdRateControl) Update(input *RateControlInput, at_time int64) int64 {
	//之后补上未初始化的处理
	amid_rate_control.ChangeBitrate(*input, at_time)
	return amid_rate_control.current_bitrate_
}

func (amid_rate_control *AimdRateControl) SetEstimate(bitrate int64, at_time int64) {
	amid_rate_control.bitrate_is_initialized_ = true
	prev_bitrate := amid_rate_control.current_bitrate_
	amid_rate_control.current_bitrate_ = amid_rate_control.ClampBitrate(bitrate)
	amid_rate_control.time_last_bitrate_change_ = at_time
	if amid_rate_control.current_bitrate_ < prev_bitrate {
		amid_rate_control.time_last_bitrate_decrease_ = at_time
	}
}

func (amid_rate_control *AimdRateControl) TimeToReduceFurther(at_time int64, estimated_throughput int64) bool {
	bitrate_reduction_interval := Clamped(amid_rate_control.rtt_, 10*1000, 200*1000)
	if at_time-amid_rate_control.time_last_bitrate_change_ >= bitrate_reduction_interval {
		return true
	}

	if amid_rate_control.bitrate_is_initialized_ {
		threshold := 0.5 * float64(amid_rate_control.current_bitrate_)
		return float64(estimated_throughput) < threshold
	}
	return false
}
