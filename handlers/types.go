package handlers

type GenericReponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DataReponse struct {
	Data any `json:"data"`
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

type CategoryData struct {
	Id    int    `json:"id"`
	User  int    `json:"user"`
	Title string `json:"title"`
}
