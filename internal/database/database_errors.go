package database

import (
    "errors"
)

var (
    ErrFailure        = errors.New("Database error")
    ErrCantDecode     = errors.New("Database can't decode result into given type")
    ErrCursorFailure  = errors.New("Database cursor error")
    ErrDBNotSupported = errors.New("Database is not supported")
    ErrCantUnmarshal  = errors.New("Can't unmarshal")
    ErrCantConvert    = errors.New("Can't convert to right type")
)
