/*
Copyright The Helm Authors.

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

// Package client provides the machinery to interact with Helm repositories
// and OCI registries hosting Helm charts alike by means of a common
// interface, the `HelmClient`. A common `Builder` interface simplifies
// creating a Helm client fit for purpose in a specific environment.
//
// In its current state this is an experimental package to explore the
// possibilities for interacting with HTTP(S) repositories and OCI registries
// alike in a Helm context. The idea is so to build the interfaces from an
// SDK-first perspective, treating the Helm CLI itself as one of possibly
// many downstream consumers wishing to interact with storage APIs hosting
// Helm charts.
package client
