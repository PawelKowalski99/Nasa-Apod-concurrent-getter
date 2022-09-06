package nasa

import (
	"context"
	"fmt"
	"github.com/PawelKowalski99/gogapps/providers/nasa/helpers"
	"github.com/PawelKowalski99/gogapps/providers/nasa/responses"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"storj.io/common/sync2"
	"time"
)

const (
	urlKey = "url"
	APODURL = "https://api.nasa.gov/planetary/apod"
)

func (n *Nasa) GetPictures(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		helper := helpers.GetPictures{
			Pictures: helpers.InitSafeStringSlice(),
		}

		reqQuery, fromTimeString, toTimeString := helper.QueryParser(r)

		from, daysDiff := helper.GetTimeRanges(w, r, fromTimeString, toTimeString)

		// Run requests concurrently
		lim := sync2.NewLimiter(n.Config.GetConcurrentRequests())

		// Make request for each day concurrently
		for day := 0; day <= daysDiff; day++ {

			apodTime := helper.AddDaysToDate(from, day)

			started := lim.Go(ctx, func() {
				timeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				n.apodRequest(timeCtx, w, r, apodTime, &helper.Pictures)
			})
			if !started {
				err = render.Render(w, r, responses.NewErrResponse(fmt.Sprintf("could not run concurrently"), http.StatusInternalServerError))
				if err != nil {
					logrus.Errorf("could not render error: %v", err)
				}
			}

		}
		lim.Wait()

		var pictureUrls []string

		for _, picture := range helper.Pictures.Read() {
			// Verify each picture json for query filtering
			if url := helper.GetValidJsonField(picture, reqQuery, urlKey); url != "" {
				pictureUrls = append(pictureUrls, url)
			}
		}

		err = render.Render(w, r, responses.PictureUrls(pictureUrls))
		if err != nil {
			logrus.Errorf("could not render error: %v", err)
		}
		return
	}
}

func (n *Nasa) apodRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, date string,  pictures *helpers.SafeStringSlice) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APODURL, nil)
	if err != nil {
		logrus.Errorf("could not create request: %v", err)
		pictures.Write("")
		err = render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusBadRequest))
		if err != nil {
			logrus.Errorf("could not render error: %v", err)
		}
		return
	}
	q := req.URL.Query()
	q.Add("date", date)
	q.Add("api_key", n.Config.GetApiKey())
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	rsp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("could not make request: %v", err)
		pictures.Write("")
		err = render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusBadRequest))
		if err != nil {
			logrus.Errorf("could not render error: %v", err)
		}
		return
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		logrus.Errorf("could not marshal body to json: %v", err)
		pictures.Write("")
		err =  render.Render(w, r, responses.NewErrResponse(err.Error(), http.StatusInternalServerError))
		if err != nil {
			logrus.Errorf("could not render error: %v", err)

		}
		return
	}

	if rsp.StatusCode == http.StatusNotFound {
		errMsg := gjson.Get(string(body), "msg")
		if errMsg.String() == fmt.Sprintf("%s %s", "No data available for date:", date) {
			err = render.Render(w, r, responses.NewErrResponse(errMsg.String(), http.StatusOK))
			if err != nil {
				logrus.Errorf("could not render error: %v", err)

			}
		}
	} else if rsp.StatusCode != http.StatusOK {
		logrus.Errorf("http status not ok")
		pictures.Write("")
		err = render.Render(w, r, responses.NewErrResponse(fmt.Sprintf(`could not get date for date: %s`, date), http.StatusNotFound))
		if err != nil {
			logrus.Errorf("could not render error: %v", err)
		}
		return
	}
	defer rsp.Body.Close()

	pictures.Write(string(body))
	return
}
