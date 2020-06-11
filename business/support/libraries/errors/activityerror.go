package errors

import "net/http"

const ACTIVITY_SYSTEM_ERROR string = "ACTIVITY_SYSTEM_ERROR"
const ACTIVITY_SYSTEM_INVALID_RELATION string = "ACTIVITY_SYSTEM_INVALID_RELATION"
const ACTIVITY_SYSTEM_INVALID_REWARD string = "ACTIVITY_SYSTEM_INVALID_REWARD"
const ACTIVITY_SYSTEM_PLAN_FINISHED string = "ACTIVITY_SYSTEM_PLAN_FINISHED"
const ACTIVITY_SYSTEM_CREATE_RELATION_FAILED string = "ACTIVITY_SYSTEM_CREATE_RELATION_FAILED"
const ACTIVITY_SYSTEM_PAIRID_NO_EXSIT string = "ACTIVITY_SYSTEM_PAIRID_NO_EXSIT"
const ACTIVITY_SYSTEM_ACCOUNT_NO_ACTIVATED string = "ACTIVITY_SYSTEM_ACCOUNT_NO_ACTIVATED"
const ACTIVITY_SYSTEM_ACCOUNT_NO_EXIST string = "ACTIVITY_SYSTEM_ACCOUNT_NO_EXIST"

var (
	ActivitySysError                  = NewError(http.StatusOK, ACTIVITY_SYSTEM_ERROR, "活动系统出错了,请稍后!")
	ActivitySysInvalidRelationError   = NewError(http.StatusOK, ACTIVITY_SYSTEM_INVALID_RELATION, "关系错误!")
	ActivitySysInvalidRewardError     = NewError(http.StatusOK, ACTIVITY_SYSTEM_INVALID_REWARD, "奖励错误!")
	ActivitySysPlanFinishedError      = NewError(http.StatusOK, ACTIVITY_SYSTEM_PLAN_FINISHED, "奖励错误!")
	ActivitySysCreateRelationError    = NewError(http.StatusOK, ACTIVITY_SYSTEM_CREATE_RELATION_FAILED, "建立关系失败!")
	ActivitySysPairNoExsitError       = NewError(http.StatusOK, ACTIVITY_SYSTEM_PAIRID_NO_EXSIT, "举荐的用户不存在")
	ActivitySysPairAccountNoActivated = NewError(http.StatusOK, ACTIVITY_SYSTEM_ACCOUNT_NO_ACTIVATED, "用户没有激活")
	ActivitySysPairAccountNoExist     = NewError(http.StatusOK, ACTIVITY_SYSTEM_ACCOUNT_NO_EXIST, "用户不存在")
)

var (
	TeamMaxMemberLimit = NewError(http.StatusOK, "TEAM_MAX_MEMBER_LIMIT", "")
	TeamNotFound       = NewError(http.StatusOK, "TEAM_NOT_FOUND", "")
	TeamDisabled       = NewError(http.StatusOK, "TEAM_DISABLED", "")
	TeamUserNotFound   = NewError(http.StatusOK, "TEAM_USER_NOT_FOUND", "")
	TeamUserDisabled   = NewError(http.StatusOK, "TEAM_USER_DISABLED", "")
	OneOfTeamMember    = NewError(http.StatusOK, "ONE_OF_TEAM_MEMBER", "")

	TaskNotCompleted = NewError(http.StatusOK, "TASK_NOT_COMPLETED", "")
	TaskHasCompleted = NewError(http.StatusOK, "TASK_HAS_COMPLETED", "")
)
