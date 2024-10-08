/*
Copyright © 2024 Omni Aura peyton@omniaura.co

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
package assert

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func NoError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		slog.Error(fmt.Sprintf(
			"%s:%d: unexpected error: %v",
			file, line, err,
		))
		os.Exit(1)
	}
}
