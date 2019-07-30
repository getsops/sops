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
	"fmt"
	"time"

	"cloud.google.com/go/internal/optional"
	"cloud.google.com/go/internal/trace"
	bq "google.golang.org/api/bigquery/v2"
)

// Model represent a reference to a BigQuery ML model.
// Within the API, models are used largely for communicating
// statistical information about a given model, as creation of models is only
// supported via BigQuery queries (e.g. CREATE MODEL .. AS ..).
//
// For more info, see documentation for Bigquery ML,
// see: https://cloud.google.com/bigquery/docs/bigqueryml
type Model struct {
	ProjectID string
	DatasetID string
	// ModelID must contain only letters (a-z, A-Z), numbers (0-9), or underscores (_).
	// The maximum length is 1,024 characters.
	ModelID string

	c *Client
}

// FullyQualifiedName returns the ID of the model in projectID:datasetID.modelid format.
func (m *Model) FullyQualifiedName() string {
	return fmt.Sprintf("%s:%s.%s", m.ProjectID, m.DatasetID, m.ModelID)
}

// Metadata fetches the metadata for a model, which includes ML training statistics.
func (m *Model) Metadata(ctx context.Context) (mm *ModelMetadata, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Model.Metadata")
	defer func() { trace.EndSpan(ctx, err) }()

	req := m.c.bqs.Models.Get(m.ProjectID, m.DatasetID, m.ModelID).Context(ctx)
	setClientHeader(req.Header())
	var model *bq.Model
	err = runWithRetry(ctx, func() (err error) {
		model, err = req.Do()
		return err
	})
	if err != nil {
		return nil, err
	}
	return bqToModelMetadata(model)
}

// Update updates mutable fields in an ML model.
func (m *Model) Update(ctx context.Context, mm ModelMetadataToUpdate, etag string) (md *ModelMetadata, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Model.Update")
	defer func() { trace.EndSpan(ctx, err) }()

	bqm, err := mm.toBQ()
	if err != nil {
		return nil, err
	}
	call := m.c.bqs.Models.Patch(m.ProjectID, m.DatasetID, m.ModelID, bqm).Context(ctx)
	setClientHeader(call.Header())
	if etag != "" {
		call.Header().Set("If-Match", etag)
	}
	var res *bq.Model
	if err := runWithRetry(ctx, func() (err error) {
		res, err = call.Do()
		return err
	}); err != nil {
		return nil, err
	}
	return bqToModelMetadata(res)
}

// Delete deletes an ML model.
func (m *Model) Delete(ctx context.Context) (err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/bigquery.Model.Delete")
	defer func() { trace.EndSpan(ctx, err) }()

	req := m.c.bqs.Models.Delete(m.ProjectID, m.DatasetID, m.ModelID).Context(ctx)
	setClientHeader(req.Header())
	return req.Do()
}

// ModelMetadata represents information about a BigQuery ML model.
type ModelMetadata struct {
	// The user-friendly description of the model.
	Description string

	// The user-friendly name of the model.
	Name string

	// The type of the model.  Possible values include:
	// "LINEAR_REGRESSION" - a linear regression model
	// "LOGISTIC_REGRESSION" - a logistic regression model
	// "KMEANS" - a k-means clustering model
	Type string

	// The creation time of the model.
	CreationTime time.Time

	// The last modified time of the model.
	LastModifiedTime time.Time

	// The expiration time of the model.
	ExpirationTime time.Time

	// The geographic location where the model resides.  This value is
	// inherited from the encapsulating dataset.
	Location string

	// The input feature columns used to train the model.
	featureColumns []*bq.StandardSqlField

	// The label columns used to train the model.  Output
	// from the model will have a "predicted_" prefix for these columns.
	labelColumns []*bq.StandardSqlField

	// Information for all training runs, ordered by increasing start times.
	trainingRuns []*bq.TrainingRun

	Labels map[string]string

	// ETag is the ETag obtained when reading metadata. Pass it to Model.Update
	// to ensure that the metadata hasn't changed since it was read.
	ETag string
}

