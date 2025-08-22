package controller

import (
	"goumang-master/db"
	"goumang-master/errorcode"
	"goumang-master/global"

	"github.com/gin-gonic/gin"
)

type Task struct{}

func (t *Task) Config(ctx *gin.Context) {
	var dbNodeList []db.GMNode
	if err := global.DefaultDB.WithContext(ctx).Where(&db.GMNode{
		Status: db.StatusEnabled,
	}).Find(&dbNodeList).Error; err != nil {
		returnErrJson(ctx, errorcode.ErrServiceException)
		return
	}
	returnSuccessJson(ctx, gin.H{
		"taskType": []map[string]any{
			{"name": db.TaskTypeMap[db.NotChoose], "value": db.NotChoose},
			{"name": db.TaskTypeMap[db.TypeCron], "value": db.TypeCron},
			{"name": db.TaskTypeMap[db.TypeDuration], "value": db.TypeDuration},
			{"name": db.TaskTypeMap[db.TypeDurationRandom], "value": db.TypeDurationRandom},
			{"name": db.TaskTypeMap[db.TypeOneTimeJobStartDateTimes], "value": db.TypeOneTimeJobStartDateTimes},
		},
		"taskStatus": []map[string]any{
			{"name": db.TaskStatusMap[db.NotChoose], "value": db.NotChoose},
			{"name": db.TaskStatusMap[db.StatusPending], "value": db.StatusPending},
			{"name": db.TaskStatusMap[db.StatusEnabled], "value": db.StatusEnabled},
			{"name": db.TaskStatusMap[db.StatusConfigError], "value": db.StatusConfigError},
			{"name": db.TaskStatusMap[db.StatusConfigExpired], "value": db.StatusConfigExpired},
			{"name": db.TaskStatusMap[db.StatusWorkerServiceError], "value": db.StatusWorkerServiceError},
		},
		"nodeList": func() []map[string]any {
			nodeList := make([]map[string]any, 0, len(dbNodeList))
			for _, node := range dbNodeList {
				nodeList = append(nodeList, map[string]any{
					"name":  node.Title,
					"value": node.ID,
				})
			}
			return nodeList
		}(),
	})
}

func (t *Task) List(ctx *gin.Context) {
	var req struct {
		Title string `json:"title"`
		Tag   string `json:"tag"`

		Query                string   `json:"query"`
		Ids                  []string `json:"ids"`
		SaleType             int8     `json:"sale_type" binding:"oneof=0 1 2 3"`
		Page                 int      `json:"page" binding:"required,min=1"`
		PageSize             int      `json:"page_size" binding:"required,min=1"`
		IsRating             int8     `json:"is_rating" binding:"oneof=-1 0 1"`
		RatingCompany        int8     `json:"rating_company"`
		IsLimited            int8     `json:"is_limited" binding:"oneof=-1 0 1"`
		IsSign               int8     `json:"is_sign" binding:"oneof=-1 0 1"`
		GTEndTime            string   `json:"gt_end_time"`
		Status               int8     `json:"status"`
		SportTypeList        []string `json:"sport_type_list"`
		ExcludeSportTypeList []string `json:"exclude_sport_type_list"`
		ChannelList          []int8   `json:"channel_list"`
		PlayerNameList       []string `json:"player_name_list"`
		BrandList            []string `json:"brand_list"`
		IssueYearList        []string `json:"issue_year_list"`
		MinPrice             float64  `json:"min_price"`
		MaxPrice             float64  `json:"max_price"`

		SortField string `json:"sort_field"` // currentPrice 价格 popularity 人气热度 listingTime 上架时间 endTime 结束时间
		SortOrder string `json:"sort_order"` // asc 正序 desc 倒序
	}

	if err := paramsValidator(ctx, &req); err != nil {
		return
	}
}
