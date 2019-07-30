// Copyright 2016 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package martianurl

// MatchHost matches two URL hosts with support for wildcards.
func MatchHost(host, match string) bool {
	// Short circuit if host is empty.
	if host == "" {
		return false
	}

	// Exact match, no need to loop.
	if host == match {
		return true
	}

	// Walk backward over the host.
	hi := len(host) - 1
	for mi := len(match) - 1; mi >= 0; mi-- {
		// Found wildcard, skip to next period.
		if match[mi] == '*' {
			for hi > 0 && host[hi] != '.' {
				hi--
			}

			// Wildcard was the leftmost part and we have walked the entire host,
			// success.
			if mi == 0 && hi == 0 {
				return true
			}

			continue
		}

		if host[hi] != match[mi] {
			return false
		}

		// We have walked the entire host, if we have not walked the entire matcher
		// (mi != 0) that means the matcher has remaining characters to match and
		// thus the host cannot match.
		if hi == 0 {
			return mi == 0
		}

		hi--
	}

	// We have walked the entire length of the matcher, but haven't finished
	// walking the host thus they cannot match.
	return false
}
