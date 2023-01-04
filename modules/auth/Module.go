package auth

type AuthModule struct {
	Controller     *AuthController
	Service        IAuthService
	UserRepository IUserRepository
}
