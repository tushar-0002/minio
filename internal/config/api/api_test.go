// Copyright (c) 2015-2026 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"strings"
	"testing"
	"time"

	"github.com/minio/minio/internal/config"
)

func presignKVS(value string) config.KVS {
	kvs := config.KVS{}
	kvs = append(kvs, DefaultKVS...)
	if value != "" {
		kvs.Set(apiPresignMaxExpiry, value)
	}
	return kvs
}

func TestLookupConfigPresignMaxExpiry(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  time.Duration
	}{
		{"default", "", 168 * time.Hour},
		{"custom", "1h", time.Hour},
		{"minimum", "1m", time.Minute},
		{"maximum", "168h", MaxPresignExpiry},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := LookupConfig(presignKVS(c.value))
			if err != nil {
				t.Fatalf("LookupConfig: %v", err)
			}
			if cfg.PresignMaxExpiry != c.want {
				t.Fatalf("PresignMaxExpiry = %v, want %v", cfg.PresignMaxExpiry, c.want)
			}
		})
	}
}

func TestLookupConfigPresignMaxExpiryInvalid(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		wantRange bool
	}{
		{"not a duration", "soon", false},
		{"below minimum", "59s", true},
		{"zero", "0s", true},
		{"negative", "-1h", true},
		{"above maximum", "200h", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := LookupConfig(presignKVS(c.value))
			if err == nil {
				t.Fatalf("expected an error for %q", c.value)
			}
			if c.wantRange && !strings.Contains(err.Error(), "between 1m0s and 168h0m0s") {
				t.Fatalf("error %q does not name the allowed range", err)
			}
		})
	}
}
