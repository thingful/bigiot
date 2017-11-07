// Copyright 2017 Thingful Ltd.
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

// Package bigiot is an attempt at porting the BIGIot client library from Java
// to Go, adapting the library where appropriate to better fit Go idioms and
// practices. This is very much a work in progress, so currently is a long way
// from supporting the same range of functionality as the Java library.
//
// Planned functionality:
//   * register an offering in the marketplace
//   * unregister an offering from the marketplace
//   * validating tokens presented by offering subscribers
//   * discovering an offering in the marketplace
//   * subscribing to an offering
//
package bigiot
