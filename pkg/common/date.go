package common

import (
    "time"
    "errors"
    "strings"
)

var (
    ErrCantParseTime = errors.New("date error")
)

func IsDateLessThenDay(d1, d2 time.Time) bool {
    if d1.Year() < d2.Year() && d1.Month() < d2.Month() && d1.Day() < d2.Day() {
           return true
    }
    return false
}

func IsDateGreaterThenDay(d1, d2 time.Time) bool {
    if d1.Year() > d2.Year() && d1.Month() > d2.Month() && d1.Day() > d2.Day() {
           return true
    }
    return false
}

func IsDateSameDay(d1, d2 time.Time) bool {
    if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
           return true
    }
    return false
}

func ParseQueryParamDate(date string) (*time.Time, error) {
    parsedDate, err := time.Parse("2.1.2006", date)
    if err != nil {
        return nil, ErrCantParseTime
    }

    return &parsedDate, nil
}

func ParseBotFilenameDate(botName string) (*time.Time, error) {
    splitBotName := strings.Split(botName, "_")[1]
    splitBotName = splitBotName[:len(splitBotName)-4]
    botDate, err := time.Parse("2006-01-02T15-04-05", splitBotName)
    if err != nil {
        return nil, ErrCantParseTime
    }

    return &botDate, nil
}
