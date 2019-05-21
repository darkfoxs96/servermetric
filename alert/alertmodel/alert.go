package alertmodel

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
)

// Errors
var (
	ErrManyOutParams    = fmt.Errorf("Alert: many out from DB parameters (max 10)")
	ErrIfEmpty          = fmt.Errorf("Alert: parameter if is empty")
	ErrThenAndElseEmpty = fmt.Errorf("Alert: parameter then and else is empty")
)

var (
	temFunc = template.FuncMap{
		"sum": func(v1, v2 int64) int64 {
			return v1 + v2
		},
		"mul": func(v1, v2 int64) int64 {
			return v1 * v2
		},
		"sub": func(v1, v2 int64) int64 {
			return v1 - v2
		},
	}
)

type OutDB struct {
	V1  interface{}
	V2  interface{}
	V3  interface{}
	V4  interface{}
	V5  interface{}
	V6  interface{}
	V7  interface{}
	V8  interface{}
	V9  interface{}
	V10 interface{}
}

type Alert struct {
	IF   *template.Template
	THEN *template.Template
	ELSE *template.Template
}

// NewAlert crete new Alert{}
func NewAlert(ifTmp, thenTmp, elseTmp string) (alert *Alert, err error) {
	alert = &Alert{}

	if ifTmp == "" {
		err = ErrIfEmpty
		return
	}

	if thenTmp == "" && elseTmp == "" {
		err = ErrThenAndElseEmpty
		return
	}

	alert.IF, err = template.New("if").Funcs(temFunc).Parse(ifTmp)
	if err != nil {
		return
	}

	alert.THEN, err = template.New("then").Funcs(temFunc).Parse(thenTmp)
	if err != nil {
		return
	}

	alert.ELSE, err = template.New("else").Funcs(temFunc).Parse(elseTmp)
	if err != nil {
		return
	}

	return
}

// Check metric
func (a *Alert) Check(DB *sql.DB, data interface{}) (msgTHEN, msgELSE []string, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = a.IF.Execute(buf, data)
	if err != nil {
		return
	}

	sqlReq := buf.String()
	row, err := DB.Query(sqlReq)
	if err != nil {
		return
	}

	columns, err := row.Columns()
	if err != nil {
		return
	}

	isTrue := false
	msgTHEN = make([]string, 0)
	msgELSE = make([]string, 0)
	out := OutDB{}
	cLen := len(columns)
	var forScan []interface{}

	for row.Next() {
		isTrue = true

		switch cLen {
		case 0:
			forScan = []interface{}{} // TODO: add action
		case 1:
			forScan = []interface{}{&out.V1}
		case 2:
			forScan = []interface{}{&out.V1, &out.V2}
		case 3:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3}
		case 4:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4}
		case 5:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5}
		case 6:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5, &out.V6}
		case 7:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5, &out.V6, &out.V7}
		case 8:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5, &out.V6, &out.V7, &out.V8}
		case 9:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5, &out.V6, &out.V7, &out.V8, &out.V9}
		case 10:
			forScan = []interface{}{&out.V1, &out.V2, &out.V3, &out.V4, &out.V5, &out.V6, &out.V7, &out.V8, &out.V9, &out.V10}
		default:
			err = ErrManyOutParams
			return
		}

		err = row.Scan(forScan...)
		if err != nil {
			return
		}

		buf.Reset()
		err = a.THEN.Execute(buf, out)
		if err != nil {
			return
		}
		if msg := buf.String(); msg != "" {
			msgTHEN = append(msgTHEN, msg)
		}

		out.V1 = nil
		out.V2 = nil
		out.V3 = nil
		out.V4 = nil
		out.V5 = nil
		out.V6 = nil
		out.V7 = nil
		out.V8 = nil
		out.V9 = nil
		out.V10 = nil
	}
	_ = row.Close()

	if !isTrue {
		buf.Reset()
		err = a.ELSE.Execute(buf, nil)
		if err != nil {
			return
		}

		if msg := buf.String(); msg != "" {
			msgELSE = append(msgELSE, msg)
		}
	}

	return
}
