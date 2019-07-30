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
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/internal/optional"
	"cloud.google.com/go/internal/trace"
	bq "google.golang.org/api/bigquery/v2"
)

// Routine represents a reference to a BigQuery routine.  There are multiple
// types of routines including stored procedures and scalar user-defined functions (UDFs).
// For more information, see the BigQuery documentation at https://cloud.google.com/bigquery/docs/
type Routine struct {
	ProjectID string
	DatasetID string
	RoutineID string

	c *Client
}

// FullyQualifiedName returns an identifer for the routine in project.dataset.routine format.
func (r *Routine) FullyQualifiedName() string {
	return fmt.Sprintf("%s.%s.%s", r.ProjectID, r.DatasetID, r.RoutineID)
}

// Create creates a Routine in the BigQuery service.
// Pass in a RoutineMetadata to define the routine.
func (r *Routine) Create(ctx context.Context, rm *RoutineMetadata) (err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Routine.Create")
	defer func() { trace.EndSpan(ctx, err) }()

	routine, err := rm.toBQ()
	if err != nil {
		return err
	}
	routine.RoutineReference = &bq.RoutineReference{
		ProjectId: r.ProjectID,
		DatasetId: r.DatasetID,
		RoutineId: r.RoutineID,
	}
	req := r.c.bqs.Routines.Insert(r.ProjectID, r.DatasetID, routine).Context(ctx)
	setClientHeader(req.Header())
	_, err = req.Do()
	return err
}

// Metadata fetches the metadata for a given Routine.
func (r *Routine) Metadata(ctx context.Context) (rm *RoutineMetadata, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Routine.Metadata")
	defer func() { trace.EndSpan(ctx, err) }()

	req := r.c.bqs.Routines.Get(r.ProjectID, r.DatasetID, r.RoutineID).Context(ctx)
	setClientHeader(req.Header())
	var routine *bq.Routine
	err = runWithRetry(ctx, func() (err error) {
		routine, err = req.Do()
		return err
	})
	if err != nil {
		return nil, err
	}
	return bqToRoutineMetadata(routine)
}

// Update modifies properties of a Routine using the API.
func (r *Routine) Update(ctx context.Context, upd *RoutineMetadataToUpdate, etag string) (rm *RoutineMetadata, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Routine.Update")
	defer func() { trace.EndSpan(ctx, err) }()

	bqr, err := upd.toBQ()
	if err != nil {
		return nil, err
	}
	//TODO: remove when routines update supports partial requests.
	bqr.RoutineReference = &bq.RoutineReference{
		ProjectId: r.ProjectID,
		DatasetId: r.DatasetID,
		RoutineId: r.RoutineID,
	}

	call := r.c.bqs.Routines.Update(r.ProjectID, r.DatasetID, r.RoutineID, bqr).Context(ctx)
	setClientHeader(call.Header())
	if etag != "" {
		call.Header().Set("If-Match", etag)
	}
	var res *bq.Routine
	if err := runWithRetry(ctx, func() (err error) {
		res, err = call.Do()
		return err
	}); err != nil {
		return nil, err
	}
	return bqToRoutineMetadata(res)
}

// Delete removes a Routine from a dataset.
func (r *Routine) Delete(ctx context.Context) (err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Model.Delete")
	defer func() { trace.EndSpan(ctx, err) }()

	req := r.c.bqs.Routines.Delete(r.ProjectID, r.DatasetID, r.RoutineID).Context(ctx)
	setClientHeader(req.Header())
	return req.Do()
}

// RoutineMetadata represents details of a given BigQuery Routine.
type RoutineMetadata struct {
	ETag string
	// Type indicates the type of routine, such as SCALAR_FUNCTION or PROCEDURE.
	Type             string
	CreationTime     time.Time
	LastModifiedTime time.Time
	// Language of the routine, such as SQL or JAVASCRIPT.
	Language string
	// The list of arguments for the the routine.
	Arguments  []*RoutineArgument
	ReturnType *StandardSQLDataType
	// For javascript routines, this indicates the paths for imported libraries.
	ImportedLibraries []string
	// Body contains the routine's body.
	// For functions, Body is the expression in the AS clause.
	//
	// For SQL functions, it is the substring inside the parentheses of a CREATE
	// FUNCTION statement.
	//
	// For JAVASCRIPT function, it is the evaluated string in the AS clause of
	// a CREATE FUNCTION statement.
	Body string
}

func (rm *RoutineMetadata) toBQ() (*bq.Routine, error) {
	r := &bq.Routine{}
	if rm == nil {
		return r, nil
	}
	r.Language = rm.Language
	r.RoutineType = rm.Type
	r.DefinitionBody = rm.Body

	var args []*bq.Argument
	for _, v := range rm.Arguments {
		bqa, err := v.toBQ()
		if err != nil {
			return nil, err
		}
		args = append(args, bqa)
	}
	r.Arguments = args
	r.ImportedLibraries = rm.ImportedLibraries
	if !rm.CreationTime.IsZero() {
		return nil, errors.New("cannot set CreationTime on create")
	}
	if !rm.LastModifiedTime.IsZero() {
		return nil, errors.New("cannot set LastModifiedTime on create")
	}
	if rm.ETag != "" {
		return nil, errors.New("cannot set ETag on create")
	}
	return r, nil
}

