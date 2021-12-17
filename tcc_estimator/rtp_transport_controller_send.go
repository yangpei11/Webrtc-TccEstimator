package tccEstimator

import (
	"github.com/pion/rtcp"
	"time"
)

type TransportPacketsFeedback struct {
	feedback_time    int64
	packet_feedbacks PackResultSlice
}

type RtpTransportControllerSend struct {
	transport_feedback_adapter_ TransportFeedbackAdapter
	control_                    GoogCcNetworkController
}

func NewRtpTransportControllerSend() *RtpTransportControllerSend {
	return &RtpTransportControllerSend{transport_feedback_adapter_: NewTransportFeedbackAdapter(), control_: NewGoogCcNetworkController()}
}

func (r *RtpTransportControllerSend) OnTransportFeedback(feedback rtcp.TransportLayerCC) int64 {
	feedback_time := time.Now().UnixNano() / 1000000 * 1000 //转换成微秒
	msg := r.transport_feedback_adapter_.ProcessTransportFeedback(feedback, feedback_time)
	bitrate := r.control_.OnTransportPacketsFeedback(msg)
	return bitrate
}

func (r *RtpTransportControllerSend) SetRtt(rtt int64) {
	r.control_.delay_based_bwe.rate_control_.rtt_ = rtt
}

func (r *RtpTransportControllerSend) AddPacket(send_time int64, size int64, sequence_number int64) {
	r.transport_feedback_adapter_.AddPacket(SendPacket{
		send_time:       send_time,
		size:            size,
		sequence_number: sequence_number,
	}, send_time)
}
