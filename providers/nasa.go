package providers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/gogoapps/providers/responses"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"math"
	"net/http"
	"storj.io/common/sync2"
	"sync"
	"time"
)

const (
	ApiKeyDefault="DEMO_KEY"
	NASA = "NASA"
	timeLayout = "2006-01-02"
)

type Nasa struct{
	ApiKey string
	ConcurrentRequests int
	L *logrus.Logger
}

func (n *Nasa) GetPictures(ctx context.Context) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		reqQuery := r.URL.Query()

		// Get from and to.
		// Rest is for unknown query filtering. From, to del is needed.
		fromTimeString := reqQuery.Get("from")
		toTimeString := reqQuery.Get("to")
		reqQuery.Del("from")
		reqQuery.Del("to")

		// Default values for time.Now date
		var from = time.Now()
		var to = time.Now().AddDate(0,0,1)
		var daysDiff = 1

		// Parse query dates
		if fromTimeString != "" && toTimeString != "" {
			// Parse values to time
			from, err = time.Parse(timeLayout, fromTimeString)
			if err != nil {
				return
			}
			to, err  = time.Parse(timeLayout, toTimeString)
			if err != nil {
				return
			}
			// Check if from is before to
			ok := from.Before(to)
			if !ok {
				render.JSON(w,r, fmt.Sprintf(`{"error": "QUERY: from is after to"}`))
				return
			}

			// Get days counter
			hoursDiff := to.Sub(from).Hours()
			daysDiff = int(math.Floor(hoursDiff / 24)) + 1
		}

		lim := sync2.NewLimiter(n.ConcurrentRequests)
		defer lim.Wait()

		// Make pictures safe slice so goroutine can append picture in json
		pictures := safeStringSlice{mu: sync.Mutex{}, data: []string{}}

		// Make request for each day concurrently
		for day := 0; day<daysDiff; day++ {

			apodTime := from.AddDate(0, 0, day).Format(timeLayout)

			started := lim.Go(ctx, func() {
				n.PictureRequest(w, r, apodTime, ctx, &pictures)
			})
			if !started {
				render.Render(w,r, responses.ErrConcurrentlyRun(http.StatusInternalServerError, errors.New("could not run concurrent func")))
			}

		}

		pictureUrls := safeStringSlice{mu: sync.Mutex{}, data: []string{}}

		var g errgroup.Group

		url := make(chan string)
		for _, picture := range pictures.data {
			g.Go(func() error {
				// Verify each picture json for query filtering
				return isValidPicture(picture, reqQuery, url)
			})
			g.Go(func() error{
				urlStr := <- url
				if urlStr != "" {
					// Add url to response
					pictureUrls.Writer(urlStr)
				}
				return nil
			})
		}

		err = g.Wait()
		if err != nil {
			render.Render(w,r,responses.ErrConcurrentlyRun(http.StatusInternalServerError, err))
		}

		render.Render(w, r, responses.PictureUrls(pictureUrls.Reader()))

		return
	}
}

func (n *Nasa) PictureRequest(w http.ResponseWriter, r *http.Request, date string, ctx context.Context, pictures *safeStringSlice) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.nasa.gov/planetary/apod"), nil)
	if err != nil {
		n.L.Errorf("could not create request: %v", err)
		pictures.Writer("")
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}
	q := req.URL.Query()
	q.Add("date", date)
	q.Add("api_key", n.ApiKey)
	req.URL.RawQuery = q.Encode()


	client := &http.Client{}

	rsp, err := client.Do(req)
	if err != nil {
		n.L.Errorf("could not make request: %v", err)
		pictures.Writer("")
		render.JSON(w, r, fmt.Sprintf(`{"error": "could not make request: %v"}`, err))
		return
	}
	if rsp.StatusCode != http.StatusOK {
		n.L.Errorf("http status not ok")
		pictures.Writer("")
		render.Render(w, r, responses.ErrStatusNotOk(rsp.StatusCode,
			fmt.Sprintf("url: %s, error: %s", req.URL.String(), rsp.Status)))
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		n.L.Errorf("could not marshal body to json: %v", err)
		pictures.Writer("")
		render.Render(w, r, responses.ErrParseBody(http.StatusInternalServerError, err))
		return
	}

	pictures.Writer(string(body))
	return
}

// isValidPicture lets query json picture due to its fields.
// If query got url == XXX AND copyright == XXX it checks if is right
func isValidPicture(picture string, query map[string][]string, url chan string) error {
	queryCounter := 0

	for key, values := range query {
		for _, value := range values {
			if gjson.Get(picture, key).String() == value {
				queryCounter++
			}
		}
	}
	if queryCounter == len(query) {
		url <- gjson.Get(picture, "url").String()
		return nil
	}
	url <- ""
	return nil
}


type safeStringSlice struct {
	mu   sync.Mutex
	data []string
}

func (o *safeStringSlice) Reader() []string {
	o.mu.Lock()
	urls := o.data
	o.mu.Unlock()
	return urls
}

func (o *safeStringSlice) Writer(data string) {
	o.mu.Lock()
	o.data = append(o.data, data)
	o.mu.Unlock()
}