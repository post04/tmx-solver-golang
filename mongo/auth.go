package mongo

import (
	"time"
)

func ValidateAPIKey(APIKey string) bool {
	for key, user := range users {
		if key == APIKey {
			if user.TotalUses >= user.MaxUses {
				return false
			}
			if user.ExpiresAt < time.Now().UnixMilli() {
				return false
			}
			return true
		}
	}
	return false
}

func ValidateAPIKeyWithoutUses(APIKey string) bool {
	us, err := GetUsers()
	if err != nil {
		return false
	}
	for key := range us {
		if key == APIKey {
			return true
		}
	}
	return false
}

func UpdateUsesCount(key string) {
	locker.Lock()
	user, ok := users[key]
	if !ok {
		locker.Unlock()
		return
	}
	user.TotalUses++
	users[key] = user
	locker.Unlock()
}

func UpdateUsers() {
	for {
		timeNow := time.Now().UnixMilli()
		locker.Lock()
		for key, user := range users {
			if user.ExpiresAt < timeNow || user.TotalUses >= user.MaxUses {
				users[key].Valid = false
				continue
			}
			UpdateUser(key, user.TotalUses)
		}
		locker.Unlock()
		time.Sleep(2 * time.Minute)
	}
}
