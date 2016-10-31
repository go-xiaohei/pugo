package model

import "github.com/tidwall/gjson"

// JSON simple reader struct
type JSON struct {
	data   []byte
	result gjson.Result
}

// NewJSON create JSON struct with json bytes
func NewJSON(data []byte) *JSON {
	j := &JSON{
		data:   data,
		result: gjson.Parse(string(data)),
	}
	return j
}

// NewJSONwithResult create JSON struct with parsed json result
func NewJSONwithResult(res gjson.Result) *JSON {
	return &JSON{
		result: res,
	}
}

// Get return a json struct to operation this json node
func (j *JSON) Get(path ...string) *JSON {
	var res gjson.Result
	if len(path) == 0 || path[0] == "" {
		res = j.result
	} else {
		res = j.result.Get(path[0])
	}
	return NewJSONwithResult(res)
}

// String get string value by path string
func (j *JSON) String(path ...string) string {
	if len(path) == 0 || path[0] == "" {
		return j.result.String()
	}
	res := j.result.Get(path[0])
	if !res.Exists() {
		return ""
	}
	return res.String()
}

// Int64 get int64 value by path
func (j *JSON) Int64(path ...string) int64 {
	if len(path) == 0 || path[0] == "" {
		return j.result.Int()
	}
	res := j.result.Get(path[0])
	if !res.Exists() {
		return 0
	}
	return res.Int()
}

// Int32 get int32 value by path
func (j *JSON) Int32(path ...string) int32 {
	return int32(j.Int64(path...))
}

// Int16 get int16 value by path
func (j *JSON) Int16(path ...string) int16 {
	return int16(j.Int64(path...))
}

// Int8 get int8 value by path
func (j *JSON) Int8(path ...string) int8 {
	return int8(j.Int64(path...))
}

// Int get int value by path
func (j *JSON) Int(path ...string) int {
	return int(j.Int64(path...))
}

// Float64 get float64 value by path
func (j *JSON) Float64(path ...string) float64 {
	if len(path) == 0 || path[0] == "" {
		return j.result.Float()
	}
	res := j.result.Get(path[0])
	if !res.Exists() {
		return 0
	}
	return res.Float()
}

// Float32 get float32 value by path
func (j *JSON) Float32(path ...string) float32 {
	return float32(j.Float64(path...))
}

// Float is alias of float64
func (j *JSON) Float(path ...string) float64 {
	return j.Float64(path...)
}

// Bool get bool value by path
func (j *JSON) Bool(path ...string) bool {
	if len(path) == 0 || path[0] == "" {
		return j.result.Bool()
	}
	res := j.result.Get(path[0])
	if !res.Exists() {
		return false
	}
	return res.Bool()
}

// Exist check existing of value by path
func (j *JSON) Exist(path ...string) bool {
	if len(path) == 0 || path[0] == "" {
		return j.result.Exists()
	}
	return j.result.Get(path[0]).Exists()
}

// Strings get string slice by path
func (j *JSON) Strings(path ...string) []string {
	var data []gjson.Result
	if len(path) == 0 || path[0] == "" {
		data = j.result.Array()
	} else {
		data = j.result.Get(path[0]).Array()
	}
	if len(data) == 0 {
		return nil
	}
	resData := make([]string, len(data))
	for i, r := range data {
		resData[i] = r.String()
	}
	return resData
}

// Ints get int64 slice by path
func (j *JSON) Ints(path ...string) []int64 {
	var data []gjson.Result
	if len(path) == 0 || path[0] == "" {
		data = j.result.Array()
	} else {
		data = j.result.Get(path[0]).Array()
	}
	if len(data) == 0 {
		return nil
	}
	resData := make([]int64, len(data))
	for i, r := range data {
		resData[i] = r.Int()
	}
	return resData
}

// Floats get float64 slice by path
func (j *JSON) Floats(path ...string) []float64 {
	var data []gjson.Result
	if len(path) == 0 || path[0] == "" {
		data = j.result.Array()
	} else {
		data = j.result.Get(path[0]).Array()
	}
	if len(data) == 0 {
		return nil
	}
	resData := make([]float64, len(data))
	for i, r := range data {
		resData[i] = r.Float()
	}
	return resData
}

// Slice get JSON struct slice by path
func (j *JSON) Slice(path ...string) []*JSON {
	var data []gjson.Result
	if len(path) == 0 || path[0] == "" {
		data = j.result.Array()
	} else {
		data = j.result.Get(path[0]).Array()
	}
	if len(data) == 0 {
		return nil
	}
	resData := make([]*JSON, len(data))
	for i, r := range data {
		resData[i] = NewJSONwithResult(r)
	}
	return resData
}

// Index get item of slice by index if current JSON is an array,
// otherwise return nil
func (j *JSON) Index(i int) *JSON {
	data := j.result.Array()
	if len(data) == 0 {
		return nil
	}
	if i >= len(data) || i < 0 {
		return nil
	}
	return NewJSONwithResult(data[i])
}

// Map get map of JSON by path
func (j *JSON) Map(path ...string) map[string]*JSON {
	var data map[string]gjson.Result
	if len(path) == 0 || path[0] == "" {
		data = j.result.Map()
	} else {
		data = j.result.Get(path[0]).Map()
	}
	if len(data) == 0 {
		return nil
	}
	res := make(map[string]*JSON)
	for k, r := range data {
		res[k] = NewJSONwithResult(r)
	}
	return res
}

// Key get JSON struct by key if current JSON is map
func (j *JSON) Key(key string) *JSON {
	res := j.result.Get(key)
	if !res.Exists() {
		return nil
	}
	return NewJSONwithResult(res)
}
