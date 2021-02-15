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

package execbase

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"meerkat/internal/query/execpb"
	"runtime/debug"
)
import "github.com/gogo/status"

// ExtractExecError extract an ExecError from an error returned from the
// grpc api. Return nil the error doesn't carry an ExecError.
func ExtractExecError(err interface{}) *execpb.ExecError {

	if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {

		st := se.GRPCStatus()

		if len(st.Details()) == 0 {
			return nil
		}

		if execError, ok := st.Details()[0].(*execpb.ExecError); ok {
			return execError
		}

		return nil

	}

	return nil

}

func NewExecError(detail string, nodeName string) *execpb.ExecError {
	id := uuid.New()
	return &execpb.ExecError{
		Id:       id[:],
		Detail:   detail,
		NodeName: nodeName,
		Stack:    string(debug.Stack()), // TODO(gvelo): get a better stacktrace (ie. from pkg/errors )
	}
}

func BuildGRPCError(execError *execpb.ExecError) error {

	st, err := status.New(codes.Canceled, execError.Detail).WithDetails(execError)

	if err != nil {
		panic(err)
	}

	return st.Err()

}
