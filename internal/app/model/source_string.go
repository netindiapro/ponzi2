// Code generated by "stringer -type=Source"; DO NOT EDIT.

package model

import "strconv"

const _Source_name = "SourceUnspecifiedSourceIEXRealTimePriceSource15MinuteDelayedPriceSourceCloseSourcePreviousClose"

var _Source_index = [...]uint8{0, 17, 39, 65, 76, 95}

func (i Source) String() string {
	if i < 0 || i >= Source(len(_Source_index)-1) {
		return "Source(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Source_name[_Source_index[i]:_Source_index[i+1]]
}
