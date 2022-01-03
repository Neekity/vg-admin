package helper

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MyTime struct {
	time.Time
}

const TimeFormat = "2006-01-02 15:04:05"

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t *MyTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = MyTime{time.Time{}}
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = MyTime{now}
	return
}

func (t MyTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *MyTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = MyTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type Column map[string]interface{}

const PageSize = 20

func SearchPaginateData(page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			return db
		}
		if page < 0 {
			page = 1
		}

		offset := (page - 1) * PageSize
		return db.Offset(offset).Limit(PageSize)
	}
}

func QueryKey(val interface{}, queryKey string, flag string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if val != nil && val != "" {
			if flag == "like" {
				val = "%" + val.(string) + "%"
			}
			return db.Where(queryKey+" "+flag+" ?", val)
		} else {
			return db
		}

	}
}

func Paginate(curPage int, totalCount int, data interface{}) (*ApiListData, error) {
	if curPage <= 0 {
		curPage = 1
	}
	lastPage := int(totalCount/PageSize) + 1
	dataList, err := ParseModelList(data)
	if err != nil {
		return nil, err
	}
	return &ApiListData{
		TotalCount: totalCount,
		List:       dataList,
		PageSize:   PageSize,
		CurPage:    curPage,
		LastPage:   lastPage,
	}, nil
}

func ParseModelList(data interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &dataList)
	if err != nil {
		return nil, err
	}
	return dataList, err
}

func ParseModelDetail(data interface{}) (map[string]interface{}, error) {
	var dataList map[string]interface{}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &dataList)
	if err != nil {
		return nil, err
	}
	return dataList, err
}
