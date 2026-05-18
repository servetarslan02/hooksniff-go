package models

type BorderRadiusConfig struct {
	TopLeft     *string `json:"topLeft,omitempty"`
	TopRight    *string `json:"topRight,omitempty"`
	BottomLeft  *string `json:"bottomLeft,omitempty"`
	BottomRight *string `json:"bottomRight,omitempty"`
}

type FontSizeConfig struct {
	Base   *string `json:"base,omitempty"`
	Small  *string `json:"small,omitempty"`
	Large  *string `json:"large,omitempty"`
}
