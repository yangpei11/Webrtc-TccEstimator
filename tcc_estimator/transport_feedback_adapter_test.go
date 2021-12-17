package tccEstimator

import (
	"fmt"
	"github.com/pion/rtcp"
	"sort"
	"testing"
	"time"
)


func TestMapShunXu(t *testing.T) {
	var a int64
	a = 22500
	fmt.Println(a / 1000 * 1000)
}

func TestFeedbackData(t *testing.T) {
	transport_feedback_adapter := TransportFeedbackAdapter{
		history_:          make(map[int64]PacketFeedback),
		last_timestamp_:   -1,
		current_offset:    -1,
		last_ack_seq_num_: -1,
	}
	i := 0
	for i < 1000 {
		send_packet := SendPacket{sequence_number: int64(i), send_time: 0, size: 0}
		transport_feedback_adapter.AddPacket(send_packet, time.Now().UnixNano()/1000)
		i++
	}
	feedback_receive_time := int64(715407752000)
	BaseSequenceNumebr := 1
	Basetime := 4859630
	var feedback rtcp.TransportLayerCC
	feedback.BaseSequenceNumber = uint16(BaseSequenceNumebr)
	feedback.ReferenceTime = uint32(Basetime)
	feedback.RecvDeltas = make([]*rtcp.RecvDelta, 0)
	delta := []int64{29750, 1250, 12750, 0, 8750, 9750, 10500, 22750, 0, 0, 750, 1750, 5000}
	for _, value := range delta {
		feedback.RecvDeltas = append(feedback.RecvDeltas, &rtcp.RecvDelta{Delta: value, Type: rtcp.TypeTCCPacketReceivedSmallDelta})
	}
	packet_result := transport_feedback_adapter.ProcessTransportFeedbackInner(feedback, feedback_receive_time)

	for _, value := range packet_result {
		fmt.Println("sequence num:", value.send_packet.sequence_number, " recv time:", value.receive_time)
	}

	feedback_receive_time = 715407854
	BaseSequenceNumebr = 14
	Basetime = 4859632
	feedback.BaseSequenceNumber = uint16(BaseSequenceNumebr)
	feedback.ReferenceTime = uint32(Basetime)
	feedback.RecvDeltas = make([]*rtcp.RecvDelta, 0)
	delta1 := []int64{10250, 0, 0, 0, 8750, 0, 11000, 6000, 0, 5250, 8250, 0, 25500, 250, 1750, 1000, 11000, 0, 250, 0}
	for _, value := range delta1 {
		feedback.RecvDeltas = append(feedback.RecvDeltas, &rtcp.RecvDelta{Delta: value, Type: rtcp.TypeTCCPacketReceivedSmallDelta})
	}
	packet_result = transport_feedback_adapter.ProcessTransportFeedbackInner(feedback, feedback_receive_time)
	for _, value := range packet_result {
		fmt.Println("sequence num:", value.send_packet.sequence_number, " recv time:", value.receive_time)
	}
}

func TestHistoryData(t *testing.T) {
	transport_feedback_adapter := TransportFeedbackAdapter{
		history_:          make(map[int64]PacketFeedback),
		last_timestamp_:   -1,
		current_offset:    -1,
		last_ack_seq_num_: -1,
	}
	i := 0
	creation_time_array := []int64{1, 2, kSendTimeHistoryWindow + 10, kSendTimeHistoryWindow * 3, kSendTimeHistoryWindow * 5}
	for i < 5 {
		send_packet := SendPacket{sequence_number: int64(i), send_time: 0, size: 0}
		transport_feedback_adapter.AddPacket(send_packet, creation_time_array[i])
		fmt.Println(transport_feedback_adapter.history_)
		i++
	}
}

