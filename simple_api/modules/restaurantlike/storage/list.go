package restaurantlikestorage

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

const timeLayout = "2006-01-02T15:04:05.999999"

func (s sqlStore) GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)

	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id;"`
		LikeCount    int `gorm:"column:count;"`
	}

	var listLikes []sqlData

	if err := s.db.Table(restaurantlikemodel.Like{}.TableName()).
		Select("restaurant_id, count(restaurant_id) as count").
		Where("restaurant_id in (?)", ids).
		Group("restaurant_id").Find(&listLikes).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLikes {
		result[item.RestaurantId] = item.LikeCount
	}

	return result, nil
}

func (s sqlStore) GetUsersLikeRestaurant(
	ctx context.Context,
	conditions map[string]interface{},
	filter *restaurantlikemodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]common.SimpleUser, error) {
	var result []restaurantlikemodel.Like
	db := s.db

	db = db.Table(restaurantlikemodel.Like{}.TableName()).Where(conditions)

	if f := filter; f != nil {
		if f.RestaurantId > 0 {
			db = db.Where("restaurant_id = ?", f.RestaurantId)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// for i := range moreKeys {
	// 	db = db.Preload(moreKeys[i])
	// }

	db = db.Preload("User")

	if cursor := paging.FakeCursor; cursor != "" {
		// NOTE: update fake cursor base on created_at
		timeCursor, err := time.Parse(timeLayout, string(base58.Decode(cursor)))

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at < ?", timeCursor.Format("2006-01-02T15:04:05")) // NOTE: format with SQL datetime format when query datetime
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Limit(paging.Limit).
		Order("created_at desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(result))

	for i, item := range result {

		users[i] = *item.User // NOTE: chổ này sẽ lỗi nếu User không tồn tại vì Preload sử dụng 2 câu query thay vì Join
		users[i].CreatedAt = item.CreatedAt
		users[i].UpdatedAt = nil
		if i == len(result)-1 {
			cursorStr := base58.Encode([]byte(item.CreatedAt.Format(timeLayout)))
			paging.NextCursor = cursorStr
		}

	}

	return users, nil
}
