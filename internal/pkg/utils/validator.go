package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var v *validator.Validate

func init() {
	v = validator.New()
	v.RegisterValidation("ISO8601Date", IsISO8601Date)
}

func Validate(c echo.Context, s interface{}) (err error) {
	if err = c.Bind(s); err != nil {
		err = fmt.Errorf("%s", "Something Went Wrong")
		return
	}

	if err = c.Validate(s); err != nil {
		c.Set("invalid-format", true)
		return
	}

	return
}

func ValidateStruct(ctx context.Context, s interface{}) (err error) {
	return v.StructCtx(ctx, s)
}

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^\\d{4}(-\\d\\d(-\\d\\d(T\\d\\d:\\d\\d(:\\d\\d)?(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?)?)?)?$"
	return regexp.MustCompile(ISO8601DateRegexString).MatchString(fl.Field().String())
}

func ValidateUUID(ctx context.Context, u string) (err error) {
	_, err = uuid.Parse(u)
	if err != nil {
		err = fmt.Errorf("invalid uuid")
		return
	}
	return
}
