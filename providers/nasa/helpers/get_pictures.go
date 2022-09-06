package helpers

import (
	"github.com/PawelKowalski99/gogapps/providers/nasa/responses"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	timeLayout = "2006-01-02"
	hoursInDay = 24
)

type GetPictures struct {
	Pictures SafeStringSlice
}

func (n *GetPictures) QueryParser(r *http.Request) (url.Values, string, string) {
	reqQuery := r.URL.Query()
	// Get from and to.
	// Rest is for unknown query filtering. From, to del is needed.
	fromTimeString := reqQuery.Get("from")
	toTimeString := reqQuery.Get("to")
	reqQuery.Del("from")
	reqQuery.Del("to")
	return reqQuery, fromTimeString, toTimeString
}

func (n *GetPictures) GetTimeRanges(w http.ResponseWriter, r *http.Request, fromTimeString, toTimeString string) (from time.Time, daysDiff int) {
	// Default values for time.Now date
	var err error
	from = time.Now()
	to := time.Now()

	// Parse query dates
	if fromTimeString != "" && toTimeString != "" {
		// Parse values to time
		from, err = time.Parse(timeLayout, fromTimeString)
		if err != nil {
			return
		}
		to, err = time.Parse(timeLayout, toTimeString)
		if err != nil {
			return
		}

		// Get days counter
		hoursDiff := to.Sub(from).Hours()
		daysDiff = int(math.Floor(hoursDiff / hoursInDay))

		// Check if from is before to
		if daysDiff < 0 {
			err = render.Render(w, r, responses.NewErrResponse("QUERY: from is after to", http.StatusInternalServerError))
			if err != nil {
				logrus.Errorf("could not render error: %v", err)
			}
			return
		}

	}

	return from, daysDiff
}

func (n *GetPictures) AddDaysToDate(date time.Time, day int) string {
	return date.AddDate(0, 0, day).Format(timeLayout)
}

// GetValidJsonField lets query json picture due to its fields.
// If query got url == XXX AND copyright == XXX it checks if is right
func (n *GetPictures) GetValidJsonField(json string, query map[string][]string, field string) string {
	queryCounter := 0
	for key, values := range query {

		connectedValues := strings.Join(values, " ")
		if strings.ReplaceAll(gjson.Get(json, key).Raw, `"`, "") == connectedValues {
			queryCounter++
		}
	}

	if queryCounter == len(query) {
		return gjson.Get(json, field).String()
	}
	return ""
}



type SafeStringSlice struct {
	mu   sync.Mutex
	data []string
}

func InitSafeStringSlice() SafeStringSlice {
	return SafeStringSlice{
		mu:   sync.Mutex{},
		data: []string{},
	}
}

func (o *SafeStringSlice) Read() []string {
	o.mu.Lock()
	urls := o.data
	o.mu.Unlock()
	return urls
}

func (o *SafeStringSlice) Write(data string) {
	o.mu.Lock()
	o.data = append(o.data, data)
	o.mu.Unlock()
}

