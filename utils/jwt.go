package utils

import (
	"errors"
	"reflect"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a JWT token from any struct or map with expiry time
func GenerateJWT(user any, expiryMinutes int, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Unix(), // Expiry time
	}

	// Convert struct (or pointer to struct) to map
	userMap, err := toMap(user)
	if err != nil {
		return "", err
	}

	// Add user data to claims
	maps.Copy(claims, userMap)

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// toMap converts a struct (or a pointer to struct) to a map using reflection
func toMap(input any) (map[string]any, error) {
	v := reflect.ValueOf(input)

	// Handle pointers by dereferencing them
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle maps directly
	if v.Kind() == reflect.Map {
		mapData, ok := input.(map[string]any)
		if !ok {
			return nil, errors.New("invalid map format")
		}
		return mapData, nil
	}

	// Ensure it's a struct
	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or a map")
	}

	output := make(map[string]any)
	t := v.Type()

	for i := range v.NumField() {
		field := t.Field(i)

		// Ignore unexported fields (private fields)
		if !v.Field(i).CanInterface() {
			continue
		}

		output[field.Name] = v.Field(i).Interface()
	}

	return output, nil
}
