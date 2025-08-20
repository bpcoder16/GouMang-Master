package tasks

import "errors"

var (
	crontabExpressionErr               = errors.New("crontab expression is invalid")
	durationExpressionErr              = errors.New("duration expression is invalid")
	durationRandomExpressionErr        = errors.New("duration random expression is invalid")
	oneTimeJobStartDateTimesExpression = errors.New("one time job start date times expression is invalid")
)
