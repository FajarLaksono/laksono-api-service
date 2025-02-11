// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package utils

import "time"

type TimeProvider interface {
	Now() time.Time
}

func ItemExists(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
