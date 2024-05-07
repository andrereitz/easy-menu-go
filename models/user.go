package models

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	Id            int     `json:"id"`
	Email         string  `json:"email"`
	BusinessName  *string `json:"business_name"`
	BusinessUrl   *string `json:"business_url"`
	BusinessColor *string `json:"business_color"`
	BusinessLogo  *string `json:"business_logo"`
}

type NewUserData struct {
	Email string `json:"email"`
	Hash  string `json:"password"`
}

type LogoMeta struct {
	Url  string `json:"url"`
	User int    `json:"user"`
}
