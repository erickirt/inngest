// Code generated by "enumer -trimprefix=WorkerConnectionSortOrder -type=WorkerConnectionSortOrder -json -text"; DO NOT EDIT.

package enums

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _WorkerConnectionSortOrderName = "DescAsc"

var _WorkerConnectionSortOrderIndex = [...]uint8{0, 4, 7}

const _WorkerConnectionSortOrderLowerName = "descasc"

func (i WorkerConnectionSortOrder) String() string {
	if i < 0 || i >= WorkerConnectionSortOrder(len(_WorkerConnectionSortOrderIndex)-1) {
		return fmt.Sprintf("WorkerConnectionSortOrder(%d)", i)
	}
	return _WorkerConnectionSortOrderName[_WorkerConnectionSortOrderIndex[i]:_WorkerConnectionSortOrderIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _WorkerConnectionSortOrderNoOp() {
	var x [1]struct{}
	_ = x[WorkerConnectionSortOrderDesc-(0)]
	_ = x[WorkerConnectionSortOrderAsc-(1)]
}

var _WorkerConnectionSortOrderValues = []WorkerConnectionSortOrder{WorkerConnectionSortOrderDesc, WorkerConnectionSortOrderAsc}

var _WorkerConnectionSortOrderNameToValueMap = map[string]WorkerConnectionSortOrder{
	_WorkerConnectionSortOrderName[0:4]:      WorkerConnectionSortOrderDesc,
	_WorkerConnectionSortOrderLowerName[0:4]: WorkerConnectionSortOrderDesc,
	_WorkerConnectionSortOrderName[4:7]:      WorkerConnectionSortOrderAsc,
	_WorkerConnectionSortOrderLowerName[4:7]: WorkerConnectionSortOrderAsc,
}

var _WorkerConnectionSortOrderNames = []string{
	_WorkerConnectionSortOrderName[0:4],
	_WorkerConnectionSortOrderName[4:7],
}

// WorkerConnectionSortOrderString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func WorkerConnectionSortOrderString(s string) (WorkerConnectionSortOrder, error) {
	if val, ok := _WorkerConnectionSortOrderNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _WorkerConnectionSortOrderNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to WorkerConnectionSortOrder values", s)
}

// WorkerConnectionSortOrderValues returns all values of the enum
func WorkerConnectionSortOrderValues() []WorkerConnectionSortOrder {
	return _WorkerConnectionSortOrderValues
}

// WorkerConnectionSortOrderStrings returns a slice of all String values of the enum
func WorkerConnectionSortOrderStrings() []string {
	strs := make([]string, len(_WorkerConnectionSortOrderNames))
	copy(strs, _WorkerConnectionSortOrderNames)
	return strs
}

// IsAWorkerConnectionSortOrder returns "true" if the value is listed in the enum definition. "false" otherwise
func (i WorkerConnectionSortOrder) IsAWorkerConnectionSortOrder() bool {
	for _, v := range _WorkerConnectionSortOrderValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for WorkerConnectionSortOrder
func (i WorkerConnectionSortOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for WorkerConnectionSortOrder
func (i *WorkerConnectionSortOrder) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("WorkerConnectionSortOrder should be a string, got %s", data)
	}

	var err error
	*i, err = WorkerConnectionSortOrderString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for WorkerConnectionSortOrder
func (i WorkerConnectionSortOrder) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for WorkerConnectionSortOrder
func (i *WorkerConnectionSortOrder) UnmarshalText(text []byte) error {
	var err error
	*i, err = WorkerConnectionSortOrderString(string(text))
	return err
}
