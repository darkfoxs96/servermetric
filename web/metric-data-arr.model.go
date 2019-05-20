package web

import (
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
)

type MetricDataArr struct {
	buf  []byte
	Data [][]interface{}
}

func (m *MetricDataArr) UnmarshalJSON(data []byte) (err error) {
	m.buf = data
	return
}

func (m *MetricDataArr) UnmarshalBuf(types []string) (err error) {
	m.Data = make([][]interface{}, 0)

	_, err2 := jsonparser.ArrayEach(m.buf, func(value []byte, dataType jsonparser.ValueType, offset int, err3 error) {
		if err != nil {
			return
		}

		i := 0
		arr := make([]interface{}, len(types), len(types))

		_, err4 := jsonparser.ArrayEach(value, func(d []byte, dataType2 jsonparser.ValueType, offset2 int, err5 error) {
			if err != nil {
				return
			}

			arr[i], err = m.getObjByType(types[i], string(d))
			i++
		})

		if err4 != nil {
			err = err4
			return
		}

		m.Data = append(m.Data, arr)
	})

	if err2 != nil {
		err = err2
	}

	return
}

func (m *MetricDataArr) getObjByType(t, data string) (obj interface{}, err error) {
	switch t {
	case "string":
		return data, nil
		// int
	case "int":
		return strconv.Atoi(data)
	case "int8":
		v, err := strconv.ParseInt(data, 10, 8)
		return int8(v), err
	case "int16":
		v, err := strconv.ParseInt(data, 10, 16)
		return int16(v), err
	case "int32":
		v, err := strconv.ParseInt(data, 10, 32)
		return int32(v), err
	case "int64":
		return strconv.ParseInt(data, 10, 64)
		// uint
	case "uint":
		v, err := strconv.ParseUint(data, 10, 64)
		return uint(v), err
	case "uint8":
		v, err := strconv.ParseUint(data, 10, 8)
		return uint8(v), err
	case "uint16":
		v, err := strconv.ParseUint(data, 10, 16)
		return uint16(v), err
	case "uint32":
		v, err := strconv.ParseUint(data, 10, 32)
		return uint32(v), err
	case "uint64":
		return strconv.ParseUint(data, 10, 64)
		// float
	case "float32":
		v, err := strconv.ParseFloat(data, 32)
		return float32(v), err
	case "float64":
		return strconv.ParseFloat(data, 64)
		// bool
	case "bool":
		return strconv.ParseBool(data)
	}

	return nil, fmt.Errorf("Metric-api: getObjByType() don't supported type: " + t)
}
