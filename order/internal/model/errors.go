package model

import "errors"

var SearchablePartsNotFound = errors.New("parts in order not founded")
var OrderNotFound = errors.New("order not found")
var OrderAlreadyPaid = errors.New("order was paid early")
var OrderCannotBeCanceled = errors.New("order was paid and cannot be canceled")
var InternalServerError = errors.New("internal server error")
