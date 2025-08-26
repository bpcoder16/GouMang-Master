package controller

import (
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"
	"sort"

	"github.com/bpcoder16/Chestnut/v2/logit"
	"github.com/gin-gonic/gin"
)

type Task struct{}

// Config 获取任务配置相关信息
func (t *Task) Config(ctx *gin.Context) {
	// 获取标签列表
	tags, err := t.getTagList(ctx)
	if err != nil {
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}

	// 构建标签列表
	tagList := buildTagList(tags)

	returnSuccessJson(ctx, gin.H{
		"tagList": tagList,
	})
}



// getTagList 获取标签列表
func (t *Task) getTagList(ctx *gin.Context) (tagList []string, err error) {
	if err = global.DefaultDB.WithContext(ctx).Model(&db.GMTask{}).
		Where("status != ?", db.StatusDeleted).
		Distinct("tag").
		Pluck("tag", &tagList).Error; err != nil {
		logit.Context(ctx).ErrorW("getTagList.error", err)
		return
	}

	// 按字母顺序排序
	sort.Strings(tagList)
	return
}

// buildTagList 构建标签列表
func buildTagList(tags []string) []map[string]any {
	tagList := make([]map[string]any, 0, len(tags)+1)
	tagList = append(tagList, map[string]any{
		"name": "全部", "value": "",
	}, map[string]any{
		"name": "无标签", "value": "NoTag",
	})
	for _, tag := range tags {
		tagList = append(tagList, map[string]any{
			"name":  tag,
			"value": tag,
		})
	}
	return tagList
}
