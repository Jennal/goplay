// Code generated by "stringer -type=Status"; DO NOT EDIT.

package pkg

import "fmt"

const (
	_Status_name_0 = "STAT_OK"
	_Status_name_1 = "STAT_ERRSTAT_ERR_WRONG_PARAMSSTAT_ERR_DECODE_FAILEDSTAT_ERR_TIMEOUTSTAT_ERR_EMPTY_RESULT"
)

var (
	_Status_index_0 = [...]uint8{0, 7}
	_Status_index_1 = [...]uint8{0, 8, 29, 51, 67, 88}
)

func (i Status) String() string {
	switch {
	case i == 0:
		return _Status_name_0
	case 144 <= i && i <= 148:
		i -= 144
		return _Status_name_1[_Status_index_1[i]:_Status_index_1[i+1]]
	default:
		return fmt.Sprintf("Status(%d)", i)
	}
}
