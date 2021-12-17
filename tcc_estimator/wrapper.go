package tccEstimator

import (
	"math"
)

type SequenceNumberUnwrapper struct{
	last_value_ int64
}

func IsNewer(value uint16, prev_value uint16)bool{
	kBreakpoint := uint16((math.MaxUint16 >> 1) + 1)
	if(value - prev_value == kBreakpoint){
		return value>prev_value
	}

	delta := value-prev_value
	return value != prev_value && delta < kBreakpoint
}

func NewSequenceNumberUnwrapper()SequenceNumberUnwrapper{
	return SequenceNumberUnwrapper{last_value_: -1}
}

func (sequence_number_unwrapper *SequenceNumberUnwrapper)Unwrap(value uint16)int64{
	if(sequence_number_unwrapper.last_value_ == -1){
		sequence_number_unwrapper.last_value_ = int64(value)
		return int64(value)
	}

	kMaxPlusOne := int64(math.MaxUint16) + 1
	cropped_last := uint16(sequence_number_unwrapper.last_value_)

	delta := int64(value)-int64(cropped_last)

	if(IsNewer(value, cropped_last)){
		if(delta < 0){
			delta += kMaxPlusOne
		}
	}else if( (delta>0) && (sequence_number_unwrapper.last_value_+delta-kMaxPlusOne) >= 0){
		delta -= kMaxPlusOne
	}

	sequence_number_unwrapper.last_value_ +=  delta
	return sequence_number_unwrapper.last_value_
}




