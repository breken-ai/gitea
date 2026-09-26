// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package org

import (
	"net/http"
	"time"

	"gitea.dev/models/organization"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/templates"
	shared_user "gitea.dev/routers/web/shared/user"
	"gitea.dev/services/context"
)

const tplByRepos templates.TplName = "org/worktime"

// parseOrgTimes contains functionality that is required in all these functions,
// like parsing the date from the request, setting default dates, etc.
func parseOrgTimes(ctx *context.Context) (unixFrom, unixTo int64) {
	rangeFrom := ctx.FormString("from")
	rangeTo := ctx.FormString("to")
	// the dates are days in the instance's UI time zone, the same as issue deadlines and the activity heatmap
	now := time.Now().In(setting.DefaultUILocation)
	if rangeFrom == "" {
		rangeFrom = now.Format("2006-01") + "-01" // defaults to start of current month
	}
	if rangeTo == "" {
		rangeTo = now.Format("2006-01-02") // defaults to today
	}

	ctx.Data["RangeFrom"] = rangeFrom
	ctx.Data["RangeTo"] = rangeTo

	timeFrom, err := time.ParseInLocation("2006-01-02", rangeFrom, setting.DefaultUILocation)
	if err != nil {
		ctx.ServerError("time.Parse", err)
	}
	timeTo, err := time.ParseInLocation("2006-01-02", rangeTo, setting.DefaultUILocation)
	if err != nil {
		ctx.ServerError("time.Parse", err)
	}
	unixFrom = timeFrom.Unix()
	unixTo = timeTo.Add(1440*time.Minute - 1*time.Second).Unix() // humans expect that we include the ending day too
	return unixFrom, unixTo
}

func Worktime(ctx *context.Context) {
	ctx.Data["PageIsOrgTimes"] = true

	unixFrom, unixTo := parseOrgTimes(ctx)
	if ctx.Written() {
		return
	}

	worktimeBy := ctx.FormString("by")
	ctx.Data["WorktimeBy"] = worktimeBy

	var worktimeSumResult any
	var err error
	switch worktimeBy {
	case "milestones":
		worktimeSumResult, err = organization.GetWorktimeByMilestones(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByMilestones"] = true
	case "members":
		worktimeSumResult, err = organization.GetWorktimeByMembers(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByMembers"] = true
	default: /* by repos */
		worktimeSumResult, err = organization.GetWorktimeByRepos(ctx, ctx.Org.Organization, unixFrom, unixTo)
		ctx.Data["WorktimeByRepos"] = true
	}
	if err != nil {
		ctx.ServerError("GetWorktime", err)
		return
	}

	if _, err := shared_user.RenderUserOrgHeader(ctx); err != nil {
		ctx.ServerError("RenderUserOrgHeader", err)
		return
	}

	ctx.Data["WorktimeSumResult"] = worktimeSumResult
	ctx.HTML(http.StatusOK, tplByRepos)
}
