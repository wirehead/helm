/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package main

import (
	"bytes"
	"regexp"
	"testing"

	"k8s.io/helm/pkg/proto/hapi/release"
)

func TestListCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		resp     []*release.Release
		expected string
		err      bool
	}{
		{
			name: "with a release",
			resp: []*release.Release{
				releaseMock(&releaseOptions{name: "thomas-guide"}),
			},
			expected: "thomas-guide",
		},
		{
			name: "list",
			args: []string{},
			resp: []*release.Release{
				releaseMock(&releaseOptions{name: "atlas"}),
			},
			expected: "NAME \tREVISION\tUPDATED                 \tSTATUS  \tCHART           \natlas\t1       \t(.*)\tDEPLOYED\tfoo-0.1.0-beta.1\n",
		},
		{
			name: "with a release, multiple flags",
			args: []string{"--deleted", "--deployed", "--failed", "-q"},
			resp: []*release.Release{
				releaseMock(&releaseOptions{name: "thomas-guide", statusCode: release.Status_DELETED}),
				releaseMock(&releaseOptions{name: "atlas-guide", statusCode: release.Status_DEPLOYED}),
			},
			// Note: We're really only testing that the flags parsed correctly. Which results are returned
			// depends on the backend. And until pkg/helm is done, we can't mock this.
			expected: "thomas-guide\natlas-guide",
		},
		{
			name: "with a release, multiple flags",
			args: []string{"--all", "-q"},
			resp: []*release.Release{
				releaseMock(&releaseOptions{name: "thomas-guide", statusCode: release.Status_DELETED}),
				releaseMock(&releaseOptions{name: "atlas-guide", statusCode: release.Status_DEPLOYED}),
			},
			// See note on previous test.
			expected: "thomas-guide\natlas-guide",
		},
	}

	var buf bytes.Buffer
	for _, tt := range tests {
		c := &fakeReleaseClient{
			rels: tt.resp,
		}
		cmd := newListCmd(c, &buf)
		cmd.ParseFlags(tt.args)
		err := cmd.RunE(cmd, tt.args)
		if (err != nil) != tt.err {
			t.Errorf("%q. expected error: %v, got %v", tt.name, tt.err, err)
		}
		re := regexp.MustCompile(tt.expected)
		if !re.Match(buf.Bytes()) {
			t.Errorf("%q. expected %q, got %q", tt.name, tt.expected, buf.String())
		}
		buf.Reset()
	}
}
