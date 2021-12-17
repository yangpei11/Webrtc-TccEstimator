package tccEstimator

import "fmt"

const(
	kBurstDeltaThresholdMs = int(5)
	kMaxBurstDurationMs = int(100)
	kArrivalTimeOffsetThresholdMs = int64(3000)
	kRecorderedResetThreshold = 3
)
type TimestampGroup struct {
	size int64
	first_timestamp uint32
	timestamp uint32
	first_arrival_ms int64
	complete_time_ms int64
	last_system_time_ms int64
}

func (timestamp_group TimestampGroup) IsFirstPacket()bool{
	return timestamp_group.complete_time_ms == -1
}

func NewDefalutTimestampGroup()TimestampGroup{
	return TimestampGroup{size:0, first_timestamp: 0, timestamp: 0, first_arrival_ms: -1, complete_time_ms: -1}
}


type InterArrival struct{
	kTimestampGroupLengthTicks uint32
	current_timestamp_group_ TimestampGroup
	prev_timestamp_group_ TimestampGroup
	timestamp_to_ms_coeff_ float64
	burst_grouping_ bool
	num_consecutive_reordered_packets_ int
}

func NewInterArrival(timestamp_group_length_ticks uint32, timestamp_to_ms_coeff float64, enable_burst_grouping bool) InterArrival{
	return InterArrival{kTimestampGroupLengthTicks:timestamp_group_length_ticks, current_timestamp_group_: NewDefalutTimestampGroup(), prev_timestamp_group_: NewDefalutTimestampGroup(),timestamp_to_ms_coeff_: timestamp_to_ms_coeff, burst_grouping_: enable_burst_grouping,num_consecutive_reordered_packets_: 0}
}

func (inter_arrival* InterArrival)PacketInOrder(timestamp uint32)bool{
	if(inter_arrival.current_timestamp_group_.IsFirstPacket()){
		return true
	}else{
		timestamp_diff  := timestamp - inter_arrival.current_timestamp_group_.first_timestamp
		return (timestamp_diff < 0x80000000)
	}
}

func (inter_arrival *InterArrival)BelongsToBurst(arrival_time_ms int64, timestamp uint32)bool{
	if(!inter_arrival.burst_grouping_){
		return false
	}

	arrival_time_delta_ms := arrival_time_ms - inter_arrival.current_timestamp_group_.complete_time_ms
	timestamp_diff := timestamp - inter_arrival.current_timestamp_group_.timestamp
	ts_delta_ms := int64(inter_arrival.timestamp_to_ms_coeff_ * float64(timestamp_diff) + 0.5)
	if ts_delta_ms == 0{
		return true
	}

	propagation_delta_ms := arrival_time_delta_ms - ts_delta_ms
	if(propagation_delta_ms < 0 && arrival_time_delta_ms <= int64(kBurstDeltaThresholdMs) && arrival_time_ms - inter_arrival.current_timestamp_group_.first_arrival_ms< int64(kMaxBurstDurationMs)){
		return true
	}
	return false
}

func (inter_arrival *InterArrival)NewTimestampGroup(arrival_time_ms int64, timestamp uint32)bool{
	if(inter_arrival.current_timestamp_group_.IsFirstPacket()){
		return false
	}else if(inter_arrival.BelongsToBurst(arrival_time_ms, timestamp)) {
		return false
	}else{
		timestamp_diff := timestamp - inter_arrival.current_timestamp_group_.first_timestamp
		return timestamp_diff > inter_arrival.kTimestampGroupLengthTicks
	}
}

func (inter_arrival *InterArrival)ComputeDeltas(timestamp uint32, arrival_time_ms int64, system_time_ms int64, packet_size int64, timestamp_delta *uint32, arrival_time_delta_ms *int64, packet_size_delta *int64)bool{
	calculated := false
	//第一包
	if(inter_arrival.current_timestamp_group_.IsFirstPacket()) {
		inter_arrival.current_timestamp_group_.timestamp = timestamp
		inter_arrival.current_timestamp_group_.first_timestamp = timestamp
		inter_arrival.current_timestamp_group_.first_arrival_ms = arrival_time_ms
	}else if(!inter_arrival.PacketInOrder(timestamp)){
		return false;
	}else if(inter_arrival.NewTimestampGroup(arrival_time_ms, timestamp)){
		if(inter_arrival.prev_timestamp_group_.complete_time_ms >= 0){
			*timestamp_delta = inter_arrival.current_timestamp_group_.timestamp - inter_arrival.prev_timestamp_group_.timestamp
			*arrival_time_delta_ms = inter_arrival.current_timestamp_group_.complete_time_ms - inter_arrival.prev_timestamp_group_.complete_time_ms

			system_time_delta_ms := inter_arrival.current_timestamp_group_.last_system_time_ms - inter_arrival.prev_timestamp_group_.last_system_time_ms
			if(*arrival_time_delta_ms - system_time_delta_ms >= kArrivalTimeOffsetThresholdMs){
				inter_arrival.Reset()
				return false
			}
			if(*arrival_time_delta_ms <0 ){
				inter_arrival.num_consecutive_reordered_packets_++
				//收到的包乱序
				if(inter_arrival.num_consecutive_reordered_packets_ >= kRecorderedResetThreshold){
					inter_arrival.Reset()
				}
				return false
			}else {
				inter_arrival.num_consecutive_reordered_packets_ = 0
			}

			*packet_size_delta = inter_arrival.current_timestamp_group_.size - inter_arrival.prev_timestamp_group_.size
			calculated = true
		}
		inter_arrival.prev_timestamp_group_ = inter_arrival.current_timestamp_group_
		inter_arrival.current_timestamp_group_.first_timestamp = timestamp
		inter_arrival.current_timestamp_group_.timestamp = timestamp
		inter_arrival.current_timestamp_group_.first_arrival_ms = arrival_time_ms
		inter_arrival.current_timestamp_group_.size = 0
	}else{
		//只需修改timestamp
		//后面改成wrap ****
		inter_arrival.current_timestamp_group_.timestamp = timestamp
	}

	inter_arrival.current_timestamp_group_.size += packet_size
	inter_arrival.current_timestamp_group_.complete_time_ms = arrival_time_ms
	inter_arrival.current_timestamp_group_.last_system_time_ms = system_time_ms

	return calculated
}

func (inter_arrival *InterArrival)Reset(){
	inter_arrival.num_consecutive_reordered_packets_ = 0
	inter_arrival.current_timestamp_group_ = NewDefalutTimestampGroup()
	inter_arrival.prev_timestamp_group_ = NewDefalutTimestampGroup()
	fmt.Println("Reset inter_arrival")
}



