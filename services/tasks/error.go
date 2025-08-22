package tasks

import "errors"

var (
	crontabExpressionErr                      = errors.New("crontab expression is invalid")
	durationExpressionErr                     = errors.New("duration expression is invalid")
	durationRandomExpressionErr               = errors.New("durationRandom expression is invalid")
	oneTimeJobStartDateTimesExpressionErr     = errors.New("oneTimeJobStartDateTimes expression is invalid")
	oneTimeJobStartDateTimesExpressionExpired = errors.New("oneTimeJobStartDateTimes expression expired")
	notSupportedTaskTypeErr                   = errors.New("not supported task type")
	notSupportedTaskMethodErr                 = errors.New("not supported task method")

	createOrUpdateJobErr = errors.New("create or update job error")

	dbTaskUUIDInvalidErr = errors.New("task uuid is invalid")
)
