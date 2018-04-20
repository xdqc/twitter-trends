package tweet

//Tweet the tweet json structure generated from online tool https://mholt.github.io/json-to-go/
//(only care about relavent json objects/arrays for this project)
type Tweet struct {
	// PossiblySensitiveEditable bool        `json:"possibly_sensitive_editable"`
	// Place                     interface{} `json:"place"`
	Text string `json:"text"`
	// IDStr                     string      `json:"id_str"`
	// Favorited                 bool        `json:"favorited"`
	// PossiblySensitive         bool        `json:"possibly_sensitive"`
	// Coordinates               interface{} `json:"coordinates"`
	// Geo                       interface{} `json:"geo"`
	// CreatedAt                 string      `json:"created_at"`
	// InReplyToStatusIDStr      interface{} `json:"in_reply_to_status_id_str"`
	// Contributors              interface{} `json:"contributors"`
	// InReplyToScreenName       interface{} `json:"in_reply_to_screen_name"`
	// Source                    string      `json:"source"`
	// InReplyToUserIDStr        interface{} `json:"in_reply_to_user_id_str"`
	// InReplyToUserID           interface{} `json:"in_reply_to_user_id"`
	// InReplyToStatusID         interface{} `json:"in_reply_to_status_id"`
	// Retweeted                 bool        `json:"retweeted"`
	// RetweetCount              int         `json:"retweet_count"`
	// Truncated                 bool        `json:"truncated"`
	User struct {
		// IsTranslator                   bool        `json:"is_translator"`
		// ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
		// ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
		// IDStr                          string      `json:"id_str"`
		// ProfileTextColor               string      `json:"profile_text_color"`
		// StatusesCount                  int         `json:"statuses_count"`
		// Following                      interface{} `json:"following"`
		// ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
		// FollowersCount                 int         `json:"followers_count"`
		// ProfileImageURL                string      `json:"profile_image_url"`
		// DefaultProfileImage            bool        `json:"default_profile_image"`
		// CreatedAt                      string      `json:"created_at"`
		// ProfileLinkColor               string      `json:"profile_link_color"`
		// Description                    string      `json:"description"`
		TimeZone string `json:"time_zone"`
		// FavouritesCount                int         `json:"favourites_count"`
		// FriendsCount                   int         `json:"friends_count"`
		// URL                            string      `json:"url"`
		// Verified                       bool        `json:"verified"`
		// ProfileBackgroundColor         string      `json:"profile_background_color"`
		// ProfileBackgroundTile          bool        `json:"profile_background_tile"`
		// ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
		// ContributorsEnabled            bool        `json:"contributors_enabled"`
		// GeoEnabled                     bool        `json:"geo_enabled"`
		// Notifications                  interface{} `json:"notifications"`
		// ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
		// Protected                      bool        `json:"protected"`
		// Location                       string      `json:"location"`
		// ListedCount                    int         `json:"listed_count"`
		// FollowRequestSent              interface{} `json:"follow_request_sent"`
		// Name                           string      `json:"name"`
		// ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
		// ID                             int         `json:"id"`
		// DefaultProfile                 bool        `json:"default_profile"`
		// ShowAllInlineMedia             bool        `json:"show_all_inline_media"`
		// Lang string `json:"lang"`
		// UtcOffset  int    `json:"utc_offset"`
		// ScreenName string `json:"screen_name"`
	} `json:"user"`
	Lang string `json:"lang"`
	// ID       int64 `json:"id"`
	Entities struct {
		// Urls []struct {
		// 	DisplayURL  string `json:"display_url"`
		// 	Indices     []int  `json:"indices"`
		// 	URL         string `json:"url"`
		// 	ExpandedURL string `json:"expanded_url"`
		// } `json:"urls"`
		Hashtags []struct {
			Text string `json:"text"`
			// Indices []int  `json:"indices"`
		} `json:"hashtags"`
		// UserMentions []interface{} `json:"user_mentions"`
	} `json:"entities"`
}
