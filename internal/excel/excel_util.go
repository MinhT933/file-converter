package excel

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func BuildHeaderIndex(header []string) map[string]int {
	indexMap := make(map[string]int)
	for i, h := range header {
		indexMap[h] = i
	}
	return indexMap
}

func MapRowByHeader(row []string, headerIndex map[string]int, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			continue
		}

		headerName, layout := parseExcelTag(tag)
		idx, ok := headerIndex[headerName]
		if !ok || idx >= len(row) {
			continue
		}

		cellValue := row[idx]
		if cellValue == "" {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		setFieldValue(fv, cellValue, layout)
	}
	return nil
}

// parseExcelTag extracts the header name and layout from the excel tag.
func parseExcelTag(tag string) (headerName, layout string) {
	parts := strings.Split(tag, ",")
	headerName = parts[0]
	layout = "2006-01-02"
	if len(parts) > 1 && strings.HasPrefix(parts[1], "layout=") {
		layout = strings.TrimPrefix(parts[1], "layout=")
	}
	return
}

// setFieldValue sets the value of a struct field based on its kind and the cell value.
func setFieldValue(fv reflect.Value, cellValue, layout string) {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(cellValue)
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setIntField(fv, cellValue)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setUintField(fv, cellValue)
		return
	case reflect.Float32, reflect.Float64:
		setFloatField(fv, cellValue)
		return
	case reflect.Complex64, reflect.Complex128:
		setComplexField(fv, cellValue)
		return
	case reflect.Bool:
		setBoolField(fv, cellValue)
		return
	default:
		if fv.Type() == reflect.TypeOf(time.Time{}) {
			setTimeField(fv, cellValue, layout)
		}
	}
}

// Helper functions to set specific kinds of fields
func setIntField(fv reflect.Value, cellValue string) {
	if iVal, err := strconv.Atoi(cellValue); err == nil {
		fv.SetInt(int64(iVal))
	}
}

// Helper functions to set specific kinds of fields
func setUintField(fv reflect.Value, cellValue string) {
	if uVal, err := strconv.ParseUint(cellValue, 10, 64); err == nil {
		fv.SetUint(uVal)
	}
}

// Helper functions to set specific kinds of fields
func setFloatField(fv reflect.Value, cellValue string) {
	if fVal, err := strconv.ParseFloat(cellValue, 64); err == nil {
		fv.SetFloat(fVal)
	}
}

// Helper functions to set specific kinds of fields
func setComplexField(fv reflect.Value, cellValue string) {
	if cVal, err := strconv.ParseComplex(cellValue, 128); err == nil {
		fv.SetComplex(cVal)
	}
}

// Helper functions to set specific kinds of fields
func setBoolField(fv reflect.Value, cellValue string) {
	if bVal, err := strconv.ParseBool(cellValue); err == nil {
		fv.SetBool(bVal)
	}
}

// Helper functions to set specific kinds of fields
func setTimeField(fv reflect.Value, cellValue, layout string) {
	if tVal, err := time.Parse(layout, cellValue); err == nil {
		fv.Set(reflect.ValueOf(tVal))
	}
}