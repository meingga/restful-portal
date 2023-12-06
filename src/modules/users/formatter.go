package users

type UserRegisterFormatter struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func FormatRegisterUser(user User) UserRegisterFormatter {
	formatter := UserRegisterFormatter{
		ID:       user.ID,
		Username: user.Username,
	}

	return formatter
}

type UserLoginFormatter struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func FormatLoginUser(user User, AccessToken string, RefreshToken string) UserLoginFormatter {
	formatter := UserLoginFormatter{
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}

	return formatter
}

type UserRefreshTokenFormatter struct {
	AccessToken string `json:"access_token"`
}

func FormatRefreshTokenUser(AccessToken string) UserRefreshTokenFormatter {
	formatter := UserRefreshTokenFormatter{
		AccessToken: AccessToken,
	}

	return formatter
}

type UserFormatter struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func FormatUser(user User) UserFormatter {
	formatter := UserFormatter{
		ID:       user.ID,
		Username: user.Username,
	}

	return formatter
}
