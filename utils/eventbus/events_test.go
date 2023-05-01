/*
Copyright 2023 QKP Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eventbus

import (
	"fmt"
	"testing"
	"time"
)

func TestNewEventBus(t *testing.T) {
	q := NewEventBus()
	q.Add("name", time.Second*10)
	q.Add("This is a test event unit", time.Minute*3)

	for !q.IsEmpty() {
		item, err := q.Remove()
		if err != nil {
			fmt.Println("Error removing item from queue:", err)
			break
		}
		fmt.Println("Got item from queue:", item)
		time.Sleep(time.Second) // Simulate processing time
	}
}
