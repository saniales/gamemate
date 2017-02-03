package models

//Represents a way to convert to a struct from a submitted
//form (or from a request).
interface FormDecodable {
    //Converts from a submitted form (or request) to his struct.
    func FromForm(c echo.Context) (FormDecodable, error)
}
