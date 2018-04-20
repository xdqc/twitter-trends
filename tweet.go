package tweet

//Tweet the tweet json structure generated from online tool https://mholt.github.io/json-to-go/ (only parse relavent entries for this project)
type Tweet struct {
	// Favorited         bool        `json:"favorited"`
	// Contributors      interface{} `json:"contributors"`
	// Truncated         bool        `json:"truncated"`
	Text string `json:"text"`
	// IsQuoteStatus     bool        `json:"is_quote_status"`
	// InReplyToStatusID interface{} `json:"in_reply_to_status_id"`
	User struct {
		// FollowRequestSent              interface{} `json:"follow_request_sent"`
		// ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
		// DefaultProfileImage            bool        `json:"default_profile_image"`
		// ID                             int64       `json:"id"`
		// Verified                       bool        `json:"verified"`
		// ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
		// ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
		// ProfileTextColor               string      `json:"profile_text_color"`
		// FollowersCount                 int         `json:"followers_count"`
		// ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
		// IDStr                          string      `json:"id_str"`
		// ProfileBackgroundColor         string      `json:"profile_background_color"`
		// ListedCount                    int         `json:"listed_count"`
		// ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
		// UtcOffset                      interface{} `json:"utc_offset"`
		// StatusesCount                  int         `json:"statuses_count"`
		// Description                    string      `json:"description"`
		// FriendsCount                   int         `json:"friends_count"`
		// Location                       string      `json:"location"`
		// ProfileLinkColor               string      `json:"profile_link_color"`
		// ProfileImageURL                string      `json:"profile_image_url"`
		// Following                      interface{} `json:"following"`
		// GeoEnabled                     bool        `json:"geo_enabled"`
		// ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
		// Name                           string      `json:"name"`
		// Lang                           string      `json:"lang"`
		// ProfileBackgroundTile          bool        `json:"profile_background_tile"`
		// FavouritesCount                int         `json:"favourites_count"`
		// ScreenName                     string      `json:"screen_name"`
		// Notifications                  interface{} `json:"notifications"`
		// URL                            interface{} `json:"url"`
		// CreatedAt                      string      `json:"created_at"`
		// ContributorsEnabled            bool        `json:"contributors_enabled"`
		TimeZone interface{} `json:"time_zone"`
		// Protected                      bool        `json:"protected"`
		// DefaultProfile                 bool        `json:"default_profile"`
		// IsTranslator                   bool        `json:"is_translator"`
	} `json:"user"`
	// FilterLevel     string      `json:"filter_level"`
	// Geo             interface{} `json:"geo"`
	// ID              int64       `json:"id"`
	// FavoriteCount   int         `json:"favorite_count"`
	Lang string `json:"lang"`
	// RetweetedStatus struct {
	// 	Contributors      interface{} `json:"contributors"`
	// 	Truncated         bool        `json:"truncated"`
	// 	Text              string      `json:"text"`
	// 	IsQuoteStatus     bool        `json:"is_quote_status"`
	// 	InReplyToStatusID interface{} `json:"in_reply_to_status_id"`
	// 	ID                int64       `json:"id"`
	// 	FavoriteCount     int         `json:"favorite_count"`
	// 	Source            string      `json:"source"`
	// 	Retweeted         bool        `json:"retweeted"`
	// 	Coordinates       interface{} `json:"coordinates"`
	// 	Entities          struct {
	// 		UserMentions []interface{} `json:"user_mentions"`
	// 		Symbols      []interface{} `json:"symbols"`
	// 		Hashtags     []interface{} `json:"hashtags"`
	// 		Urls         []interface{} `json:"urls"`
	// 	} `json:"entities"`
	// 	InReplyToScreenName interface{} `json:"in_reply_to_screen_name"`
	// 	IDStr               string      `json:"id_str"`
	// 	RetweetCount        int         `json:"retweet_count"`
	// 	InReplyToUserID     interface{} `json:"in_reply_to_user_id"`
	// 	Favorited           bool        `json:"favorited"`
	// 	User                struct {
	// 		FollowRequestSent              interface{} `json:"follow_request_sent"`
	// 		ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	// 		DefaultProfileImage            bool        `json:"default_profile_image"`
	// 		ID                             int         `json:"id"`
	// 		Verified                       bool        `json:"verified"`
	// 		ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
	// 		ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
	// 		ProfileTextColor               string      `json:"profile_text_color"`
	// 		FollowersCount                 int         `json:"followers_count"`
	// 		ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
	// 		IDStr                          string      `json:"id_str"`
	// 		ProfileBackgroundColor         string      `json:"profile_background_color"`
	// 		ListedCount                    int         `json:"listed_count"`
	// 		ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
	// 		UtcOffset                      interface{} `json:"utc_offset"`
	// 		StatusesCount                  int         `json:"statuses_count"`
	// 		Description                    string      `json:"description"`
	// 		FriendsCount                   int         `json:"friends_count"`
	// 		Location                       interface{} `json:"location"`
	// 		ProfileLinkColor               string      `json:"profile_link_color"`
	// 		ProfileImageURL                string      `json:"profile_image_url"`
	// 		Following                      interface{} `json:"following"`
	// 		GeoEnabled                     bool        `json:"geo_enabled"`
	// 		ProfileBannerURL               string      `json:"profile_banner_url"`
	// 		ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
	// 		Name                           string      `json:"name"`
	// 		Lang                           string      `json:"lang"`
	// 		ProfileBackgroundTile          bool        `json:"profile_background_tile"`
	// 		FavouritesCount                int         `json:"favourites_count"`
	// 		ScreenName                     string      `json:"screen_name"`
	// 		Notifications                  interface{} `json:"notifications"`
	// 		URL                            string      `json:"url"`
	// 		CreatedAt                      string      `json:"created_at"`
	// 		ContributorsEnabled            bool        `json:"contributors_enabled"`
	// 		TimeZone                       interface{} `json:"time_zone"`
	// 		Protected                      bool        `json:"protected"`
	// 		DefaultProfile                 bool        `json:"default_profile"`
	// 		IsTranslator                   bool        `json:"is_translator"`
	// 	} `json:"user"`
	// 	Geo                  interface{} `json:"geo"`
	// 	InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	// 	Lang                 string      `json:"lang"`
	// 	CreatedAt            string      `json:"created_at"`
	// 	FilterLevel          string      `json:"filter_level"`
	// 	InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	// 	Place                interface{} `json:"place"`
	// } `json:"retweeted_status"`
	Entities struct {
		// UserMentions []struct {
		// 	ID         int    `json:"id"`
		// 	Indices    []int  `json:"indices"`
		// 	IDStr      string `json:"id_str"`
		// 	ScreenName string `json:"screen_name"`
		// 	Name       string `json:"name"`
		// } `json:"user_mentions"`
		// Symbols  []interface{} `json:"symbols"`
		Hashtags []interface{} `json:"hashtags"`
		// Urls     []interface{} `json:"urls"`
	} `json:"entities"`
	// InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	// Retweeted            bool        `json:"retweeted"`
	// Coordinates          interface{} `json:"coordinates"`
	// TimestampMs          string      `json:"timestamp_ms"`
	// Source               string      `json:"source"`
	// InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	// InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
	// IDStr                string      `json:"id_str"`
	// Place                interface{} `json:"place"`
	// RetweetCount         int         `json:"retweet_count"`
	// CreatedAt            string      `json:"created_at"`
	// InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
}
