package helpers_test

import (
	helpers2 "github.com/PawelKowalski99/gogapps/providers/nasa/helpers"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetValidJsonField(t *testing.T) {
	type In struct {
		json        string
		queryValues map[string][]string
		field       string
	}

	tests := []struct {
		name   string
		in     In
		expOut string
	}{
		{
			name:   "TC_01 Get url without queries",
			in:     In{
				json: `{
					"copyright": "Fritz\nHelmut Hemmerich",
					"date": "2022-08-07",
					"explanation": "What's that green streak in front of the Andromeda galaxy? A meteor. While photographing the Andromeda galaxy in 2016, near the peak of the Perseid Meteor Shower, a small pebble from deep space crossed right in front of our Milky Way Galaxy's far-distant companion. The small meteor took only a fraction of a second to pass through this 10-degree field.  The meteor flared several times while braking violently upon entering Earth's atmosphere.  The green color was created, at least in part, by the meteor's gas glowing as it vaporized. Although the exposure was timed to catch a Perseid meteor, the orientation of the imaged streak seems a better match to a meteor from the Southern Delta Aquariids, a meteor shower that peaked a few weeks earlier.  Not coincidentally, the Perseid Meteor Shower peaks later this week, although this year the meteors will have to outshine a sky brightened by a nearly full moon.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_1948.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "Meteor before Galaxy",
					"url": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg"
				}`,
				queryValues: map[string][]string{},
				field: "url",
			},
			expOut: "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg",
		},
		{
			name:   "TC_02 Get url with one query",
			in:     In{
				json: `{
					"copyright": "Fritz\nHelmut Hemmerich",
					"date": "2022-08-07",
					"explanation": "What's that green streak in front of the Andromeda galaxy? A meteor. While photographing the Andromeda galaxy in 2016, near the peak of the Perseid Meteor Shower, a small pebble from deep space crossed right in front of our Milky Way Galaxy's far-distant companion. The small meteor took only a fraction of a second to pass through this 10-degree field.  The meteor flared several times while braking violently upon entering Earth's atmosphere.  The green color was created, at least in part, by the meteor's gas glowing as it vaporized. Although the exposure was timed to catch a Perseid meteor, the orientation of the imaged streak seems a better match to a meteor from the Southern Delta Aquariids, a meteor shower that peaked a few weeks earlier.  Not coincidentally, the Perseid Meteor Shower peaks later this week, although this year the meteors will have to outshine a sky brightened by a nearly full moon.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_1948.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "Meteor before Galaxy",
					"url": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg"
				}`,
				queryValues: map[string][]string{
					"media_type": {"image"},
				},
				field: "url",
			},
			expOut: "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg",
		},
		{
			name:   `TC_03 Get url with two queries and \n`,
			in:     In{
				json: `{
					"copyright": "Fritz\nHelmut Hemmerich",
					"date": "2022-08-07",
					"explanation": "What's that green streak in front of the Andromeda galaxy? A meteor. While photographing the Andromeda galaxy in 2016, near the peak of the Perseid Meteor Shower, a small pebble from deep space crossed right in front of our Milky Way Galaxy's far-distant companion. The small meteor took only a fraction of a second to pass through this 10-degree field.  The meteor flared several times while braking violently upon entering Earth's atmosphere.  The green color was created, at least in part, by the meteor's gas glowing as it vaporized. Although the exposure was timed to catch a Perseid meteor, the orientation of the imaged streak seems a better match to a meteor from the Southern Delta Aquariids, a meteor shower that peaked a few weeks earlier.  Not coincidentally, the Perseid Meteor Shower peaks later this week, although this year the meteors will have to outshine a sky brightened by a nearly full moon.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_1948.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "Meteor before Galaxy",
					"url": "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg"
				}`,
				queryValues: map[string][]string{
					"media_type": {"image"},
					"copyright": {`Fritz\nHelmut Hemmerich`},
				},
				field: "url",
			},
			expOut: "https://apod.nasa.gov/apod/image/2208/MeteorM31_hemmerich_960.jpg",
		},
		{
			name:   "TC_03 Get url with two queries",
			in:     In{
				json: `{
					"copyright": "Joan Josep Isach Cogollos",
					"date": "2022-08-04",
					"explanation": "In 1716, English astronomer Edmond Halley noted, \"This is but a little Patch, but it shows itself to the naked Eye, when the Sky is serene and the Moon absent.\" Of course, M13 is now less modestly recognized as the Great Globular Cluster in Hercules, one of the brightest globular star clusters in the northern sky. Sharp telescopic views like this one reveal the spectacular cluster's hundreds of thousands of stars. At a distance of 25,000 light-years, the cluster stars crowd into a region 150 light-years in diameter. Approaching the cluster core upwards of 100 stars could be contained in a cube just 3 light-years on a side. For comparison, the closest star to the Sun is over 4 light-years away. The remarkable range of brightness recorded in this image follows stars into the dense cluster core. Distant background galaxies in the medium-wide field of view include NGC 6207 at the upper left.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "M13: The Great Globular Cluster in Hercules",
					"url": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma1024.jpg"
					}`,
				queryValues: map[string][]string{
					"media_type": {"image"},
					"copyright": {`Joan Josep Isach Cogollos`},
				},
				field: "url",
			},
			expOut: "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma1024.jpg",
		},
		{
			name:   "TC_04 Get date with two queries",
			in:     In{
				json: `{
					"copyright": "Joan Josep Isach Cogollos",
					"date": "2022-08-04",
					"explanation": "In 1716, English astronomer Edmond Halley noted, \"This is but a little Patch, but it shows itself to the naked Eye, when the Sky is serene and the Moon absent.\" Of course, M13 is now less modestly recognized as the Great Globular Cluster in Hercules, one of the brightest globular star clusters in the northern sky. Sharp telescopic views like this one reveal the spectacular cluster's hundreds of thousands of stars. At a distance of 25,000 light-years, the cluster stars crowd into a region 150 light-years in diameter. Approaching the cluster core upwards of 100 stars could be contained in a cube just 3 light-years on a side. For comparison, the closest star to the Sun is over 4 light-years away. The remarkable range of brightness recorded in this image follows stars into the dense cluster core. Distant background galaxies in the medium-wide field of view include NGC 6207 at the upper left.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "M13: The Great Globular Cluster in Hercules",
					"url": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma1024.jpg"
					}`,
				queryValues: map[string][]string{
					"media_type": {"image"},
					"copyright": {`Joan Josep Isach Cogollos`},
				},
				field: "date",
			},
			expOut: "2022-08-04",
		},
		{
			name:   "TC_04 Get empty with two queries, value not valid",
			in:     In{
				json: `{
					"copyright": "Joan Josep Isach Cogollos",
					"date": "2022-08-04",
					"explanation": "In 1716, English astronomer Edmond Halley noted, \"This is but a little Patch, but it shows itself to the naked Eye, when the Sky is serene and the Moon absent.\" Of course, M13 is now less modestly recognized as the Great Globular Cluster in Hercules, one of the brightest globular star clusters in the northern sky. Sharp telescopic views like this one reveal the spectacular cluster's hundreds of thousands of stars. At a distance of 25,000 light-years, the cluster stars crowd into a region 150 light-years in diameter. Approaching the cluster core upwards of 100 stars could be contained in a cube just 3 light-years on a side. For comparison, the closest star to the Sun is over 4 light-years away. The remarkable range of brightness recorded in this image follows stars into the dense cluster core. Distant background galaxies in the medium-wide field of view include NGC 6207 at the upper left.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "M13: The Great Globular Cluster in Hercules",
					"url": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma1024.jpg"
					}`,
				queryValues: map[string][]string{
					"media_type": {"embed"},
					"copyright": {`Joan Josep Isach Cogollos`},
				},
				field: "date",
			},
			expOut: "",
		},
		{
			name:   "TC_04 Get empty with two queries, key not valid",
			in:     In{
				json: `{
					"copyright": "Joan Josep Isach Cogollos",
					"date": "2022-08-04",
					"explanation": "In 1716, English astronomer Edmond Halley noted, \"This is but a little Patch, but it shows itself to the naked Eye, when the Sky is serene and the Moon absent.\" Of course, M13 is now less modestly recognized as the Great Globular Cluster in Hercules, one of the brightest globular star clusters in the northern sky. Sharp telescopic views like this one reveal the spectacular cluster's hundreds of thousands of stars. At a distance of 25,000 light-years, the cluster stars crowd into a region 150 light-years in diameter. Approaching the cluster core upwards of 100 stars could be contained in a cube just 3 light-years on a side. For comparison, the closest star to the Sun is over 4 light-years away. The remarkable range of brightness recorded in this image follows stars into the dense cluster core. Distant background galaxies in the medium-wide field of view include NGC 6207 at the upper left.",
					"hdurl": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma.jpg",
					"media_type": "image",
					"service_version": "v1",
					"title": "M13: The Great Globular Cluster in Hercules",
					"url": "https://apod.nasa.gov/apod/image/2208/M13_final2_sinfirma1024.jpg"
					}`,
				queryValues: map[string][]string{
					"media_type": {"embed"},
					"name": {`Joan`},
				},
				field: "date",
			},
			expOut: "",
		},
	}

		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := helpers2.GetValidJsonField(tc.in.json, tc.in.queryValues, tc.in.field)

			if out != tc.expOut {
				t.Errorf("want: %s, got: %s", tc.expOut, out)
			}
		})
	}
}