func TestPacketResultSort(t *testing.T) {
	packet_result := PackResultSlice{}
	packet_result = append(packet_result, PacketResult{receive_time: 2, send_packet: SendPacket{sequence_number: 5, send_time: 11}})
	packet_result = append(packet_result, PacketResult{receive_time: 2, send_packet: SendPacket{sequence_number: 1000, send_time: 10}})
	packet_result = append(packet_result, PacketResult{receive_time: 1, send_packet: SendPacket{sequence_number: 5, send_time: 10}})
	packet_result = append(packet_result, PacketResult{receive_time: 1, send_packet: SendPacket{sequence_number: 6, send_time: 10}})
	sort.Sort(packet_result)
	for _, value := range packet_result {
		fmt.Println("recv time: ", value.receive_time, " sendtime: ", value.send_packet.send_time, " sequence_number: ", value.send_packet.sequence_number)
	}
}

func TestAcknowledgeBitrateEstimator_IncomingPacketFeedbackVector(t *testing.T) {
	acknowledge_bitrate_estimator := NewAcknowledgeBitrateEstimator()

	packet_result0 := PackResultSlice{
		PacketResult{receive_time: 1061906417000, send_packet: SendPacket{size: 23}},
		PacketResult{receive_time: 1061906417000, send_packet: SendPacket{size: 1099}},
		PacketResult{receive_time: 1061906439000, send_packet: SendPacket{size: 23}},
		PacketResult{receive_time: 1061906439000, send_packet: SendPacket{size: 1135}},
		PacketResult{receive_time: 1061906439000, send_packet: SendPacket{size: 1070}},
		PacketResult{receive_time: 1061906447000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061906455000, send_packet: SendPacket{size: 1070}},
		PacketResult{receive_time: 1061906465000, send_packet: SendPacket{size: 1010}},
		PacketResult{receive_time: 1061906474000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061906475000, send_packet: SendPacket{size: 1053}},
		PacketResult{receive_time: 1061906478000, send_packet: SendPacket{size: 1070}},
		PacketResult{receive_time: 1061906484000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061906507000, send_packet: SendPacket{size: 249}},
		PacketResult{receive_time: 1061906507000, send_packet: SendPacket{size: 1068}},
		PacketResult{receive_time: 1061906509000, send_packet: SendPacket{size: 1071}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result0)
	packet_result1 := PackResultSlice{
		PacketResult{receive_time: 1061906517000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061906527000, send_packet: SendPacket{size: 270}},
		PacketResult{receive_time: 1061906527000, send_packet: SendPacket{size: 1071}},
		PacketResult{receive_time: 1061906539000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061906548000, send_packet: SendPacket{size: 23}},
		PacketResult{receive_time: 1061906561000, send_packet: SendPacket{size: 982}},
		PacketResult{receive_time: 1061906561000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061906561000, send_packet: SendPacket{size: 1048}},
		PacketResult{receive_time: 1061906563000, send_packet: SendPacket{size: 530}},
		PacketResult{receive_time: 1061906563000, send_packet: SendPacket{size: 982}},
		PacketResult{receive_time: 1061906567000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061906569000, send_packet: SendPacket{size: 982}},
		PacketResult{receive_time: 1061906571000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061906573000, send_packet: SendPacket{size: 221}},
		PacketResult{receive_time: 1061906573000, send_packet: SendPacket{size: 242}},
		PacketResult{receive_time: 1061906573000, send_packet: SendPacket{size: 1020}},
		PacketResult{receive_time: 1061906585000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061906586000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061906586000, send_packet: SendPacket{size: 24}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result1)
	packet_result2 := PackResultSlice{
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1024}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1024}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1025}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1025}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 990}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1025}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1025}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1025}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 24}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906619000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906620000, send_packet: SendPacket{size: 1040}},
		PacketResult{receive_time: 1061906620000, send_packet: SendPacket{size: 1040}},
		PacketResult{receive_time: 1061906620000, send_packet: SendPacket{size: 502}},
		PacketResult{receive_time: 1061906626000, send_packet: SendPacket{size: 958}},
		PacketResult{receive_time: 1061906626000, send_packet: SendPacket{size: 1021}},
		PacketResult{receive_time: 1061906658000, send_packet: SendPacket{size: 1049}},
		PacketResult{receive_time: 1061906658000, send_packet: SendPacket{size: 929}},
		PacketResult{receive_time: 1061906658000, send_packet: SendPacket{size: 930}},
		PacketResult{receive_time: 1061906658000, send_packet: SendPacket{size: 671}},
		PacketResult{receive_time: 1061906658000, send_packet: SendPacket{size: 700}},
		PacketResult{receive_time: 1061906664000, send_packet: SendPacket{size: 672}},
		PacketResult{receive_time: 1061906664000, send_packet: SendPacket{size: 674}},
		PacketResult{receive_time: 1061906664000, send_packet: SendPacket{size: 673}},
		PacketResult{receive_time: 1061906664000, send_packet: SendPacket{size: 1023}},
		PacketResult{receive_time: 1061906666000, send_packet: SendPacket{size: 964}},
		PacketResult{receive_time: 1061906668000, send_packet: SendPacket{size: 674}},
		PacketResult{receive_time: 1061906668000, send_packet: SendPacket{size: 673}},
		PacketResult{receive_time: 1061906673000, send_packet: SendPacket{size: 1023}},
		PacketResult{receive_time: 1061906688000, send_packet: SendPacket{size: 648}},
		PacketResult{receive_time: 1061906688000, send_packet: SendPacket{size: 619}},
		PacketResult{receive_time: 1061906688000, send_packet: SendPacket{size: 620}},
		PacketResult{receive_time: 1061906695000, send_packet: SendPacket{size: 856}},
		PacketResult{receive_time: 1061906695000, send_packet: SendPacket{size: 827}},
		PacketResult{receive_time: 1061906695000, send_packet: SendPacket{size: 827}},
		PacketResult{receive_time: 1061906696000, send_packet: SendPacket{size: 828}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result2)
	packet_result3 := PackResultSlice{
		PacketResult{receive_time: 1061906738000, send_packet: SendPacket{size: 875}},
		PacketResult{receive_time: 1061906738000, send_packet: SendPacket{size: 846}},
		PacketResult{receive_time: 1061906738000, send_packet: SendPacket{size: 847}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 1051}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 544}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 516}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 1023}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 1023}},
		PacketResult{receive_time: 1061906776000, send_packet: SendPacket{size: 1023}},
		PacketResult{receive_time: 1061906804000, send_packet: SendPacket{size: 1015}},
		PacketResult{receive_time: 1061906804000, send_packet: SendPacket{size: 886}},
		PacketResult{receive_time: 1061906804000, send_packet: SendPacket{size: 857}},
		PacketResult{receive_time: 1061906804000, send_packet: SendPacket{size: 858}},
		PacketResult{receive_time: 1061906804000, send_packet: SendPacket{size: 986}},
		PacketResult{receive_time: 1061906806000, send_packet: SendPacket{size: 987}},
		PacketResult{receive_time: 1061906807000, send_packet: SendPacket{size: 987}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result3)
	packet_result4 := PackResultSlice{
		PacketResult{receive_time: 1061906841000, send_packet: SendPacket{size: 903}},
		PacketResult{receive_time: 1061906841000, send_packet: SendPacket{size: 916}},
		PacketResult{receive_time: 1061906841000, send_packet: SendPacket{size: 875}},
		PacketResult{receive_time: 1061906842000, send_packet: SendPacket{size: 875}},
		PacketResult{receive_time: 1061906842000, send_packet: SendPacket{size: 887}},
		PacketResult{receive_time: 1061906842000, send_packet: SendPacket{size: 887}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 888}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 1166}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 1022}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 1138}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 1138}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 993}},
		PacketResult{receive_time: 1061906874000, send_packet: SendPacket{size: 993}},
		PacketResult{receive_time: 1061906899000, send_packet: SendPacket{size: 994}},
		PacketResult{receive_time: 1061906900000, send_packet: SendPacket{size: 989}},
		PacketResult{receive_time: 1061906900000, send_packet: SendPacket{size: 990}},
		PacketResult{receive_time: 1061906900000, send_packet: SendPacket{size: 960}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result4)
	packet_result5 := PackResultSlice{
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 961}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 961}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906939000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906940000, send_packet: SendPacket{size: 962}},
		PacketResult{receive_time: 1061906940000, send_packet: SendPacket{size: 1067}},
		PacketResult{receive_time: 1061906940000, send_packet: SendPacket{size: 1104}},
		PacketResult{receive_time: 1061906940000, send_packet: SendPacket{size: 1038}},
		PacketResult{receive_time: 1061906940000, send_packet: SendPacket{size: 1038}},
		PacketResult{receive_time: 1061906978000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061906978000, send_packet: SendPacket{size: 1075}},
		PacketResult{receive_time: 1061906978000, send_packet: SendPacket{size: 1076}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result5)
	packet_result6 := PackResultSlice{
		PacketResult{receive_time: 1061907046000, send_packet: SendPacket{size: 1076}},
		PacketResult{receive_time: 1061907046000, send_packet: SendPacket{size: 1076}},
		PacketResult{receive_time: 1061907047000, send_packet: SendPacket{size: 1076}},
		PacketResult{receive_time: 1061907052000, send_packet: SendPacket{size: 1039}},
		PacketResult{receive_time: 1061907052000, send_packet: SendPacket{size: 1074}},
		PacketResult{receive_time: 1061907087000, send_packet: SendPacket{size: 1010}},
		PacketResult{receive_time: 1061907088000, send_packet: SendPacket{size: 1010}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result6)
	packet_result7 := PackResultSlice{
		PacketResult{receive_time: 1061907124000, send_packet: SendPacket{size: 1011}},
		PacketResult{receive_time: 1061907124000, send_packet: SendPacket{size: 1045}},
		PacketResult{receive_time: 1061907146000, send_packet: SendPacket{size: 1011}},
		PacketResult{receive_time: 1061907146000, send_packet: SendPacket{size: 1045}},
		PacketResult{receive_time: 1061907181000, send_packet: SendPacket{size: 1046}},
		PacketResult{receive_time: 1061907181000, send_packet: SendPacket{size: 1046}},
		PacketResult{receive_time: 1061907208000, send_packet: SendPacket{size: 1046}},
		PacketResult{receive_time: 1061907208000, send_packet: SendPacket{size: 1139}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result7)
	packet_result8 := PackResultSlice{
		PacketResult{receive_time: 1061907236000, send_packet: SendPacket{size: 1111}},
		PacketResult{receive_time: 1061907237000, send_packet: SendPacket{size: 1111}},
		PacketResult{receive_time: 1061907267000, send_packet: SendPacket{size: 988}},
		PacketResult{receive_time: 1061907269000, send_packet: SendPacket{size: 959}},
		PacketResult{receive_time: 1061907269000, send_packet: SendPacket{size: 959}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result8)
	packet_result9 := PackResultSlice{
		PacketResult{receive_time: 1061907340000, send_packet: SendPacket{size: 959}},
		PacketResult{receive_time: 1061907340000, send_packet: SendPacket{size: 959}},
		PacketResult{receive_time: 1061907340000, send_packet: SendPacket{size: 960}},
		PacketResult{receive_time: 1061907346000, send_packet: SendPacket{size: 840}},
		PacketResult{receive_time: 1061907375000, send_packet: SendPacket{size: 811}},
		PacketResult{receive_time: 1061907375000, send_packet: SendPacket{size: 812}},
		PacketResult{receive_time: 1061907375000, send_packet: SendPacket{size: 812}},
		PacketResult{receive_time: 1061907405000, send_packet: SendPacket{size: 941}},
		PacketResult{receive_time: 1061907405000, send_packet: SendPacket{size: 912}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result9)
	packet_result10 := PackResultSlice{
		PacketResult{receive_time: 1061907440000, send_packet: SendPacket{size: 912}},
		PacketResult{receive_time: 1061907440000, send_packet: SendPacket{size: 913}},
		PacketResult{receive_time: 1061907440000, send_packet: SendPacket{size: 913}},
		PacketResult{receive_time: 1061907483000, send_packet: SendPacket{size: 913}},
		PacketResult{receive_time: 1061907483000, send_packet: SendPacket{size: 815}},
		PacketResult{receive_time: 1061907483000, send_packet: SendPacket{size: 1047}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result10)
	packet_result11 := PackResultSlice{
		PacketResult{receive_time: 1061907522000, send_packet: SendPacket{size: 787}},
		PacketResult{receive_time: 1061907522000, send_packet: SendPacket{size: 787}},
		PacketResult{receive_time: 1061907522000, send_packet: SendPacket{size: 787}},
		PacketResult{receive_time: 1061907548000, send_packet: SendPacket{size: 1019}},
		PacketResult{receive_time: 1061907548000, send_packet: SendPacket{size: 1019}},
		PacketResult{receive_time: 1061907583000, send_packet: SendPacket{size: 1019}},
		PacketResult{receive_time: 1061907583000, send_packet: SendPacket{size: 1019}},
		PacketResult{receive_time: 1061907608000, send_packet: SendPacket{size: 806}},
		PacketResult{receive_time: 1061907608000, send_packet: SendPacket{size: 1111}},
		PacketResult{receive_time: 1061907608000, send_packet: SendPacket{size: 777}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result11)
	packet_result12 := PackResultSlice{
		PacketResult{receive_time: 1061907655000, send_packet: SendPacket{size: 778}},
		PacketResult{receive_time: 1061907655000, send_packet: SendPacket{size: 778}},
		PacketResult{receive_time: 1061907655000, send_packet: SendPacket{size: 1082}},
		PacketResult{receive_time: 1061907686000, send_packet: SendPacket{size: 1082}},
		PacketResult{receive_time: 1061907686000, send_packet: SendPacket{size: 1082}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result12)
	packet_result13 := PackResultSlice{
		PacketResult{receive_time: 1061907717000, send_packet: SendPacket{size: 1083}},
		PacketResult{receive_time: 1061907717000, send_packet: SendPacket{size: 1107}},
		PacketResult{receive_time: 1061907749000, send_packet: SendPacket{size: 827}},
		PacketResult{receive_time: 1061907749000, send_packet: SendPacket{size: 827}},
		PacketResult{receive_time: 1061907749000, send_packet: SendPacket{size: 827}},
		PacketResult{receive_time: 1061907781000, send_packet: SendPacket{size: 1078}},
		PacketResult{receive_time: 1061907781000, send_packet: SendPacket{size: 1079}},
		PacketResult{receive_time: 1061907810000, send_packet: SendPacket{size: 1079}},
		PacketResult{receive_time: 1061907812000, send_packet: SendPacket{size: 1079}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result13)
	packet_result14 := PackResultSlice{
		PacketResult{receive_time: 1061907852000, send_packet: SendPacket{size: 966}},
		PacketResult{receive_time: 1061907852000, send_packet: SendPacket{size: 953}},
		PacketResult{receive_time: 1061907878000, send_packet: SendPacket{size: 937}},
		PacketResult{receive_time: 1061907878000, send_packet: SendPacket{size: 937}},
		PacketResult{receive_time: 1061907906000, send_packet: SendPacket{size: 938}},
		PacketResult{receive_time: 1061907908000, send_packet: SendPacket{size: 924}},
		PacketResult{receive_time: 1061907908000, send_packet: SendPacket{size: 924}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result14)
	packet_result15 := PackResultSlice{
		PacketResult{receive_time: 1061907988000, send_packet: SendPacket{size: 924}},
		PacketResult{receive_time: 1061907988000, send_packet: SendPacket{size: 925}},
		PacketResult{receive_time: 1061907988000, send_packet: SendPacket{size: 925}},
		PacketResult{receive_time: 1061907988000, send_packet: SendPacket{size: 987}},
		PacketResult{receive_time: 1061908013000, send_packet: SendPacket{size: 804}},
		PacketResult{receive_time: 1061908013000, send_packet: SendPacket{size: 805}},
		PacketResult{receive_time: 1061908013000, send_packet: SendPacket{size: 805}},
		PacketResult{receive_time: 1061908013000, send_packet: SendPacket{size: 958}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result15)
	packet_result16 := PackResultSlice{
		PacketResult{receive_time: 1061908077000, send_packet: SendPacket{size: 958}},
		PacketResult{receive_time: 1061908077000, send_packet: SendPacket{size: 958}},
		PacketResult{receive_time: 1061908077000, send_packet: SendPacket{size: 958}},
		PacketResult{receive_time: 1061908089000, send_packet: SendPacket{size: 959}},
		PacketResult{receive_time: 1061908097000, send_packet: SendPacket{size: 833}},
		PacketResult{receive_time: 1061908097000, send_packet: SendPacket{size: 947}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result16)
	packet_result17 := PackResultSlice{
		PacketResult{receive_time: 1061908119000, send_packet: SendPacket{size: 768}},
		PacketResult{receive_time: 1061908119000, send_packet: SendPacket{size: 768}},
		PacketResult{receive_time: 1061908120000, send_packet: SendPacket{size: 769}},
		PacketResult{receive_time: 1061908161000, send_packet: SendPacket{size: 919}},
		PacketResult{receive_time: 1061908161000, send_packet: SendPacket{size: 943}},
		PacketResult{receive_time: 1061908190000, send_packet: SendPacket{size: 842}},
		PacketResult{receive_time: 1061908190000, send_packet: SendPacket{size: 871}},
		PacketResult{receive_time: 1061908190000, send_packet: SendPacket{size: 1104}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result17)
	packet_result18 := PackResultSlice{
		PacketResult{receive_time: 1061908229000, send_packet: SendPacket{size: 817}},
		PacketResult{receive_time: 1061908229000, send_packet: SendPacket{size: 1106}},
		PacketResult{receive_time: 1061908229000, send_packet: SendPacket{size: 1074}},
		PacketResult{receive_time: 1061908274000, send_packet: SendPacket{size: 1078}},
		PacketResult{receive_time: 1061908274000, send_packet: SendPacket{size: 1078}},
		PacketResult{receive_time: 1061908274000, send_packet: SendPacket{size: 1068}},
		PacketResult{receive_time: 1061908306000, send_packet: SendPacket{size: 826}},
		PacketResult{receive_time: 1061908306000, send_packet: SendPacket{size: 898}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result18)
	packet_result19 := PackResultSlice{
		PacketResult{receive_time: 1061908334000, send_packet: SendPacket{size: 1050}},
		PacketResult{receive_time: 1061908367000, send_packet: SendPacket{size: 1056}},
		PacketResult{receive_time: 1061908397000, send_packet: SendPacket{size: 860}},
		PacketResult{receive_time: 1061908397000, send_packet: SendPacket{size: 860}},
		PacketResult{receive_time: 1061908397000, send_packet: SendPacket{size: 888}},
		PacketResult{receive_time: 1061908397000, send_packet: SendPacket{size: 1065}},
		PacketResult{receive_time: 1061908397000, send_packet: SendPacket{size: 821}},
	}
	acknowledge_bitrate_estimator.IncomingPacketFeedbackVector(packet_result19)
}
