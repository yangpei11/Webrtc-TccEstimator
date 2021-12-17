package tccEstimator

import (
	"math"
)


type LinkCapacityEstimator struct {
	estimate_kbps_	float64
	deviation_kbps_ float64
}

func (link_capacity_estimator *LinkCapacityEstimator)has_estimate()bool{
	if(math.Abs(-1-link_capacity_estimator.estimate_kbps_) < 0.1){
		return false
	}else{
		return true
	}
}

func (link_capacity_estimator *LinkCapacityEstimator)UpperBound()int64{
	if(link_capacity_estimator.has_estimate()){
		return int64((link_capacity_estimator.estimate_kbps_ + 3 * link_capacity_estimator.deviation_estimate_kbps())*1000)
	}
	return math.MaxInt64
}

func (link_capacity_estimator *LinkCapacityEstimator)deviation_estimate_kbps()float64{
	return math.Sqrt(link_capacity_estimator.deviation_kbps_ * link_capacity_estimator.estimate_kbps_)
}

func (link_capacity_estimator *LinkCapacityEstimator)Reset(){
	link_capacity_estimator.estimate_kbps_ = -1
}

func (link_capacity_estimator *LinkCapacityEstimator)LowerBound()int64{
	if(link_capacity_estimator.has_estimate()){
		return int64(math.Max(0.0, link_capacity_estimator.estimate_kbps_- 3*link_capacity_estimator.deviation_estimate_kbps())*1000)
	}
	return 0
}

func (link_capacity_estimator *LinkCapacityEstimator)Update(capacity_sample int64, alpha float64){
	sample_kbps := math.Round(float64(capacity_sample)/1000.0)
	if(!link_capacity_estimator.has_estimate()){
		link_capacity_estimator.estimate_kbps_ = sample_kbps
	}else{
		link_capacity_estimator.estimate_kbps_ = (1-alpha)*link_capacity_estimator.estimate_kbps_ + alpha * sample_kbps
	}

	norm := math.Max(link_capacity_estimator.estimate_kbps_, 1.0)
	error_kbps := link_capacity_estimator.estimate_kbps_ - float64(sample_kbps)
	link_capacity_estimator.deviation_kbps_ = (1-alpha)*link_capacity_estimator.deviation_kbps_ + alpha * error_kbps*error_kbps/norm
	if(link_capacity_estimator.deviation_kbps_ < 0.4){
		link_capacity_estimator.deviation_kbps_ = 0.4
	}else if(link_capacity_estimator.deviation_kbps_ > 2.5){
		link_capacity_estimator.deviation_kbps_ = 2.5
	}

}

func (link_capacity_estimator *LinkCapacityEstimator)OnOveruseDetected(acknowledged_rate int64){
	link_capacity_estimator.Update(acknowledged_rate, 0.05)
}

