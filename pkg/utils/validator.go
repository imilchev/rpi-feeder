package utils

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var Validate = validator.New()
