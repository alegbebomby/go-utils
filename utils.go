package library

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"log"
	"math"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const notSetError = "is not set"
const StandardDateFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

var (
	once      sync.Once
	netClient *http.Client
)

func GetString(payload map[string]interface{}, name string, defaults string) (string, error) {

	if payload[name] == nil {

		return defaults, errors.New(fmt.Sprintf("value %s not set", name))
	}

	v := reflect.ValueOf(payload[name])

	switch v.Kind() {

	case reflect.Invalid:
		return defaults, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int()), nil
		//return strconv.FormatInt(v.Int(), 10), nil

	case reflect.Float64, reflect.Float32:
		return fmt.Sprintf("%0.f", v.Float()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil

	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil

	case reflect.String:
		return v.String(), nil

	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value", nil
	}
}

func GetFloat(payload map[string]interface{}, name string, defaults float64) (float64, error) {

	if payload[name] == nil {

		return defaults, errors.New(fmt.Sprintf("value %s not set", name))
	}

	v := reflect.ValueOf(payload[name])

	switch v.Kind() {

	case reflect.Invalid:
		return defaults, errors.New(fmt.Sprintf("value %s not set", name))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), nil

	case reflect.Float64, reflect.Float32:
		return v.Float(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), nil
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		//return strconv.FormatBool(v.Bool())
		return defaults, errors.New(fmt.Sprintf("value %s not set", name))

	case reflect.String:
		val, err := strconv.ParseFloat(v.String(), 64)

		if err != nil {

			return defaults, err
		}

		return val, nil

	default: // reflect.Array, reflect.Struct, reflect.Interface
		//return v.Type().String() + " value"
		return defaults, errors.New(fmt.Sprintf("value %s not set", name))
	}

}

func GetInt64(payload map[string]interface{}, name string, defaults int64) (int64, error) {

	if payload[name] == nil {

		return defaults, errors.New(fmt.Sprintf("value %s not set", name))
	}

	v := reflect.ValueOf(payload[name])

	switch v.Kind() {

	case reflect.Invalid:
		return defaults, errors.New(fmt.Sprintf(notSetError))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil

	case reflect.Float64, reflect.Float32:
		return int64(v.Float()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint()), nil
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		//return strconv.FormatBool(v.Bool())
		return defaults, errors.New(fmt.Sprintf(notSetError))

	case reflect.String:
		val, err := strconv.ParseInt(v.String(), 10, 64)

		if err != nil {

			return defaults, err
		}

		return val, nil

	default: // reflect.Array, reflect.Struct, reflect.Interface
		//return v.Type().String() + " value"
		return defaults, errors.New(fmt.Sprintf(notSetError))
	}

}

func GetBool(payload map[string]interface{}, name string, defaults bool) (bool, error) {

	if payload[name] == nil {

		return defaults, fmt.Errorf("value %s not set", name)
	}

	v := reflect.ValueOf(payload[name])

	switch v.Kind() {

	case reflect.Invalid:
		return defaults, errors.New(fmt.Sprintf(notSetError))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 1, nil

	case reflect.Float64, reflect.Float32:
		return int64(v.Float()) == 1, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint()) == 1, nil

	case reflect.Bool:
		return v.Bool(), nil

	case reflect.String:
		val, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {

			return strings.TrimSpace(strings.ToLower(v.String())) == "true", nil
		}

		return val == 1, nil

	default:
		return defaults, errors.New(fmt.Sprintf(notSetError))
	}
}

func GetInt64Value(payload interface{}, defaults int64) (int64, error) {

	if payload == nil {

		return defaults, errors.New(fmt.Sprintf(notSetError))
	}

	v := reflect.ValueOf(payload)

	switch v.Kind() {

	case reflect.Invalid:
		return defaults, errors.New(fmt.Sprintf(notSetError))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil

	case reflect.Float64, reflect.Float32:
		return int64(v.Float()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint()), nil
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return defaults, errors.New(fmt.Sprintf(notSetError))

	case reflect.String:
		val, err := strconv.ParseInt(v.String(), 10, 64)

		if err != nil {

			return defaults, err
		}

		return val, nil

	default: // reflect.Array, reflect.Struct, reflect.Interface
		//return v.Type().String() + " value"
		return defaults, errors.New(fmt.Sprintf(notSetError))
	}

}

func ToMysql(t time.Time) string {

	return t.Format(StandardDateFormat)
}

func CombinedDateTime(dateDate time.Time, timeString string) time.Time {

	dt := dateDate.Format(DateFormat)

	stringTime := fmt.Sprintf("%s %s", dt, timeString)

	newLayout := StandardDateFormat
	myTime, err := time.Parse(newLayout, stringTime)

	if err != nil {

		log.Printf("Got error converting string to time %s", err.Error())
		return dateDate
	}

	return myTime
}

func StringToTime(stringTime string) time.Time {

	if len(strings.TrimSpace(stringTime)) == 0 {

		return time.Now()
	}

	parts := strings.Split(stringTime," ")
	if len(parts) == 1 {

		stringTime = fmt.Sprintf("%s 00:00:00",parts[0])
	}

	newLayout := StandardDateFormat
	myTime, err := time.Parse(newLayout, stringTime)
	if err != nil {

		log.Printf("StringToTime error converting %s to time %s",stringTime,err.Error())
	}

	return myTime
}

func MysqlNow() string {

	return time.Now().Format(StandardDateFormat)
}

func Today() string {

	return time.Now().Format(DateFormat)
}

func CalculateTotalPages(total int, perPage int) int {

	if total == 0 || perPage == 0 {

		return 1
	}

	if total <= perPage {

		return 1
	}

	totalPages := math.Round(float64(total) / float64(perPage))

	if total > perPage && total % perPage > 0 {

		totalPages++

	}

	return int(totalPages)
}

