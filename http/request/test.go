package request

type (
	TestPost struct {
		TestParam string `json:"test_param" validate:"required"`
	}
)
