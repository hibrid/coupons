package coupon

type CouponConfig struct {
	Prefix        string `json:"prefix"`
	Suffix        string `json:"suffix"`
	Pattern       string `json:"pattern"`
	Divider       string `json:"divider"`
	CharacterSet  string `json:"characterSet"`
	Count         int    `json:"count"`
	MinimumLength int    `json:"minimumLength"`
	VanityName    string `json:"vanityName,omitempty"`
}
