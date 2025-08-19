package bootstrap

import (
	"encoding/json"
	"fmt"
	"goumang-master/global"
	"strconv"
	"strings"

	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/core/file/operations"
)

const (
	defaultVersion = 10000
)

func initDB() {
	fileVersionValue, err := operations.ReadOrCreate[int](env.RootPath()+"/data/version", defaultVersion)
	if err != nil {
		panic(err)
	}

	switch global.AppBizConfig.GormDBDriver {
	//case "mysql":
	//	initMySQL(fileVersionValue)
	default:
		initSQLite(fileVersionValue)
	}
}

//func initMySQL(fileVersionValue int) {
//	sqlFiles, err := operations.ListFilesSorted(env.RootPath() + "/migrate/mysql")
//	if err != nil {
//		panic(err)
//	}
//}

type sqlFile struct {
	CreateOrAlter  string
	InsertOrUpdate string
}

func initSQLite(fileVersionValue int) {
	sqlFilePaths, err := operations.ListFilesSorted(env.RootPath() + "/migrate/sqlite")
	if err != nil {
		panic(err)
	}

	pendingSQLFiles := make(map[int]sqlFile, 10)
	for _, sqlFilePath := range sqlFilePaths {
		tmpArr := strings.Split(sqlFilePath[:len(sqlFilePath)-4], "_")
		if len(tmpArr) == 2 {
			indexInt, errS := strconv.Atoi(tmpArr[0])
			if errS != nil {
				panic(errS)
			}
			if indexInt > fileVersionValue {
				var sqlFileValue sqlFile
				var isOK bool
				if sqlFileValue, isOK = pendingSQLFiles[indexInt]; !isOK {
					sqlFileValue = sqlFile{}
				}

				switch tmpArr[1] {
				case "createOrAlter":
					sqlFileValue.CreateOrAlter = sqlFilePath
				case "insertOrUpdate":
					sqlFileValue.InsertOrUpdate = sqlFilePath
				}
				pendingSQLFiles[indexInt] = sqlFileValue
			}
		}
	}

	fileVersionValue++
	for {
		if sqlFileValue, isOK := pendingSQLFiles[fileVersionValue]; isOK {
			var sqlValue string
			if len(sqlFileValue.CreateOrAlter) > 0 {
				createOrAlterValue, errR := operations.ReadFile(env.RootPath() + "/migrate/sqlite/" + sqlFileValue.CreateOrAlter)
				if errR != nil {
					panic(errR)
				}
				sqlValue += createOrAlterValue
			}
			if len(sqlFileValue.InsertOrUpdate) > 0 {
				insertOrUpdateValue, errR := operations.ReadFile(env.RootPath() + "/migrate/sqlite/" + sqlFileValue.InsertOrUpdate)
				if errR != nil {
					panic(errR)
				}
				sqlValue += insertOrUpdateValue
			}
			sqlValueList := strings.Split(sqlValue, ";")
			//errDB := global.DefaultDB.Transaction(func(tx *gorm.DB) error {
			for _, sql := range sqlValueList {
				fmt.Println(sql)
				//if errE := tx.Exec(sql).Error; err != nil {
				//	return errE
				//}
			}
			//	return nil
			//})
			//if errDB != nil {
			//	panic(errDB)
			//}
			_ = operations.WriteFile[int](env.RootPath()+"/data/version", fileVersionValue)
		} else {
			break
		}
		fileVersionValue = fileVersionValue + 1
	}

	sqlFilesJ, _ := json.Marshal(pendingSQLFiles)
	//fmt.Println(fileVersionValue)
	fmt.Println(string(sqlFilesJ))
}
