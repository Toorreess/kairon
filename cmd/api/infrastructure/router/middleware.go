package router

import (
	"encoding/json"
	"kairon/cmd/api/presenter"
	"kairon/config"
	"kairon/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/idtoken"
)

func RequestLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${status}] ${method} ${host}${uri} \n",
	})
}

func AuthLogger(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract and validate authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				projectIdSgmts := strings.Split(config.C.ProjectID, "-")
				envTag := projectIdSgmts[len(projectIdSgmts)-1]

				if os.Getenv("DEBUG") == "true" && (envTag == "dev" || envTag == "test" || envTag == "playground") {
					c.Set("user", "jtp")
					c.Set("role", "admin")
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, presenter.APIResponse(http.StatusUnauthorized, ""))
			}

			// Process Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, presenter.APIResponse(http.StatusUnauthorized, ""))
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Try Firebase auth first
			firebaseToken, err := authClient.VerifyIDToken(c.Request().Context(), token)
			if err == nil {
				role := firebaseToken.Claims["role"].(string)

				c.Set("user", firebaseToken.UID)
				c.Set("role", role)
				return next(c)
			}

			// Fallback to Google ID token validation
			if _, err := idtoken.Validate(c.Request().Context(), token, config.C.ProjectID); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, presenter.APIResponse(http.StatusUnauthorized, ""))
			}

			// Set integrator values
			c.Set("user", "integrator")
			c.Set("role", "superadmin")
			return next(c)
		}
	}
}

func CheckRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, presenter.APIResponse(http.StatusForbidden, ""))
			}

			for _, allowedRole := range roles {
				if allowedRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, presenter.APIResponse(http.StatusForbidden, ""))
		}
	}
}

func validatedChanges[T any](h func(c echo.Context, t T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		changes := make(map[string]any)
		err := json.NewDecoder(c.Request().Body).Decode(&changes)
		if err != nil {
			return echo.NewHTTPError(400, presenter.APIResponse(400, err.Error()))
		}

		var t T
		if err := utils.Map2Struct(changes, &t); err != nil {
			return echo.NewHTTPError(400, presenter.APIResponse(400, err.Error()))
		}

		val := reflect.ValueOf(t)
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
		allowedFields := make(map[string]bool)
		requiredFields := make(map[string]bool)
		e := val.Type()
		for i := 0; i < e.NumField(); i++ {
			fieldName := e.Field(i).Name
			fieldTags := e.Field(i).Tag
			if fieldTags.Get("json") != "-" && fieldTags.Get("json") != "" {
				fieldName = fieldTags.Get("json")
			} else {
				fieldName = strings.ToLower(fieldName)
			}
			if fieldTags.Get("updateAllowed") == "true" {
				allowedFields[fieldName] = true
			}

			if strings.Contains(fieldTags.Get("validate"), "required") {
				requiredFields[fieldName] = true
			}

			if strings.Contains(fieldTags.Get("validate"), "oneof") {
				oneofTag := fieldTags.Get("validate")
				oneofValues := strings.Split(strings.Split(oneofTag, "oneof=")[1], " ")
				fieldValue := val.Field(i).Interface()

				validValue := false
				if utils.IsEmpty(fieldValue) {
					validValue = true
				} else {
					for _, v := range oneofValues {
						if reflect.DeepEqual(fieldValue, v) {
							validValue = true
							break
						}
					}
				}

				if !validValue {
					return echo.NewHTTPError(http.StatusBadRequest,
						presenter.APIResponse(http.StatusBadRequest,
							fmt.Sprintf("Invalid value for field '%s'. Must be one of: %s", fieldName, strings.Join(oneofValues, ", "))),
					)
				}
			}

		}

		log.Println("updates: ", changes)
		log.Print("af: ", allowedFields)
		log.Print("rf: ", requiredFields)

		for k := range changes {
			if _, ok := allowedFields[k]; !ok {
				return echo.NewHTTPError(400, presenter.APIResponse(400, fmt.Sprintf("Invalid update field '%s', This field cannot be updated", k)))
			}

			if utils.IsEmpty(changes[k]) {
				if _, ok := requiredFields[k]; ok {
					return echo.NewHTTPError(400, presenter.APIResponse(400, fmt.Sprintf("Invalid update for field '%s', Cannot unset this field", k)))
				}
			}
		}

		c.Set("requestMap", changes)

		return h(c, t)
	}
}

func validated[T any](h func(c echo.Context, t T) error) echo.HandlerFunc {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	return func(c echo.Context) error {
		var t T
		if err := c.Bind(&t); err != nil {
			log.Println(err.Error())
			return echo.NewHTTPError(400, presenter.APIResponse(400, err.Error()))
		}

		if err := validate.Struct(t); err != nil {
			log.Println(err.Error())
			return echo.NewHTTPError(400, presenter.APIResponse(400, err.Error()))
		}

		return h(c, t)
	}
}
