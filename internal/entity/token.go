package entity

type Token struct {
	IDToken string `json:"id_token"`
	User    User   `json:"user"`
}
