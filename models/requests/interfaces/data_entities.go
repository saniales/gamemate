package requestInterfaces

import "github.com/labstack/echo"

//FormDecodable represents a struct convertible from a form submitted.
type FormDecodable interface {
	FromForm(echo.Context) error
}
