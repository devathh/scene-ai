package dtos

type Token struct {
	Access     string `json:"access"`
	AccessTTL  int64  `json:"access-ttl"`
	Refresh    string `json:"refresh"`
	RefreshTTL int64  `json:"refresh-ttl"`
}
