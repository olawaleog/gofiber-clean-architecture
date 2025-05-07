package model

// language=json
var _ = `
{"description":"3 Idowu Taylor Street, Lagos, Nigeria","matched_substrings":[{"length":1,"offset":0},{"length":19,"offset":2}],"place_id":"ChIJB5xBsy71OxARMU6yL5z2NOU","reference":"ChIJB5xBsy71OxARMU6yL5z2NOU","structured_formatting":{"main_text":"3 Idowu Taylor Street","main_text_matched_substrings":[{"length":1,"offset":0},{"length":19,"offset":2}],"secondary_text":"Lagos, Nigeria"},"terms":[{"offset":0,"value":"3"},{"offset":2,"value":"Idowu Taylor Street"},{"offset":23,"value":"Lagos"},{"offset":30,"value":"Nigeria"}],"types":["geocode","premise"]}

`

type AddressModel struct {
	Description       string `json:"description"`
	MatchedSubstrings []struct {
		Length int `json:"length"`
		Offset int `json:"offset"`
	} `json:"matched_substrings"`
	PlaceId              string `json:"place_id"`
	Reference            string `json:"reference"`
	StructuredFormatting struct {
		MainText                  string `json:"main_text"`
		MainTextMatchedSubstrings []struct {
			Length int `json:"length"`
			Offset int `json:"offset"`
		} `json:"main_text_matched_substrings"`
		SecondaryText string `json:"secondary_text"`
	} `json:"structured_formatting"`
	Terms []struct {
		Offset int    `json:"offset"`
		Value  string `json:"value"`
	} `json:"terms"`
	Types  []string `json:"types"`
	UserId uint
}

type AddressResponseModel struct {
	Id          uint   `json:"id"`
	Description string `json:"description"`
	PlaceId     string `json:"placeId"`
}
