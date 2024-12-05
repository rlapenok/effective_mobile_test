package swagger_client

type SwaggerClientError struct {
	Code int
	Err  error
}

func (e SwaggerClientError) Error() string {
	return e.Err.Error()
}
