// Lute HTTP - HTTP Server for Lute.
// Copyright (c) 2019-present, b3log.org
//
// Lute HTTP is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of the License at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"strings"
	"testing"
)

func TestNewLuteEnablesCallout(t *testing.T) {
	html := newLute().Markdown("", []byte("> [!NOTE] Title\n> Content"))
	output := string(html)
	if !strings.Contains(output, `<div class="callout" data-subtype="NOTE">`) {
		t.Fatalf("Callout was not rendered: %s", output)
	}
}
