package providers

import (
	"context"
	"fmt"
	"github.com/PawelKowalski99/gogapps/helpers"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/PawelKowalski99/gogapps/providers/responses"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"storj.io/common/sync2"
)

const (
	ApiKeyDefault = "DEMO_KEY"
	NASA          = "NASA"
	timeLayout    = "2006-01-02"
	urlKey          = "url"
)

type Nasa struct {
	ApiKey             string
	ConcurrentRequests int
	L                  *logrus.Logger
}

func (n *Nasa) GetPictures(ctx context.Context) http.HandlerFunc {

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
		var to = time.Now()
		var daysDiff = 0


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
			daysDiff = int(math.Floor(hoursDiff / 24))

			// Check if from is before to
			if daysDiff < 0 {
				 err = render.Render(w, r, responses.NewErrResponse("QUERY: from is after to", http.StatusInternalServerError))
				if err != nil {
					n.L.Errorf("could not render error: %v", err)
				}
				 return
			}

		}

		lim := sync2.NewLimiter(n.ConcurrentRequests)

		// Make pictures safe slice so goroutine can append picture in json
		pictures := safeStringSlice{mu: sync.Mutex{}, data: []string{}}

		// Make request for each day concurrently
		for day := 0; day <= daysDiff; day++ {

			apodTime := from.AddDate(0, 0, day).Format(timeLayout)

			started := lim.Go(ctx, func() {
				n.PictureRequest(w, r, apodTime, ctx, &pictures)
			})
			if !started {
				err = render.Render(w, r, responses.NewErrResponse(fmt.Sprintf("could not run concurrently"), http.StatusInternalServerError))
				if err != nil {
					n.L.Errorf("could not render error: %v", err)
				}
			}

		}
		lim.Wait()

		var pictureUrls []string

		for _, picture := range pictures.data {
			// Verify each picture json for query filtering
			if url := helpers.GetValidJsonField(picture, reqQuery, urlKey); url != "" {
				pictureUrls = append(pictureUrls, url)
			}
		}

		err = render.Render(w, r, responses.PictureUrls(pictureUrls))
		if err != nil {
			n.L.Errorf("could not render error: %v", err)
		}
		return
	}
}

func (n *Nasa) PictureRequest(w http.ResponseWriter, r *http.Request, date string, ctx context.Context, pictures *safeStringSlice) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.nasa.gov/planetary/apod"), nil)
	if err != nil {
		n.L.Errorf("could not create request: %v", err)
		pictures.Writer("")
		err = render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusBadRequest))
		if err != nil {
			n.L.Errorf("could not render error: %v", err)
		}
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
		err = render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusBadRequest))
		if err != nil {
			n.L.Errorf("could not render error: %v", err)
		}
		return
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		n.L.Errorf("could not marshal body to json: %v", err)
		pictures.Writer("")
		err =  render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusInternalServerError))
		if err != nil {
			n.L.Errorf("could not render error: %v", err)

		}
		return
	}

	if rsp.StatusCode == http.StatusNotFound {
		errMsg := gjson.Get(string(body), "msg")
		if errMsg.String() == fmt.Sprintf("%s %s", "No data available for date:", date) {
			err = render.Render(w, r, responses.NewErrResponse(errMsg.String(), http.StatusOK))
			if err != nil {
				n.L.Errorf("could not render error: %v", err)

			}
		}
	} else if rsp.StatusCode != http.StatusOK {
		n.L.Errorf("http status not ok")
		pictures.Writer("")
		err = render.Render(w, r, responses.NewErrResponse(fmt.Sprintf(`could not get date for date: %s`, date), http.StatusNotFound))
		if err != nil {
			n.L.Errorf("could not render error: %v", err)
		}
		return
	}
	defer rsp.Body.Close()

	pictures.Writer(string(body))
	return
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
