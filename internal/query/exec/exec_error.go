// Copyright 2020 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)
import "github.com/gogo/status"

// extractExecError extract an ExecError from an error returned from the
// grpc api. Return nil the error doesn't carry an ExecError.
func extractExecError(err error) *ExecError {

	if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {

		st := se.GRPCStatus()

		if len(st.Details()) == 0 {
			return nil
		}

		if execError, ok := st.Details()[0].(*ExecError); ok {
			return execError
		}

		return nil

	}

	return nil

}

func newExecError(detail string) *ExecError {
	id := uuid.New()
	return &ExecError{
		Id:     id[:],
		Detail: detail,
	}
}

func (m *ExecError) Err() error {

	st, err := status.New(codes.Canceled, m.Detail).WithDetails(m)

	if err != nil {
		panic(err)
	}

	return st.Err()

}
