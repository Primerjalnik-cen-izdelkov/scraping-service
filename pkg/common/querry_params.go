package common

import (
    "time"
    "reflect"
    "errors"

    "github.com/ahmetcanozcan/fet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

    "github.com/pasztorpisti/qs"
)

var (
    ErrObjectIDConversion = errors.New("Can't convert string to the primitive.ObjectID")
)

// Query documents by given id(s). If id(s) is not provided this function does
// nothing. Otherwise we try to convert id(s) to mongoDB id object and add a
// query to match documents with specified id(s).
func QueryID(u []fet.Updater, id []string, fieldName string) ([]fet.Updater, error) {
    if id != nil {
        ids := []fet.Updater{}
        for _, i := range id {
            pid, err := primitive.ObjectIDFromHex(i)
            if err != nil {
                return nil, ErrObjectIDConversion
            }
            ids = append(ids, fet.Field(fieldName).Eq(pid, fet.IfNotNil))
        }
        u = append(u, fet.Or(ids...))
    }
    return u, nil
}

// Query document by given string slice. It matches document field 'fieldName'
// with strings from the slice.
func QueryStringSlice(u []fet.Updater, slice []string, fieldName string) []fet.Updater {
    if slice != nil {
        elements := []fet.Updater{}
        for _, s := range slice {
            elements = append(elements, fet.Field(fieldName).Eq(s, fet.IfNotNil))
        }
        u = append(u, fet.Or(elements...))
    }
    return u
}

// Query document by given int slice. It matches document field 'fieldName'
// with ints from the slice.
func QueryIntSlice(u []fet.Updater, slice []int, fieldName string) []fet.Updater {
    if slice != nil {
        elements := []fet.Updater{}
        for _, s := range slice {
            elements = append(elements, fet.Field(fieldName).Eq(s, fet.IfNotNil))
        }
        u = append(u, fet.Or(elements...))
    }
    return u
}

// Query document by given 'dates' slice. It matches document field 'fieldName'
// with dates from the slice. We only check if the year, month and day match
// and ignore rest of the date.
func QuerySameDay(u []fet.Updater, dates []time.Time, fieldName string) []fet.Updater {
    if dates != nil {
        elements := []fet.Updater{}
        for _, d := range dates {
            tempElement := []fet.Updater{}
            tempElement = append(tempElement, fet.Field(fieldName).Gt(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC), fet.IfNotNil))
            tempElement = append(tempElement, fet.Field(fieldName).Lt(time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, time.UTC), fet.IfNotNil))
            elements = append(elements, fet.And(tempElement...))
        }
        u = append(u, fet.Or(elements...))
    }
    return u
}

// Query document by given 'date' time.Time. It matches document field 'fieldName'
// where date is greater than the given 'date'.
func QueryDateGreater(u []fet.Updater, date time.Time, fieldName string) []fet.Updater {
    if !date.IsZero() {
        u = append(u, fet.Field(fieldName).Gt(date, fet.IfNotNil))
    }
    return u
}

// Query document by given 'date' time.Time. It matches document field 'fieldName'
// where date is lesser than the given 'date'.
func QueryDateLess(u []fet.Updater, date time.Time, fieldName string) []fet.Updater {
    if !date.IsZero() {
        u = append(u, fet.Field(fieldName).Lt(date, fet.IfNotNil))
    }
    return u
}

// Query document by given 'i' int. It matches document field 'fieldName'
// where number is greater than the given 'i'.
func QueryIntGreater(u []fet.Updater, i int, fieldName string) []fet.Updater {
    if i > 0 {
        u = append(u, fet.Field(fieldName).Gt(i, fet.IfNotNil))
    }
    return u
}

// Query document by given 'i' int. It matches document field 'fieldName'
// where number is lesser than the given 'i'.
func QueryIntLess(u []fet.Updater, i int, fieldName string) []fet.Updater {
    if i > 0 {
        u = append(u, fet.Field(fieldName).Lt(i, fet.IfNotNil))
    }
    return u
}

// Query document by given 'b' bool. It matches document field 'fieldName'
// where bool is the same as the given 'b'.
func QueryBoolSlice(u []fet.Updater, slice []bool, fieldName string) []fet.Updater {
    if slice != nil {
        elements := []fet.Updater{}
        for _, b := range slice {
            elements = append(elements, fet.Field(fieldName).Eq(b, fet.IfNotNil))
        }
        u = append(u, fet.Or(elements...))
    }
    return u
}

// Set the limit of the number of documents returned.
func QueryOptsLimit(opts *options.FindOptions, limit int) *options.FindOptions {
    if limit > 0 {
        opts = opts.SetLimit(int64(limit))
    }
    return opts
}

// Sort returned documents on given 'sort' strings.
func QueryOptsSort(opts *options.FindOptions, sort []string) *options.FindOptions {
    if sort != nil {
        d := bson.D{}
        for _, s := range sort {
            n := 1
            if s[0] == '-' {
                n = -1
                s = s[1:]
            }
            d = append(d, bson.E{s, n})
        }
        opts = opts.SetSort(d)
    }
    return opts
}

// Set projection on returned documents on given 'field' strings.
func QueryOptsProjection(opts *options.FindOptions, field []string) *options.FindOptions {
    if field != nil {
        d := bson.D{}
        for _, f := range field {
            n := 1
            if f[0] == '-' {
                n = 0
                f = f[1:]
                if f == "id" {
                    f = "_id"
                }
            }
            d = append(d, bson.E{f, n})
        }
        opts = opts.SetProjection(d)
    }
    return opts
}

// Set projection on returned documents on given 'field' strings.
func QueryOneOptsProjection(opts *options.FindOneOptions, field []string) *options.FindOneOptions {
    if field != nil {
        d := bson.D{}
        for _, f := range field {
            n := 1
            if f[0] == '-' {
                n = 0
                f = f[1:]
                if f == "id" {
                    f = "_id"
                }
            }
            d = append(d, bson.E{f, n})
        }
        opts = opts.SetProjection(d)
    }
    return opts
}

// Custom unmarshaler for our time.Time format ("2.1.2006").
var CustomUnmarshaler = qs.NewUnmarshaler(&qs.UnmarshalOptions{
    UnmarshalerFactory: &unmarshalerFactory{qs.NewDefaultUnmarshalOptions().UnmarshalerFactory},
})

type unmarshalerFactory struct {
    orig qs.UnmarshalerFactory
}

var timeType = reflect.TypeOf(time.Time(time.Unix(0, 0)))

func (f *unmarshalerFactory) Unmarshaler(t reflect.Type, opts *qs.UnmarshalOptions) (qs.Unmarshaler, error) {
    switch t {
    case timeType:
        return timeMarshalerInstance, nil
    default:
        return f.orig.Unmarshaler(t, opts)
    }
}

var timeMarshalerInstance = &timeMarshaler{}

type timeMarshaler struct{}

func (o *timeMarshaler) Unmarshal(v reflect.Value, a []string, opts *qs.UnmarshalOptions) error {
    s, err := opts.SliceToString(a)
    if err != nil {
    }
    t, err := time.Parse("2.1.2006", s)
    if err != nil {

    }
    v.Set(reflect.ValueOf(t))
    return nil
}
