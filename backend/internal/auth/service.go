package auth

var globalJWTService *JWTService

// SetGlobalJWTService sets the global JWT service instance
func SetGlobalJWTService(svc *JWTService) {
	globalJWTService = svc
}

// GetGlobalJWTService returns the global JWT service instance
func GetGlobalJWTService() *JWTService {
	return globalJWTService
}