/*
# For details see man 4 crontabs

# Example of job definition:
# .---------------- minute (0 - 59)
# |  .------------- hour (0 - 23)
# |  |  .---------- day of month (1 - 31)
# |  |  |  .------- month (1 - 12) OR jan,feb,mar,apr ...
# |  |  |  |  .---- day of week (0 - 6) (Sunday=0 or 7) OR sun,mon,tue,wed,thu,fri,sat
# |  |  |  |  |
# *  *  *  *  * user-name  command to be executed
*/

func CronString(repeatType string, repeatIntervalValue string, date string, sendTime string) (cron *string, error error) {

	repeatIntervalValue = strings.TrimSpace(repeatIntervalValue)

	if repeatType == "EVERY_MINUTE" {

		// expected format
		// * * * * *
		cron := fmt.Sprintf("* * * * *")
		return &cron, nil

	}

	if strings.Compare(repeatType, "EVERY_HOUR") == 0 {

		// expected format
		// min * * * *
		cron := fmt.Sprintf("%s * * * *", repeatIntervalValue)
		return &cron, nil

	}

	if strings.Compare(repeatType, "EVERY_DAY") == 0 {

		//sendTime := strings.TrimSpace(sendTime)

		parts := strings.Split(sendTime, ":")

		if len(parts) < 2 {

			err := fmt.Sprintf("invalid send time %s ", sendTime)
			log.Printf("got error parsing time error %s ", err)
			return nil, errors.New(err)
		}
		// expected format
		// min hour * * *
		cron := fmt.Sprintf("%s %s * * *", strings.TrimSpace(parts[1]), strings.TrimSpace(parts[0]))
		return &cron, nil
	}

	if len(date) == 0 {

		date = Today()
	}

	sendTime = fmt.Sprintf("%s:00", sendTime)

	mydate := fmt.Sprintf("%s %s", date, sendTime)

	myDate, err := time.Parse(StandardDateFormat, mydate)

	if err != nil {

		log.Printf("got error parsing date %s error %s ", mydate, err.Error())
		return nil, err
	}

	_, month, day := myDate.Date()
	hour, min, _ := myDate.Clock()

	if strings.Compare(repeatType, "NO_REPEAT") == 0 {

		// expected format
		// min hour day month *

		cron := fmt.Sprintf("%d %d %d %d *", min, hour, day, month)

		return &cron, nil
	}

	if strings.Compare(repeatType, "EVERY_WEEK") == 0 {

		// expected format
		// min hour * * daysOfWeek
		cron := fmt.Sprintf("%d %d * * %s", min, hour, repeatIntervalValue)
		return &cron, nil

	}

	if strings.Compare(repeatType, "EVERY_MONTH") == 0 {

		// expected format
		// min hour day 1-12 *
		cron := fmt.Sprintf("%d %d %d 1-12 *", min, hour, day)
		return &cron, nil
	}

	return nil, errors.New("invalid repeat type values passed")

}

func DateLayout() string {

	return DateFormat
}

func RemoveInvalidCharacters(message string) string {

	invalids := map[string]string{
		"/[áàâãªä]/u": "a",
		"/[ÁÀÂÃÄ]/u":  "A",
		"/[ÍÌÎÏ]/u":   "I",
		"/[íìîï]/u":   "i",
		"/[éèêë]/u":   "e",
		"/[ÉÈÊË]/u":   "E",
		"/[óòôõºö]/u": "o",
		"/[ÓÒÔÕÖ]/u":  "O",
		"/[úùûü]/u":   "u",
		"/[ÚÙÛÜ]/u":   "U",
		"/ç/":         "c",
		"/Ç/":         "C",
		"/ñ/":         "n",
		"/Ñ/":         "N",
		"/–/":         "-", // UTF-8 hyphen to "normal" hyphen
		"/[’‘‹›‚]/u":  " ", // Literally a single quote
		"/[“”«»„]/u":  " ", // Double quote
		"/ /":         " ", // nonbreaking space (equiv. to 0x160)
	}

	for k, v := range invalids {

		message = strings.Replace(message, k, v, -1)
	}

	re := regexp.MustCompile("[[:^ascii:]]")
	message = re.ReplaceAllLiteralString(message, "")

	return message
	///[[:^print:]]/
}

func ToHuman(t time.Time) string {

	return t.Format("3:04PM Mon Jan")
}

func ToMysqlDateTime(t time.Time) string {

	return t.Format(StandardDateFormat)

}

func ToMysqlDate(t time.Time) string {

	return t.Format(DateFormat)

}


func Contains(elems []string, elem string) bool {

	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func NextMonth(t time.Time) time.Time {

	myTime := now.New(t)
	myT := myTime.EndOfMonth()
	myTT := myT.Add(48 * time.Hour)
	return myTT
}

func NewNetClient() *http.Client {

	once.Do(func() {

		var netTransport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 60 * time.Second,
			}).Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 60 * time.Second,
		}

		netClient = &http.Client{
			Timeout:   time.Second * 60,
			Transport: netTransport,
		}
	})

	return netClient
}

func HTTPPost(url string, headers map[string] string, payload interface{}) (httpStatus int, response string) {

	if payload == nil {

		payload = "{}"
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if headers != nil {

		for k,v := range headers {

			req.Header.Set(k,v)
		}
	}

	resp, err := NewNetClient().Do(req)
	if err != nil {

		log.Printf("got error making http request %s", err.Error())
		return 0, ""
	}

	st := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Printf("got error making http request %s",err.Error())
		return st,""
	}

	return st, string(body)
}