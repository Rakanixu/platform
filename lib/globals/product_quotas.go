package globals

import "time"

const (
	QUOTA_TIME_LIMITER        = time.Hour
	QUOTA_TIME_LIMITER_STRING = "hour"

	PRODUCT_TYPE_PERSONAL   = "personal"
	PRODUCT_TYPE_TEAM       = "team"
	PRODUCT_TYPE_ENTERPRISE = "enterprise"
)

var PRODUCT_QUOTAS = struct {
	M map[string]map[string]map[string]interface{}
}{
	M: map[string]map[string]map[string]interface{}{
		PRODUCT_TYPE_PERSONAL: map[string]map[string]interface{}{
			DATASOURCE_SERVICE_NAME: map[string]interface{}{
				"label":      "Datasource service",
				"icon":       "file:cloud-queue",
				"handler":    0,
				"subscriber": 0,
			},
			CRAWLER_SERVICE_NAME: map[string]interface{}{
				"label":      "Discovery service",
				"icon":       "action:explore",
				"handler":    0,
				"subscriber": 0,
			},
			NOTIFICATION_SERVICE_NAME: map[string]interface{}{
				"label":      "Notification service",
				"icon":       "action:announcement",
				"handler":    0,
				"subscriber": 0,
			},
			FILE_SERVICE_NAME: map[string]interface{}{
				"label":      "File service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			USER_SERVICE_NAME: map[string]interface{}{
				"label":      "User service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			CHANNEL_SERVICE_NAME: map[string]interface{}{
				"label":      "Channel service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			QUOTA_SERVICE_NAME: map[string]interface{}{
				"label":      "Quota service",
				"icon":       "action:lock-outline",
				"handler":    0,
				"subscriber": 0,
			},
			PROFILE_SERVICE_NAME: map[string]interface{}{
				"label":      "Profile service",
				"icon":       "image:photo",
				"handler":    0,
				"subscriber": 0,
			},
			THUMBNAIL_SERVICE_NAME: map[string]interface{}{
				"label":      "Thumbnail service",
				"icon":       "image:photo-size-select-actual",
				"handler":    0,
				"subscriber": 0,
			},
			AUDIO_SERVICE_NAME: map[string]interface{}{
				"label":      "Speech to text service",
				"icon":       "image:audiotrack",
				"handler":    0,
				"subscriber": 1,
			},
			DOCUMENT_SERVICE_NAME: map[string]interface{}{
				"label":      "Content extraction service",
				"icon":       "action:find-in-page",
				"handler":    0,
				"subscriber": 0,
			},
			IMAGE_SERVICE_NAME: map[string]interface{}{
				"label":      "Image content service",
				"icon":       "image:photo-library",
				"handler":    0,
				"subscriber": 10,
			},
			ENTITIES_SERVICE_NAME: map[string]interface{}{
				"label":      "Entity extraction service",
				"icon":       "action:description",
				"handler":    0,
				"subscriber": 10,
			},
			SENTIMENT_SERVICE_NAME: map[string]interface{}{
				"label":      "Sentiment extraction service",
				"icon":       "social:mood",
				"handler":    0,
				"subscriber": 10,
			},
			TRANSLATE_SERVICE_NAME: map[string]interface{}{
				"label": "Translation service",
				// TODO: change this
				"icon":       "social:mood",
				"handler":    20,
				"subscriber": 20,
			},
		},
		PRODUCT_TYPE_TEAM: map[string]map[string]interface{}{
			DATASOURCE_SERVICE_NAME: map[string]interface{}{
				"label":      "Datasource service",
				"icon":       "file:cloud-queue",
				"handler":    0,
				"subscriber": 0,
			},
			CRAWLER_SERVICE_NAME: map[string]interface{}{
				"label":      "Discovery service",
				"icon":       "action:explore",
				"handler":    0,
				"subscriber": 0,
			},
			NOTIFICATION_SERVICE_NAME: map[string]interface{}{
				"label":      "Notification service",
				"icon":       "action:announcement",
				"handler":    0,
				"subscriber": 0,
			},
			FILE_SERVICE_NAME: map[string]interface{}{
				"label":      "File service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			USER_SERVICE_NAME: map[string]interface{}{
				"label":      "User service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			CHANNEL_SERVICE_NAME: map[string]interface{}{
				"label":      "Channel service",
				"icon":       "editor:insert-drive-file",
				"handler":    0,
				"subscriber": 0,
			},
			QUOTA_SERVICE_NAME: map[string]interface{}{
				"label":      "Quota service",
				"icon":       "action:lock-outline",
				"handler":    0,
				"subscriber": 0,
			},
			PROFILE_SERVICE_NAME: map[string]interface{}{
				"label":      "Profile service",
				"icon":       "image:photo",
				"handler":    0,
				"subscriber": 0,
			},
			THUMBNAIL_SERVICE_NAME: map[string]interface{}{
				"label":      "Thumbnail service",
				"icon":       "image:photo-size-select-actual",
				"handler":    0,
				"subscriber": 0,
			},
			AUDIO_SERVICE_NAME: map[string]interface{}{
				"label":      "Speech to text service",
				"icon":       "image:audiotrack",
				"handler":    0,
				"subscriber": 5,
			},
			DOCUMENT_SERVICE_NAME: map[string]interface{}{
				"label":      "Content extraction service",
				"icon":       "action:find-in-page",
				"handler":    0,
				"subscriber": 0,
			},
			IMAGE_SERVICE_NAME: map[string]interface{}{
				"label":      "Image content service",
				"icon":       "image:photo-library",
				"handler":    0,
				"subscriber": 100,
			},
			ENTITIES_SERVICE_NAME: map[string]interface{}{
				"label":      "Entity extraction service",
				"icon":       "action:description",
				"handler":    0,
				"subscriber": 50,
			},
			SENTIMENT_SERVICE_NAME: map[string]interface{}{
				"label":      "Sentiment extraction service",
				"icon":       "social:mood",
				"handler":    0,
				"subscriber": 50,
			},
			TRANSLATE_SERVICE_NAME: map[string]interface{}{
				"label": "Translation service",
				// TODO: change this
				"icon":       "social:mood",
				"handler":    20,
				"subscriber": 20,
			},
		},
	},
}
