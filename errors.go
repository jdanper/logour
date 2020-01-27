package main

import "errors"

var (
	errInvalidJSON = errors.New("invalid json")
	errParseJSON = errors.New("unable to parse json")
	errEmptyMandatoryField = errors.New("event must contains all mandatory fields")
	errInvalidKind = errors.New("invalid event kind")
)
