package tccEstimator

type GoogCcNetworkController struct {
	delay_based_bwe DelayBasedBwe
	acknowledged_bitrate_estimator AcknowledgeBitrateEstimator
	delay_based_bitrate int64
	is_first_test bool
}

func NewGoogCcNetworkController() GoogCcNetworkController{
	return GoogCcNetworkController{
				delay_based_bwe: NewDelayBasedBwe(),
				acknowledged_bitrate_estimator:NewAcknowledgeBitrateEstimator(),
				delay_based_bitrate: 450000,
				is_first_test: true,
			}
}

func (goog_cc_network_controller *GoogCcNetworkController)TestVersionOnTransportPacketsFeedback(report TransportPacketsFeedback,probe_birate int64, rtt int64)Result{
	goog_cc_network_controller.acknowledged_bitrate_estimator.IncomingPacketFeedbackVector(report.packet_feedbacks)
	acknowledged_bitrate := goog_cc_network_controller.acknowledged_bitrate_estimator.GetEstimator()
	result := goog_cc_network_controller.delay_based_bwe.IncomingPacketFeedbackVector(report, acknowledged_bitrate, probe_birate, rtt)
	return result
}

func (goog_cc_network_controller *GoogCcNetworkController)OnTransportPacketsFeedback(report TransportPacketsFeedback)int64{
	goog_cc_network_controller.acknowledged_bitrate_estimator.IncomingPacketFeedbackVector(report.packet_feedbacks)
	acknowledged_bitrate := goog_cc_network_controller.acknowledged_bitrate_estimator.GetEstimator()
	var result Result
	if(goog_cc_network_controller.is_first_test){
		//后面改掉rrt
		result = goog_cc_network_controller.delay_based_bwe.IncomingPacketFeedbackVector(report, acknowledged_bitrate, 450000, 0)
		goog_cc_network_controller.is_first_test = false
	}else{
		result = goog_cc_network_controller.delay_based_bwe.IncomingPacketFeedbackVector(report, acknowledged_bitrate, -1, 0)
	}
	if(result.updated){
		goog_cc_network_controller.delay_based_bitrate = result.target_bitrate
		return result.target_bitrate
	}else{
		return goog_cc_network_controller.delay_based_bitrate
	}
}

func (goog_cc_network_controller *GoogCcNetworkController)SetRtt(rtt int64){
	goog_cc_network_controller.delay_based_bwe.rate_control_.rtt_ = rtt
}
