package model

type UserModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phoneNumber"`
	EmailAddress string `json:"emailAddress"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
	IsActive     bool
	Id           uint   `json:"id"`
	Token        string `json:"token"`
	Address      string `json:"address"`
	Image        string `json:"image"`
	FileName     string `json:"fileName"`

	Street      string `json:"street"`
	Region      string `json:"region"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
	RefineryId  uint   `json:"refineryId"`
	CountryCode string `json:"countryCode"`
	OtpCode     string `json:"otpCode"`
}

type LoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordModel struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
