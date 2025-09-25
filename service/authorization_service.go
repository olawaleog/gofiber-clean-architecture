package service

import ()

// AuthorizationService defines methods for authorization checks
type AuthorizationService interface {
	IsAdmin(claims map[string]interface{}) bool
	HasRole(claims map[string]interface{}, role string) bool
	CanAccessRefinery(claims map[string]interface{}, refineryId uint) bool
}

type authorizationServiceImpl struct{}

// NewAuthorizationService creates a new authorization service instance
func NewAuthorizationService() AuthorizationService {
	return &authorizationServiceImpl{}
}

// IsAdmin checks if the user has admin privileges
func (a *authorizationServiceImpl) IsAdmin(claims map[string]interface{}) bool {
	if claims == nil {
		return false
	}

	// Check the role in claims
	if claims["roles"] != nil {
		roles, ok := claims["roles"].([]interface{})
		if ok && len(roles) > 0 {
			for _, role := range roles {
				if roleMap, ok := role.(map[string]interface{}); ok {
					if roleVal, exists := roleMap["role"].(string); exists && roleVal == "admin" {
						return true
					}
				}
			}
		}
	}

	// Check user role field directly
	if userRole, exists := claims["role"].(string); exists && userRole == "admin" {
		return true
	}

	return false
}

// HasRole checks if the user has a specific role
func (a *authorizationServiceImpl) HasRole(claims map[string]interface{}, role string) bool {
	if claims == nil {
		return false
	}

	// Check the role in claims
	if claims["roles"] != nil {
		roles, ok := claims["roles"].([]interface{})
		if ok && len(roles) > 0 {
			for _, r := range roles {
				if roleMap, ok := r.(map[string]interface{}); ok {
					if roleVal, exists := roleMap["role"].(string); exists && roleVal == role {
						return true
					}
				}
			}
		}
	}

	// Check user role field directly
	if userRole, exists := claims["role"].(string); exists && userRole == role {
		return true
	}

	return false
}

// CanAccessRefinery checks if the user can access a specific refinery
func (a *authorizationServiceImpl) CanAccessRefinery(claims map[string]interface{}, refineryId uint) bool {
	if claims == nil {
		return false
	}

	// Admin can access all refineries
	if a.IsAdmin(claims) {
		return true
	}

	// Check if user is associated with this refinery
	if userRefineryId, ok := claims["refineryId"].(float64); ok {
		return uint(userRefineryId) == refineryId
	}

	return false
}
