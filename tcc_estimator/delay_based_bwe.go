package tccEstimator

import "fmt"

type BandwidthUsage int32

const (
	kBwNormal     BandwidthUsage = 0
	kBwUnderusing BandwidthUsage = 1
	kBwOverusing  BandwidthUsage = 2

	kStreamTimeOut                  = 2 * 1000000
	kAbsSendTimeFraction            = 18
	kAbsSendTimeInterArrivalUpshift = 8
	kInterArrivalShift              = kAbsSendTimeInterArrivalUpshift + kAbsSendTimeFraction
	kTimestampGroupLengthMs         = 5
	kTimestampGroupTicks            = (kTimestampGroupLengthMs << kInterArrivalShift) / 1000

	kTimestampMs = float64(1000.0 / (1 << kInterArrivalShift))
)

type Result struct {
	updated                bool
	probe                  bool
	target_bitrate         int64
	recovered_from_overuse bool
	choose                 int
}

type RateControlInput struct {
	estimated_throughput int64
	bw_state             BandwidthUsage
}

func NewDefaultResult() Result {
	return Result{updated: false, probe: false, target_bitrate: 0, recovered_from_overuse: false}
}

func NewResultWithParams(probe bool, target_bitrate int64) Result {
	return Result{updated: true, probe: probe, target_bitrate: target_bitrate, recovered_from_overuse: false}
}

type DelayBasedBwe struct {
	active_delay_detector_ TrendlineEstimator
	last_seen_packet_      int64
	video_inter_arrival    InterArrival
	rate_control_          AimdRateControl
	prev_state_            BandwidthUsage
	pre_bitrate_           int64
}

func NewDelayBasedBwe() DelayBasedBwe {
	return DelayBasedBwe{
		rate_control_:          NewAimdRateControl(),
		active_delay_detector_: NewTrendlineEstimator(),
		video_inter_arrival:    NewInterArrival((5<<26)/1000, 1000.0/(1<<26), true),
		last_seen_packet_:      -1,
		prev_state_:            kBwNormal,
		pre_bitrate_:           0,
	}
}

func (delay_based_bwe *DelayBasedBwe) IncomingPacketFeedbackVector(msg TransportPacketsFeedback, acked_bitrate int64, probe_bitrate int64, rtt int64) Result {
	if len(msg.packet_feedbacks) == 0 {
		return NewDefaultResult()
	}
	delayed_feedback := true
	recovered_from_overuse := false
	prev_detector_state := delay_based_bwe.active_delay_detector_.State()

	for _, value := range msg.packet_feedbacks {
		delayed_feedback = false
		delay_based_bwe.IncomingPacketFeedback(value, msg.feedback_time)
		if prev_detector_state == kBwUnderusing && delay_based_bwe.active_delay_detector_.State() == kBwNormal {
			recovered_from_overuse = true
		}
		prev_detector_state = delay_based_bwe.active_delay_detector_.State()
	}

	if delayed_feedback {
		return NewDefaultResult()
	}

	return delay_based_bwe.MayUpdateEstimate(acked_bitrate, recovered_from_overuse, msg.feedback_time, probe_bitrate, rtt)

}

func (delay_based_bwe *DelayBasedBwe) IncomingPacketFeedback(packet_feedback PacketResult, at_time int64) {
	//第一次要初始化
	if delay_based_bwe.last_seen_packet_ == -1 && at_time-delay_based_bwe.last_seen_packet_ > kStreamTimeOut {
		delay_based_bwe.active_delay_detector_ = NewTrendlineEstimator()
		delay_based_bwe.video_inter_arrival = NewInterArrival(kTimestampGroupTicks, kTimestampMs, true)
		fmt.Println("First inits")
	}

	delay_based_bwe.last_seen_packet_ = at_time
	send_time_24bits := uint64(((uint64(packet_feedback.send_packet.send_time/1000)<<kAbsSendTimeFraction)+500)/1000) & 0x00FFFFFF
	timestamp := uint32(send_time_24bits << kAbsSendTimeInterArrivalUpshift)

	timestamp_delta := uint32(0)
	recv_delta_ms := int64(0)
	size_deleta := int64(0)
	calculated_deltas := delay_based_bwe.video_inter_arrival.ComputeDeltas(timestamp, (packet_feedback.receive_time / 1000), at_time/1000, packet_feedback.send_packet.size, &timestamp_delta, &recv_delta_ms, &size_deleta)

	send_delta_ms := (1000.0 * float64(timestamp_delta)) / (1 << kInterArrivalShift)
	delay_based_bwe.active_delay_detector_.Update(float64(recv_delta_ms), send_delta_ms, packet_feedback.send_packet.send_time/1000, packet_feedback.receive_time/1000, packet_feedback.send_packet.size, calculated_deltas)
}

func (delay_base_bwe *DelayBasedBwe) MayUpdateEstimate(acked_bitrate int64, recovered_from_overuse bool, at_time int64, probe_bitrate int64, rtt int64) Result {
	var result Result
	result.updated = false
	result.target_bitrate = 0
	var choose int
	//delay_base_bwe.active_delay_detector_.hypothesis_ = state
	//delay_base_bwe.rate_control_.rtt_ = rtt
	if delay_base_bwe.active_delay_detector_.State() == kBwOverusing {
		if acked_bitrate != -1 && delay_base_bwe.rate_control_.TimeToReduceFurther(at_time, acked_bitrate) {
			result.updated = delay_base_bwe.UpdateEstimate(at_time, acked_bitrate, &result.target_bitrate)
			choose = 1
		} else {
			choose = 0
		}
		//result.updated = delay_base_bwe.UpdateEstimate(at_time,acked_bitrate, &result.target_bitrate)
	} else {
		if probe_bitrate != -1 {
			result.probe = true
			result.updated = true
			result.target_bitrate = probe_bitrate
			delay_base_bwe.rate_control_.SetEstimate(probe_bitrate, at_time)
			choose = 2
		} else {
			result.updated = delay_base_bwe.UpdateEstimate(at_time, acked_bitrate, &result.target_bitrate)
			result.recovered_from_overuse = recovered_from_overuse
			choose = 3
		}
	}

	//一些状态改变
	dectector_state := delay_base_bwe.active_delay_detector_.State()
	if (result.updated && delay_base_bwe.pre_bitrate_ != result.target_bitrate) || dectector_state != delay_base_bwe.prev_state_ {
		var bitrate int64
		if result.updated {
			bitrate = result.target_bitrate
		} else {
			bitrate = delay_base_bwe.pre_bitrate_
		}
		delay_base_bwe.pre_bitrate_ = bitrate
		delay_base_bwe.prev_state_ = dectector_state
	}

	result.choose = choose
	return result
	//return result
}

func (delay_based_bwe *DelayBasedBwe) UpdateEstimate(at_time int64, acked_bitrate int64, target_bitate *int64) bool {
	input := &RateControlInput{estimated_throughput: acked_bitrate, bw_state: delay_based_bwe.active_delay_detector_.State()}
	*target_bitate = delay_based_bwe.rate_control_.Update(input, at_time)
	return delay_based_bwe.rate_control_.bitrate_is_initialized_
}
