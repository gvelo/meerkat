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

package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	//s := ` StormEvents
	//      | limit 10
	//      | count`

	s := `Heartbeat
         | where TimeGenerated > start_time and TimeGenerated < end_time
         | summarize heartbeat_per_hour=count() by bin_at(TimeGenerated, 1h, start_time), Computer
         | extend available_per_hour=iff(heartbeat_per_hour>0, true, false)
         | summarize total_available_hours=countif(available_per_hour==true) by Computer 
         | extend total_number_of_buckets=round((end_time-start_time)/1h)+1
         | extend availability_rate=total_available_hours*100/total_number_of_buckets
`

	fmt.Println(s)

	p := NewParser()

	stmt, err := p.Parse(s)

	if err != nil {
		panic(err)
	}

	PrintAST(stmt)

}
