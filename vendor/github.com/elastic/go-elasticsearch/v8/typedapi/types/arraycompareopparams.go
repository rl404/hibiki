// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/a0da620389f06553c0727f98f95e40dbb564fcca

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/quantifier"
)

// ArrayCompareOpParams type.
//
// https://github.com/elastic/elasticsearch-specification/blob/a0da620389f06553c0727f98f95e40dbb564fcca/specification/watcher/_types/Conditions.ts#L27-L30
type ArrayCompareOpParams struct {
	Quantifier quantifier.Quantifier `json:"quantifier"`
	Value      FieldValue            `json:"value"`
}

func (s *ArrayCompareOpParams) UnmarshalJSON(data []byte) error {

	dec := json.NewDecoder(bytes.NewReader(data))

	for {
		t, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		switch t {

		case "quantifier":
			if err := dec.Decode(&s.Quantifier); err != nil {
				return err
			}

		case "value":
			if err := dec.Decode(&s.Value); err != nil {
				return err
			}

		}
	}
	return nil
}

// NewArrayCompareOpParams returns a ArrayCompareOpParams.
func NewArrayCompareOpParams() *ArrayCompareOpParams {
	r := &ArrayCompareOpParams{}

	return r
}
