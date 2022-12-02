package common

import (
    "time"
	"errors"
    "strconv"
    "fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// handle date - if the date is given in ISO? format or as unix time

func QuerryParamParseTime(timeStr string) (primitive.DateTime, error) {
    timeInt, err := strconv.ParseInt(timeStr, 10, 64)
    if errors.Is(err, strconv.ErrSyntax) {
        t, err := time.Parse("2.1.2006-15:4:5", timeStr)
        if err != nil {
            t, err = time.Parse("2.1.2006", timeStr)
            if err != nil {
                return -1, err
            }
        }
        timeInt = t.Unix()
    } else if err != nil {
        fmt.Println("Some other error was found", err)
        return -1, err
    }

    return primitive.NewDateTimeFromTime(time.Unix(timeInt, 0)), nil
}
