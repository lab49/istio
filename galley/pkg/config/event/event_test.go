// Copyright 2019 Istio Authors
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

package event

import (
	"strings"
	"testing"

	"istio.io/istio/galley/pkg/config/collection"
	"istio.io/istio/galley/pkg/config/resource"

	. "github.com/onsi/gomega"

	"github.com/gogo/protobuf/types"
)

func TestEvent_String(t *testing.T) {
	e := resource.Entry{
		Metadata: resource.Metadata{
			Name:    resource.NewName("ns1", "rs1"),
			Version: "v1",
		},
		Item: &types.Empty{},
	}

	tests := []struct {
		i   Event
		exp string
	}{
		{
			i:   Event{},
			exp: "[Event](None)",
		},
		{
			i:   Event{Kind: Added, Entry: &e},
			exp: "[Event](Added: /ns1/rs1)",
		},
		{
			i:   Event{Kind: Updated, Entry: &e},
			exp: "[Event](Updated: /ns1/rs1)",
		},
		{
			i:   Event{Kind: Deleted, Entry: &e},
			exp: "[Event](Deleted: /ns1/rs1)",
		},
		{
			i:   Event{Kind: FullSync, Source: collection.NewName("foo")},
			exp: "[Event](FullSync: foo)",
		},
		{
			i:   Event{Kind: Kind(99), Source: collection.NewName("foo")},
			exp: "[Event](<<Unknown Kind 99>>)",
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			g := NewGomegaWithT(t)
			actual := tc.i.String()
			g.Expect(strings.TrimSpace(actual)).To(Equal(strings.TrimSpace(tc.exp)))
		})
	}
}

func TestEvent_DetailedString(t *testing.T) {
	e := resource.Entry{
		Metadata: resource.Metadata{
			Name:    resource.NewName("ns1", "rs1"),
			Version: "v1",
		},
		Item: &types.Empty{},
	}

	tests := []struct {
		i      Event
		prefix string
	}{
		{
			i:      Event{},
			prefix: "[Event](None",
		},
		{
			i:      Event{Kind: Added, Entry: &e},
			prefix: "[Event](Added: /ns1/rs1",
		},
		{
			i:      Event{Kind: Updated, Entry: &e},
			prefix: "[Event](Updated: /ns1/rs1",
		},
		{
			i:      Event{Kind: Deleted, Entry: &e},
			prefix: "[Event](Deleted: /ns1/rs1",
		},
		{
			i:      Event{Kind: FullSync, Source: collection.NewName("foo")},
			prefix: "[Event](FullSync: foo",
		},
		{
			i:      Event{Kind: Kind(99), Source: collection.NewName("foo")},
			prefix: "[Event](<<Unknown Kind 99>>",
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			g := NewGomegaWithT(t)
			actual := tc.i.String()
			actual = strings.TrimSpace(actual)
			expected := strings.TrimSpace(tc.prefix)
			g.Expect(actual).To(HavePrefix(expected))
		})
	}
}

func TestEvent_Clone(t *testing.T) {
	g := NewGomegaWithT(t)

	r := resource.Entry{
		Metadata: resource.Metadata{
			Name: resource.NewName("ns1", "rs1"),
			Labels: map[string]string{
				"foo": "bar",
			},
			Version: "v1",
		},
		Item: &types.Empty{},
	}

	e := Event{Kind: Added, Source: collection.NewName("boo"), Entry: &r}

	g.Expect(e.Clone()).To(Equal(e))
}

func TestEvent_FullSyncFor(t *testing.T) {
	g := NewGomegaWithT(t)

	e := FullSyncFor(collection.NewName("boo"))

	expected := Event{
		Kind:   FullSync,
		Source: collection.NewName("boo"),
	}
	g.Expect(e).To(Equal(expected))
}