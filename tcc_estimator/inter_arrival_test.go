package tccEstimator

import (
	"fmt"
	"testing"
)

type packet struct {
	SendTime int64
	RecvTime int64
	FeedTime int64
	size int64
}

func TestArrive(t *testing.T) {
	inter_arriavl := NewInterArrival((5<<26)/1000, 1000.0/(1<<26), true)
	input := []packet{
			packet{SendTime:1416332544, RecvTime:1150293295, FeedTime:1150293280, size:23},
			packet{SendTime:1416399616, RecvTime:1150293295, FeedTime:1150293280, size:1096},
			packet{SendTime:1417070848, RecvTime:1150293295, FeedTime:1150293280, size:23},
			packet{SendTime:1417070848, RecvTime:1150293296, FeedTime:1150293280, size:1108},
			packet{SendTime:1417741824, RecvTime:1150293296, FeedTime:1150293280, size:1067},
			packet{SendTime:1418345728, RecvTime:1150293296, FeedTime:1150293280, size:1079},
			packet{SendTime:1419016960, RecvTime:1150293297, FeedTime:1150293280, size:1067},
			packet{SendTime:1419687936, RecvTime:1150293309, FeedTime:1150293280, size:973},
			packet{SendTime:1420023552, RecvTime:1150293309, FeedTime:1150293280, size:1079},
			packet{SendTime:1420292096, RecvTime:1150293309, FeedTime:1150293280, size:1165},
			packet{SendTime:1420627456, RecvTime:1150293309, FeedTime:1150293280, size:395},
			packet{SendTime:1420694528, RecvTime:1150293317, FeedTime:1150293280, size:1068},
			packet{SendTime:1421164288, RecvTime:1150293351, FeedTime:1150293280, size:1079},
			packet{SendTime:1422506496, RecvTime:1150293351, FeedTime:1150293280, size:751},
			packet{SendTime:1422573824, RecvTime:1150293351, FeedTime:1150293280, size:374},
			packet{SendTime:1423244800, RecvTime:1150293351, FeedTime:1150293280, size:771},
			packet{SendTime:1423311872, RecvTime:1150293351, FeedTime:1150293280, size:1068},
			packet{SendTime:1423915776, RecvTime:1150293408, FeedTime:1150293386, size:392},
			packet{SendTime:1423983104, RecvTime:1150293408, FeedTime:1150293386, size:1079},
			packet{SendTime:1424587008, RecvTime:1150293408, FeedTime:1150293386, size:1018},
			packet{SendTime:1424654080, RecvTime:1150293408, FeedTime:1150293386, size:23},
			packet{SendTime:1424654080, RecvTime:1150293408, FeedTime:1150293386, size:944},
			packet{SendTime:1425325056, RecvTime:1150293408, FeedTime:1150293386, size:1079},
			packet{SendTime:1425996288, RecvTime:1150293408, FeedTime:1150293386, size:415},
			packet{SendTime:1426063360, RecvTime:1150293408, FeedTime:1150293386, size:670},
			packet{SendTime:1426734336, RecvTime:1150293408, FeedTime:1150293386, size:945},
			packet{SendTime:1426801664, RecvTime:1150293409, FeedTime:1150293386, size:1079},
			packet{SendTime:1427472640, RecvTime:1150293415, FeedTime:1150293386, size:945},
			packet{SendTime:1428210944, RecvTime:1150293428, FeedTime:1150293386, size:405},
			packet{SendTime:1428278016, RecvTime:1150293428, FeedTime:1150293386, size:849},
			packet{SendTime:1428278016, RecvTime:1150293429, FeedTime:1150293386, size:1080},
			packet{SendTime:1428881920, RecvTime:1150293432, FeedTime:1150293386, size:367},
			packet{SendTime:1428948992, RecvTime:1150293432, FeedTime:1150293386, size:346},
			packet{SendTime:1430157056, RecvTime:1150293455, FeedTime:1150293386, size:1081},
			packet{SendTime:1430224128, RecvTime:1150293455, FeedTime:1150293386, size:1081},
			packet{SendTime:1430291200, RecvTime:1150293455, FeedTime:1150293386, size:1082},
			packet{SendTime:1430358272, RecvTime:1150293455, FeedTime:1150293386, size:1080},
			packet{SendTime:1430358272, RecvTime:1150293456, FeedTime:1150293386, size:364},
			packet{SendTime:1430425344, RecvTime:1150293456, FeedTime:1150293386, size:387},
			packet{SendTime:1430760960, RecvTime:1150293464, FeedTime:1150293386, size:578},
			packet{SendTime:1430828032, RecvTime:1150293464, FeedTime:1150293386, size:1000},
			packet{SendTime:1430828032, RecvTime:1150293464, FeedTime:1150293386, size:377},
			packet{SendTime:1430828032, RecvTime:1150293464, FeedTime:1150293386, size:24},
			packet{SendTime:1430895104, RecvTime:1150293465, FeedTime:1150293386, size:1137},
			packet{SendTime:1431432192, RecvTime:1150293477, FeedTime:1150293386, size:550},
			packet{SendTime:1431499264, RecvTime:1150293477, FeedTime:1150293386, size:1137},
			packet{SendTime:1431566336, RecvTime:1150293477, FeedTime:1150293386, size:1137},
			packet{SendTime:1432103168, RecvTime:1150293495, FeedTime:1150293386, size:1137},
			packet{SendTime:1432237312, RecvTime:1150293495, FeedTime:1150293386, size:1137},
			packet{SendTime:1432774144, RecvTime:1150293495, FeedTime:1150293386, size:645},
			packet{SendTime:1432841472, RecvTime:1150293495, FeedTime:1150293386, size:1017},
			packet{SendTime:1432841472, RecvTime:1150293495, FeedTime:1150293386, size:617},
			packet{SendTime:1433512448, RecvTime:1150293521, FeedTime:1150293481, size:723},
			packet{SendTime:1433579520, RecvTime:1150293522, FeedTime:1150293481, size:743},
			packet{SendTime:1433579520, RecvTime:1150293522, FeedTime:1150293481, size:24},
			packet{SendTime:1433646592, RecvTime:1150293522, FeedTime:1150293481, size:989},
			packet{SendTime:1433646592, RecvTime:1150293522, FeedTime:1150293481, size:989},
			packet{SendTime:1434250752, RecvTime:1150293522, FeedTime:1150293481, size:989},
			packet{SendTime:1434250752, RecvTime:1150293522, FeedTime:1150293481, size:989},
			packet{SendTime:1434317824, RecvTime:1150293522, FeedTime:1150293481, size:990},
			packet{SendTime:1434921728, RecvTime:1150293522, FeedTime:1150293481, size:601},
			packet{SendTime:1434988800, RecvTime:1150293523, FeedTime:1150293481, size:629},
			packet{SendTime:1434988800, RecvTime:1150293523, FeedTime:1150293481, size:627},
			packet{SendTime:1435055872, RecvTime:1150293523, FeedTime:1150293481, size:601},
			packet{SendTime:1436129792, RecvTime:1150293550, FeedTime:1150293481, size:990},
			packet{SendTime:1436263936, RecvTime:1150293550, FeedTime:1150293481, size:642},
			packet{SendTime:1436398080, RecvTime:1150293550, FeedTime:1150293481, size:821},
			packet{SendTime:1436599552, RecvTime:1150293550, FeedTime:1150293481, size:972},
			packet{SendTime:1436733696, RecvTime:1150293550, FeedTime:1150293481, size:989},
			packet{SendTime:1437136384, RecvTime:1150293561, FeedTime:1150293481, size:624},
			packet{SendTime:1437203456, RecvTime:1150293561, FeedTime:1150293481, size:623},
			packet{SendTime:1437203456, RecvTime:1150293561, FeedTime:1150293481, size:595},
			packet{SendTime:1437203456, RecvTime:1150293567, FeedTime:1150293481, size:596},
			packet{SendTime:1437270528, RecvTime:1150293593, FeedTime:1150293481, size:598},
			packet{SendTime:1437270528, RecvTime:1150293594, FeedTime:1150293481, size:599},
			packet{SendTime:1437337600, RecvTime:1150293594, FeedTime:1150293481, size:595},
			packet{SendTime:1437337600, RecvTime:1150293594, FeedTime:1150293481, size:595},
			packet{SendTime:1439216640, RecvTime:1150293594, FeedTime:1150293481, size:713},
			packet{SendTime:1439283712, RecvTime:1150293594, FeedTime:1150293481, size:685},
			packet{SendTime:1439283712, RecvTime:1150293595, FeedTime:1150293481, size:685},
			packet{SendTime:1441431296, RecvTime:1150293622, FeedTime:1150293590, size:638},
			packet{SendTime:1441498368, RecvTime:1150293622, FeedTime:1150293590, size:1043},
			packet{SendTime:1441498368, RecvTime:1150293622, FeedTime:1150293590, size:1015},
			packet{SendTime:1441565440, RecvTime:1150293622, FeedTime:1150293590, size:610},
			packet{SendTime:1442840576, RecvTime:1150293637, FeedTime:1150293590, size:610},
			packet{SendTime:1444249856, RecvTime:1150293664, FeedTime:1150293590, size:842},
			packet{SendTime:1444316928, RecvTime:1150293664, FeedTime:1150293590, size:798},
			packet{SendTime:1444316928, RecvTime:1150293664, FeedTime:1150293590, size:770},
			packet{SendTime:1444384000, RecvTime:1150293664, FeedTime:1150293590, size:814},
			packet{SendTime:1444988160, RecvTime:1150293677, FeedTime:1150293590, size:814},
			packet{SendTime:1446464512, RecvTime:1150293710, FeedTime:1150293687, size:976},
			packet{SendTime:1446531584, RecvTime:1150293710, FeedTime:1150293687, size:1036},
			packet{SendTime:1446598656, RecvTime:1150293710, FeedTime:1150293687, size:1008},
			packet{SendTime:1447135488, RecvTime:1150293710, FeedTime:1150293687, size:948},
			packet{SendTime:1447202560, RecvTime:1150293711, FeedTime:1150293687, size:948},
			packet{SendTime:1448544768, RecvTime:1150293723, FeedTime:1150293687, size:1148},
			packet{SendTime:1448611840, RecvTime:1150293723, FeedTime:1150293687, size:1120},
			packet{SendTime:1448679168, RecvTime:1150293734, FeedTime:1150293687, size:1120},
			packet{SendTime:1449618688, RecvTime:1150293754, FeedTime:1150293687, size:1013},
			packet{SendTime:1449685760, RecvTime:1150293754, FeedTime:1150293687, size:985},
	}

	for _, value := range input{
		timestamp_delta := uint32(0)
		recv_delta_ms := int64(0)
		size_delta := int64(0)

		valid := inter_arriavl.ComputeDeltas(uint32(value.SendTime),value.RecvTime, value.FeedTime,value.size, &timestamp_delta, &recv_delta_ms, &size_delta)
		fmt.Println(valid, " ", timestamp_delta, " ", recv_delta_ms, " ", size_delta)
	}
}
