package tccEstimator

type AcknowledgeBitrateEstimator struct {
	bitrate_estimator_ BitrateEstimator
}

func NewAcknowledgeBitrateEstimator() AcknowledgeBitrateEstimator{
	return AcknowledgeBitrateEstimator{bitrate_estimator_: NewBitrateEstimator()}
}

func (acknowledged_bitrate_estimator *AcknowledgeBitrateEstimator)IncomingPacketFeedbackVector(packet_feedback_slice PackResultSlice){
	for _, value := range(packet_feedback_slice){
		packet_size := value.send_packet.size
		acknowledged_bitrate_estimator.bitrate_estimator_.Update(value.receive_time/1000, packet_size, false)
	}
}

//返回预测的kbps
func (acknowledged_bitrate_estimator *AcknowledgeBitrateEstimator) GetEstimator()int64{
	if(acknowledged_bitrate_estimator.bitrate_estimator_.bitrate_estimate_kbps_ == -1){
		return -1
	}else{
		return int64(acknowledged_bitrate_estimator.bitrate_estimator_.bitrate_estimate_kbps_*1000)
	}
}


