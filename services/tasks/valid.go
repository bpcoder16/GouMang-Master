package tasks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/logit"
	robfigCron "github.com/robfig/cron/v3"
)

func IsValidCrontabExpression(ctx context.Context, expression string) (err error) {
	p := robfigCron.NewParser(robfigCron.SecondOptional | robfigCron.Minute | robfigCron.Hour | robfigCron.Dom | robfigCron.Month | robfigCron.Dow | robfigCron.Descriptor)
	withLocation := fmt.Sprintf("CRON_TZ=%s %s", env.TimeLocation().String(), expression)
	_, err = p.Parse(withLocation)
	if err != nil {
		logit.Context(ctx).WarnW("IsValidCrontabExpression.Err", "["+expression+"]"+err.Error())
		return crontabExpressionErr
	}
	return
}

func IsValidDurationExpression(ctx context.Context, expression string) (durationMillisecond time.Duration, err error) {
	var durationMillisecondInt int
	durationMillisecondInt, err = strconv.Atoi(expression)
	if err != nil {
		logit.Context(ctx).WarnW("IsValidDurationExpression.Err", "["+expression+"]"+err.Error())
		err = durationExpressionErr
		return
	}

	durationMillisecond = time.Duration(durationMillisecondInt) * time.Millisecond
	return
}

func IsValidDurationRandomExpression(expression string) (minDurationMillisecond, maxDurationMillisecond time.Duration, err error) {
	expressionList := strings.Split(expression, ",")
	if len(expressionList) != 2 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}
	var minDurationMillisecondInt, maxDurationMillisecondInt int
	minDurationMillisecondInt, err = strconv.Atoi(expressionList[0])
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
		return
	}
	maxDurationMillisecondInt, err = strconv.Atoi(expressionList[1])
	if err != nil {
		err = errors.New("[" + expression + "]" + err.Error())
		return
	}
	if minDurationMillisecondInt >= maxDurationMillisecondInt || minDurationMillisecondInt <= 0 || maxDurationMillisecondInt <= 0 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}
	minDurationMillisecond = time.Duration(minDurationMillisecondInt) * time.Millisecond
	maxDurationMillisecond = time.Duration(maxDurationMillisecondInt) * time.Millisecond
	return
}

func IsValidOneTimeJobStartDateTimesExpression(expression string) (timeList []time.Time, err error) {
	expressionList := strings.Split(expression, ",")
	if len(expressionList) == 0 {
		err = errors.New("[" + expression + "]invalid expression")
		return
	}

	timeList = make([]time.Time, 0, len(expressionList))
	for _, timeStr := range expressionList {
		startAt, errT := time.ParseInLocation(time.DateTime, timeStr, env.TimeLocation())
		if errT != nil {
			err = errors.New("[" + expression + "]" + errT.Error())
			return
		}
		timeList = append(timeList, startAt)
	}

	maxTime := timeList[0]
	for _, t := range timeList[1:] {
		if t.After(maxTime) {
			maxTime = t
		}
	}
	if maxTime.Before(time.Now()) {
		err = errors.New("[" + expression + "]max time is too early")
	}

	return
}
