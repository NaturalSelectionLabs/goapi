// Code generated by "enumer -type=ErrCode -values -transform=snake -trimprefix=ErrCode -json"; DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _ErrCodeName = "not_foundinvalid_paraminternal_error"

var _ErrCodeIndex = [...]uint8{0, 9, 22, 36}

const _ErrCodeLowerName = "not_foundinvalid_paraminternal_error"

func (i Code) String() string {
	if i < 0 || i >= Code(len(_ErrCodeIndex)-1) {
		return fmt.Sprintf("ErrCode(%d)", i)
	}
	return _ErrCodeName[_ErrCodeIndex[i]:_ErrCodeIndex[i+1]]
}

func (Code) Values() []string {
	return ErrCodeStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ErrCodeNoOp() {
	var x [1]struct{}
	_ = x[CodeNotFound-(0)]
	_ = x[CodeInvalidParam-(1)]
	_ = x[CodeInternalError-(2)]
}

var _ErrCodeValues = []Code{CodeNotFound, CodeInvalidParam, CodeInternalError}

var _ErrCodeNameToValueMap = map[string]Code{
	_ErrCodeName[0:9]:        CodeNotFound,
	_ErrCodeLowerName[0:9]:   CodeNotFound,
	_ErrCodeName[9:22]:       CodeInvalidParam,
	_ErrCodeLowerName[9:22]:  CodeInvalidParam,
	_ErrCodeName[22:36]:      CodeInternalError,
	_ErrCodeLowerName[22:36]: CodeInternalError,
}

var _ErrCodeNames = []string{
	_ErrCodeName[0:9],
	_ErrCodeName[9:22],
	_ErrCodeName[22:36],
}

// ErrCodeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ErrCodeString(s string) (Code, error) {
	if val, ok := _ErrCodeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ErrCodeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ErrCode values", s)
}

// ErrCodeValues returns all values of the enum
func ErrCodeValues() []Code {
	return _ErrCodeValues
}

// ErrCodeStrings returns a slice of all String values of the enum
func ErrCodeStrings() []string {
	strs := make([]string, len(_ErrCodeNames))
	copy(strs, _ErrCodeNames)
	return strs
}

// IsAErrCode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Code) IsAErrCode() bool {
	for _, v := range _ErrCodeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ErrCode
func (i Code) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ErrCode
func (i *Code) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ErrCode should be a string, got %s", data)
	}

	var err error
	*i, err = ErrCodeString(s)
	return err
}