// RoutineArgument represents an argument supplied to a routine such as a UDF or
// stored procedured.
type RoutineArgument struct {
	// The name of this argument.  Can be absent for function return argument.
	Name string
	// Kind indicates the kind of argument represented.
	// Possible values:
	//   ARGUMENT_KIND_UNSPECIFIED
	//   FIXED_TYPE - The argument is a variable with fully specified
	//     type, which can be a struct or an array, but not a table.
	//   ANY_TYPE - The argument is any type, including struct or array,
	//     but not a table.
	Kind string
	// Mode is optional, and indicates whether an argument is input or output.
	// Mode can only be set for procedures.
	//
	// Possible values:
	//   MODE_UNSPECIFIED
	//   IN - The argument is input-only.
	//   OUT - The argument is output-only.
	//   INOUT - The argument is both an input and an output.
	Mode string
	// DataType provides typing information.  Unnecessary for ANY_TYPE Kind
	// arguments.
	DataType *StandardSQLDataType
}

func (ra *RoutineArgument) toBQ() (*bq.Argument, error) {
	if ra == nil {
		return nil, nil
	}
	a := &bq.Argument{
		Name:         ra.Name,
		ArgumentKind: ra.Kind,
		Mode:         ra.Mode,
	}
	if ra.DataType != nil {
		dt, err := ra.DataType.toBQ()
		if err != nil {
			return nil, err
		}
		a.DataType = dt
	}
	return a, nil
}

func bqToRoutineArgument(bqa *bq.Argument) (*RoutineArgument, error) {
	arg := &RoutineArgument{
		Name: bqa.Name,
		Kind: bqa.ArgumentKind,
		Mode: bqa.Mode,
	}
	dt, err := bqToStandardSQLDataType(bqa.DataType)
	if err != nil {
		return nil, err
	}
	arg.DataType = dt
	return arg, nil
}

func bqToArgs(in []*bq.Argument) ([]*RoutineArgument, error) {
	var out []*RoutineArgument
	for _, a := range in {
		arg, err := bqToRoutineArgument(a)
		if err != nil {
			return nil, err
		}
		out = append(out, arg)
	}
	return out, nil
}

func routineArgumentsToBQ(in []*RoutineArgument) ([]*bq.Argument, error) {
	var out []*bq.Argument
	for _, inarg := range in {
		arg, err := inarg.toBQ()
		if err != nil {
			return nil, err
		}
		out = append(out, arg)
	}
	return out, nil
}

// RoutineMetadataToUpdate governs updating a routine.
type RoutineMetadataToUpdate struct {
	Arguments         []*RoutineArgument
	Type              optional.String
	Language          optional.String
	Body              optional.String
	ImportedLibraries []string
	ReturnType        *StandardSQLDataType
}

func (rm *RoutineMetadataToUpdate) toBQ() (*bq.Routine, error) {
	r := &bq.Routine{}
	forceSend := func(field string) {
		r.ForceSendFields = append(r.ForceSendFields, field)
	}
	nullField := func(field string) {
		r.NullFields = append(r.NullFields, field)
	}
	if rm.Arguments != nil {
		if len(rm.Arguments) == 0 {
			nullField("Arguments")
		} else {
			args, err := routineArgumentsToBQ(rm.Arguments)
			if err != nil {
				return nil, err
			}
			r.Arguments = args
			forceSend("Arguments")
		}
	}
	if rm.Type != nil {
		r.RoutineType = optional.ToString(rm.Type)
		forceSend("RoutineType")
	}
	if rm.Language != nil {
		r.Language = optional.ToString(rm.Language)
		forceSend("Language")
	}
	if rm.Body != nil {
		r.DefinitionBody = optional.ToString(rm.Body)
		forceSend("DefinitionBody")
	}
	if rm.ImportedLibraries != nil {
		if len(rm.ImportedLibraries) == 0 {
			nullField("ImportedLibraries")
		} else {
			r.ImportedLibraries = rm.ImportedLibraries
			forceSend("ImportedLibraries")
		}
	}
	if rm.ReturnType != nil {
		dt, err := rm.ReturnType.toBQ()
		if err != nil {
			return nil, err
		}
		r.ReturnType = dt
		forceSend("ReturnType")
	}
	return r, nil
}

func bqToRoutineMetadata(r *bq.Routine) (*RoutineMetadata, error) {
	meta := &RoutineMetadata{
		ETag:              r.Etag,
		Type:              r.RoutineType,
		CreationTime:      unixMillisToTime(r.CreationTime),
		LastModifiedTime:  unixMillisToTime(r.LastModifiedTime),
		Language:          r.Language,
		ImportedLibraries: r.ImportedLibraries,
		Body:              r.DefinitionBody,
	}
	args, err := bqToArgs(r.Arguments)
	if err != nil {
		return nil, err
	}
	meta.Arguments = args
	ret, err := bqToStandardSQLDataType(r.ReturnType)
	if err != nil {
		return nil, err
	}
	meta.ReturnType = ret
	return meta, nil
}
