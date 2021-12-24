package tccEstimator

import (
	"fmt"
	"github.com/pion/rtcp"
	"sort"
	"sync"
)

const (
	BaseRefrenceTimeScale = rtcp.TypeTCCDeltaScaleFactor * 256
	//窗口大小60s = 60 *1000000
	kSendTimeHistoryWindow = 60 * 1000000
)

type SendPacket struct {
	send_time       int64
	size            int64
	sequence_number int64
}
type PacketFeedback struct {
	creation_time int64
	sent          SendPacket
	receive_time  int64
}

type TransportFeedbackAdapter struct {
	history_          map[int64]PacketFeedback
	last_timestamp_   int64 //初始化-1
	current_offset    int64
	last_ack_seq_num_ int64
	mutex             sync.Mutex
	sequence_num_unwrapper_ SequenceNumberUnwrapper
}

type PacketResult struct {
	send_packet  SendPacket
	receive_time int64
}

func NewTransportFeedbackAdapter() TransportFeedbackAdapter {
	return TransportFeedbackAdapter{
		history_:          make(map[int64]PacketFeedback),
		last_timestamp_:   -1,
		current_offset:    -1,
		last_ack_seq_num_: -1,
		sequence_num_unwrapper_: NewSequenceNumberUnwrapper(),
	}
}

type PackResultSlice []PacketResult

func (packet_result PackResultSlice) Less(i, j int) bool {
	if packet_result[i].receive_time != packet_result[j].receive_time {
		return packet_result[i].receive_time < packet_result[j].receive_time
	}

	if packet_result[i].send_packet.send_time != packet_result[j].send_packet.send_time {
		return packet_result[i].send_packet.send_time < packet_result[j].send_packet.send_time
	}

	return packet_result[i].send_packet.sequence_number < packet_result[j].send_packet.sequence_number
}

func (s PackResultSlice) Len() int      { return len(s) }
func (s PackResultSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

//每次发送包需要调用此处
func (transport_feedback_adpater *TransportFeedbackAdapter) AddPacket(send_packet SendPacket, creation_time int64) {
	transport_feedback_adpater.mutex.Lock()
	defer transport_feedback_adpater.mutex.Unlock()

	var packet PacketFeedback
	packet.sent = send_packet
	packet.sent.sequence_number = transport_feedback_adpater.sequence_num_unwrapper_.Unwrap(uint16(packet.sent.sequence_number))
	packet.creation_time = creation_time

	for key, value := range transport_feedback_adpater.history_ {
		if creation_time-value.creation_time <= kSendTimeHistoryWindow {
			break
		}
		//超过窗口大小，删除记录
		delete(transport_feedback_adpater.history_, key)
	}
	transport_feedback_adpater.history_[send_packet.sequence_number] = packet
}

func (transport_feedback_adpater *TransportFeedbackAdapter) ProcessTransportFeedbackInner(feedback rtcp.TransportLayerCC, feedback_receive_time int64) PackResultSlice {
	transport_feedback_adpater.mutex.Lock()
	defer transport_feedback_adpater.mutex.Unlock()

	if transport_feedback_adpater.last_timestamp_ == -1 {
		transport_feedback_adpater.current_offset = feedback_receive_time
	} else {
		delta := (int64(feedback.ReferenceTime*BaseRefrenceTimeScale) - transport_feedback_adpater.last_timestamp_) / 1000 * 1000
		transport_feedback_adpater.current_offset += delta
	}

	transport_feedback_adpater.last_timestamp_ = int64(feedback.ReferenceTime * BaseRefrenceTimeScale)
	packet_result_vector := make(PackResultSlice, 0)
	fail_lookups := 0
	//ignored := 0
	packet_offset := int64(0)
	for index, value := range feedback.RecvDeltas {
		seq_num := transport_feedback_adpater.sequence_num_unwrapper_.Unwrap(uint16(index + int(feedback.BaseSequenceNumber))) //需要warp
		if seq_num > transport_feedback_adpater.last_ack_seq_num_ {
			transport_feedback_adpater.last_ack_seq_num_ = seq_num
		}
		send_packet_info, ok := transport_feedback_adpater.history_[seq_num]
		if !ok {
			fail_lookups++
			continue
		}
		/*
			if(send_packet_info.sent.send_time == -1){

			}*/

		if value.Type == rtcp.TypeTCCPacketReceivedSmallDelta || value.Type == rtcp.TypeTCCPacketReceivedLargeDelta {
			packet_offset += value.Delta
			send_packet_info.receive_time = transport_feedback_adpater.current_offset + packet_offset/1000*1000
			delete(transport_feedback_adpater.history_, seq_num)
		} else {
			//没有收到的包还是不放到PacketResult里面
			continue
		}

		var result PacketResult
		result.receive_time = send_packet_info.receive_time
		result.send_packet = send_packet_info.sent
		packet_result_vector = append(packet_result_vector, result)

		if fail_lookups > 0 {
			fmt.Println("failed to lookup ", fail_lookups)
		}
	}

	//排好序
	sort.Sort(packet_result_vector)
	return packet_result_vector
}

func (transport_feedback_adpater *TransportFeedbackAdapter) ProcessTransportFeedback(feedback rtcp.TransportLayerCC, feedback_receive_time int64) TransportPacketsFeedback {
	var msg TransportPacketsFeedback
	msg.packet_feedbacks = transport_feedback_adpater.ProcessTransportFeedbackInner(feedback, feedback_receive_time)
	msg.feedback_time = feedback_receive_time
	return msg
}
