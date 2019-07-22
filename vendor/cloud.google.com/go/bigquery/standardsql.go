// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquery

import (
	"fmt"

	bq "google.golang.org/api/bigquery/v2"
)

// StandardSQLDataType conveys type information using the Standard SQL type
// system.
type StandardSQLDataType struct {
	// ArrayElementType indicates the type of an array's elements, when the
	// TypeKind is ARRAY.
	ArrayElementType *StandardSQLDataType
	// StructType indicates the struct definition (fields), when the
	// TypeKind is STRUCT.
	StructType *StandardSQLStructType
	// The top-level type of this type definition.
	// Can be any standard SQL data type.  For more information about BigQuery
	// data types, see
	// https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types
	//
	// Additional information is available in the REST documentation:
	// https://cloud.google.com/bigquery/docs/reference/rest/v2/StandardSqlDataType
	TypeKind string
}

func (ssdt *StandardSQLDataType) toBQ() (*bq.StandardSqlDataType, error) {
	if ssdt == nil {
		return nil, nil
	}
	bqdt := &bq.StandardSqlDataType{
		TypeKind: ssdt.TypeKind,
	}
	if ssdt.ArrayElementType != nil {
		dt, err := ssdt.ArrayElementType.toBQ()
		if err != nil {
			return nil, err
		}
		bqdt.ArrayElementType = dt
	}
	if ssdt.StructType != nil {
		dt, err := ssdt.StructType.toBQ()
		if err != nil {
			return nil, err
		}
		bqdt.StructType = dt
	}
	return bqdt, nil
}

func bqToStandardSQLDataType(bqdt *bq.StandardSqlDataType) (*StandardSQLDataType, error) {
	if bqdt == nil {
		return nil, nil
	}
	ssdt := &StandardSQLDataType{
		TypeKind: bqdt.TypeKind,
	}

	if bqdt.ArrayElementType != nil {
		dt, err := bqToStandardSQLDataType(bqdt.ArrayElementType)
		if err != nil {
			return nil, err
		}
		ssdt.ArrayElementType = dt
	}
	if bqdt.StructType != nil {
		st, err := bqToStandardSQLStructType(bqdt.StructType)
		if err != nil {
			return nil, err
		}
		ssdt.StructType = st
	}
	return ssdt, nil
}

// StandardSQLField represents a field using the Standard SQL data type system.
type StandardSQLField struct {
	// The name of this field.  Can be absent for struct fields.
	Name string
	// Data type for the field.
	Type *StandardSQLDataType
}

func (ssf *StandardSQLField) toBQ() (*bq.StandardSqlField, error) {
	if ssf == nil {
		return nil, nil
	}
	bqf := &bq.StandardSqlField{
		Name: ssf.Name,
	}
	if ssf.Type != nil {
		dt, err := ssf.Type.toBQ()
		if err != nil {
			return nil, err
		}
		bqf.Type = dt
	}
	return bqf, nil
}

func bqToStandardSQLField(bqf *bq.StandardSqlField) (*StandardSQLField, error) {
	if bqf == nil {
		return nil, nil
	}
	t, err := bqToStandardSQLDataType(bqf.Type)
	if err != nil {
		return nil, err
	}
	return &StandardSQLField{
		Name: bqf.Name,
		Type: t,
	}, nil
}

// StandardSQLStructType represents a structure type, which is a list of Standard SQL fields.
// For more information, see:
// https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#struct-type
type StandardSQLStructType struct {
	Fields []*StandardSQLField
}

func (ssst *StandardSQLStructType) toBQ() (*bq.StandardSqlStructType, error) {
	if ssst == nil {
		return nil, nil
	}
	fields, err := standardSQLStructFieldsToBQ(ssst.Fields)
	if err != nil {
		return nil, err
	}
	return &bq.StandardSqlStructType{
		Fields: fields,
	}, nil
}

func bqToStandardSQLStructType(bqst *bq.StandardSqlStructType) (*StandardSQLStructType, error) {
	if bqst == nil {
		return nil, nil
	}
	var fields []*StandardSQLField
	for _, v := range bqst.Fields {
		f, err := bqToStandardSQLField(v)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return &StandardSQLStructType{
		Fields: fields,
	}, nil
}

func standardSQLStructFieldsToBQ(fields []*StandardSQLField) ([]*bq.StandardSqlField, error) {
	var bqFields []*bq.StandardSqlField
	for _, v := range fields {
		bqf, err := v.toBQ()
		if err != nil {
			return nil, fmt.Errorf("error converting struct fields: %v", err)
		}
		bqFields = append(bqFields, bqf)
	}
	return bqFields, nil
}
