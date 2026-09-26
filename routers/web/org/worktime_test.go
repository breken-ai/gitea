// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package org

import (
	"testing"
	"time"

	"gitea.dev/modules/setting"
	"gitea.dev/modules/test"
	"gitea.dev/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestParseOrgTimesUsesUILocation(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	defer test.MockVariableValue(&setting.DefaultUILocation, loc)()

	ctx, _ := contexttest.MockContext(t, "org3/-/worktime?from=2026-09-01&to=2026-09-30")
	unixFrom, unixTo := parseOrgTimes(ctx)
	assert.False(t, ctx.Written())
	assert.Equal(t, time.Date(2026, 9, 1, 0, 0, 0, 0, loc).Unix(), unixFrom)
	assert.Equal(t, time.Date(2026, 9, 30, 23, 59, 59, 0, loc).Unix(), unixTo)
}