// TrainingRun represents information about a single training run for a BigQuery ML model.
// Experimental:  This information may be modified or removed in future versions of this package.
type TrainingRun bq.TrainingRun

// RawTrainingRuns exposes the underlying training run stats for a model using types from
// "google.golang.org/api/bigquery/v2", which are subject to change without warning.
// It is EXPERIMENTAL and subject to change or removal without notice.
func (mm *ModelMetadata) RawTrainingRuns() []*TrainingRun {
	if mm.trainingRuns == nil {
		return nil
	}
	var runs []*TrainingRun

	for _, v := range mm.trainingRuns {
		r := TrainingRun(*v)
		runs = append(runs, &r)
	}
	return runs
}

// RawLabelColumns exposes the underlying label columns used to train an ML model and uses types from
// "google.golang.org/api/bigquery/v2", which are subject to change without warning.
// It is EXPERIMENTAL and subject to change or removal without notice.
func (mm *ModelMetadata) RawLabelColumns() ([]*StandardSQLField, error) {
	return bqToModelCols(mm.labelColumns)
}

// RawFeatureColumns exposes the underlying feature columns used to train an ML model and uses types from
// "google.golang.org/api/bigquery/v2", which are subject to change without warning.
// It is EXPERIMENTAL and subject to change or removal without notice.
func (mm *ModelMetadata) RawFeatureColumns() ([]*StandardSQLField, error) {
	return bqToModelCols(mm.featureColumns)
}

func bqToModelCols(s []*bq.StandardSqlField) ([]*StandardSQLField, error) {
	if s == nil {
		return nil, nil
	}
	var cols []*StandardSQLField
	for _, v := range s {
		c, err := bqToStandardSQLField(v)
		if err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func bqToModelMetadata(m *bq.Model) (*ModelMetadata, error) {
	md := &ModelMetadata{
		Description:      m.Description,
		Name:             m.FriendlyName,
		Type:             m.ModelType,
		Location:         m.Location,
		Labels:           m.Labels,
		ExpirationTime:   unixMillisToTime(m.ExpirationTime),
		CreationTime:     unixMillisToTime(m.CreationTime),
		LastModifiedTime: unixMillisToTime(m.LastModifiedTime),
		featureColumns:   m.FeatureColumns,
		labelColumns:     m.LabelColumns,
		trainingRuns:     m.TrainingRuns,
		ETag:             m.Etag,
	}
	return md, nil
}

// ModelMetadataToUpdate is used when updating an ML model's metadata.
// Only non-nil fields will be updated.
type ModelMetadataToUpdate struct {
	// The user-friendly description of this model.
	Description optional.String

	// The user-friendly name of this model.
	Name optional.String

	// The time when this model expires.  To remove a model's expiration,
	// set ExpirationTime to NeverExpire.  The zero value is ignored.
	ExpirationTime time.Time

	labelUpdater
}

func (mm *ModelMetadataToUpdate) toBQ() (*bq.Model, error) {
	m := &bq.Model{}
	forceSend := func(field string) {
		m.ForceSendFields = append(m.ForceSendFields, field)
	}

	if mm.Description != nil {
		m.Description = optional.ToString(mm.Description)
		forceSend("Description")
	}

	if mm.Name != nil {
		m.FriendlyName = optional.ToString(mm.Name)
		forceSend("FriendlyName")
	}

	if !validExpiration(mm.ExpirationTime) {
		return nil, invalidTimeError(mm.ExpirationTime)
	}
	if mm.ExpirationTime == NeverExpire {
		m.NullFields = append(m.NullFields, "ExpirationTime")
	} else if !mm.ExpirationTime.IsZero() {
		m.ExpirationTime = mm.ExpirationTime.UnixNano() / 1e6
		forceSend("ExpirationTime")
	}
	labels, forces, nulls := mm.update()
	m.Labels = labels
	m.ForceSendFields = append(m.ForceSendFields, forces...)
	m.NullFields = append(m.NullFields, nulls...)
	return m, nil
}
