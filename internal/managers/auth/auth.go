package auth

type AuthManager struct {
	jwtKey []byte
}

func NewAuthMan(jwtKey []byte) *AuthManager {
	return &AuthManager{
		jwtKey: jwtKey,
	}
}

func (a *AuthManager) JWTKey() []byte {
	return a.jwtKey
}
