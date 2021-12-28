# Transport-CC Algorithm

## Description
This is a Google's Transport-CC Implemention In go, which can be used to evaluate bandwidth.

## How to Use
rtp_transport_controller_send.go provied some Interfaces.
You can use Function **OnTransportFeedback** to get a bandwith from tcc algorithm, but you have to do three things to make sure the result correctly before. 
1. First, you must input rtcp packet which is defined as **rtcp.TransportLayerCC**. 
2. Second, When you receive **RR** Rtcp Packet, you should parse this packet to get **rtt** and **SetRtt**.
3. When you send rtp packet, you must **AddPacket**.

## File Structure
* rtp_transport_controller_send.go

RtpTransportControllerSend include TransportFeedbackAdapter and GoogCcNetworkController. NewRtpTransportControllerSend can create a object, we can use this object to make bandwith.

* transport_feedback_adapter.go

TransportFeedbackAdapter's **history** save the sending packet record. While receiving rtcp packet, return a vector of packet_result and feedback time.

* goog_cc_network_control.go
 
Congestion control module include a delay based rate estimator and window flow based rate estimator. There is a difference that webrtc use a ProbeEstimator to make Bandwidth detection in the beginning, but we First input 450000 as probe rate.

* delay_based_bwe.go

DelayBasedBwe include three modules, rate_control_, active_delay_detector_ and video_inter_arrival. 

* inter_arrival.go

InterArrival provide a arrival-time model, calculate delta between package groups in instead of single packet.

* trendline_estimator.go

TrendlineEstimator provide a over-use detector to maintain a status of current network.

* aimd_rate_control.go

AimdRateControl make a reasonable bandwidth from status of current network.

* acknowledged_bitrate_estimator.go

AcknowledgeBitrateEstimator calcalute a window flow bandwidth rate.

## Data Flow Chart
* see Figure FlowChart.png

