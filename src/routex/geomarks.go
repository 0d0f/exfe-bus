package routex

import (
	"logger"
	"net/http"
	"strconv"
	"time"
)

const CrossPlaceTag = "xplace"

func (m RouteMap) HandleSearchGeomarks() []Geomark {
	crossIdStr := m.Vars()["cross_id"]
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return nil
	}
	ret := make([]Geomark, 0)
	tag := m.Request().URL.Query().Get("tags")
	if tag == "" {
		return ret
	}
	data, err := m.geomarksRepo.Get(crossId)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", crossId, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		return ret
	}
	for _, geomark := range data {
		ok := false
		for _, t := range geomark.Tags {
			if t == tag {
				ok = true
			}
			if t == CrossPlaceTag {
				ok = false
				break
			}
		}
		if ok {
			ret = append(ret, geomark)
		}
	}
	return ret
}

func (m RouteMap) HandleGetGeomarks() []Geomark {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth(true)
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	data, err := m.geomarksRepo.Get(int64(token.Cross.ID))
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		var lat, lng float64
		if token.Cross.Place != nil {
			if lng, err = strconv.ParseFloat(token.Cross.Place.Lng, 64); err != nil {
				token.Cross.Place = nil
			} else if lat, err = strconv.ParseFloat(token.Cross.Place.Lat, 64); err != nil {
				token.Cross.Place = nil
			}
		}
		if token.Cross.Place != nil {
			createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
			if err != nil {
				updatedAt = time.Now()
			}
			destinaion := Geomark{
				Id:          "destination",
				Type:        "location",
				CreatedAt:   createdAt.Unix(),
				CreatedBy:   token.Cross.By.Id(),
				UpdatedAt:   updatedAt.Unix(),
				UpdatedBy:   token.Cross.By.Id(),
				Tags:        []string{"destination", CrossPlaceTag},
				Icon:        "http://panda.0d0f.com/static/img/map_pin_blue@2x.png",
				Title:       token.Cross.Place.Title,
				Description: token.Cross.Place.Description,
				Longitude:   lng,
				Latitude:    lat,
			}
			go func() {
				m.geomarksRepo.Set(int64(token.Cross.ID), destinaion)
			}()
			data = []Geomark{destinaion}
		}
	}

	for i, d := range data {
		for _, t := range d.Tags {
			if t == CrossPlaceTag && token.Cross.Place != nil {
				d.Latitude, _ = strconv.ParseFloat(token.Cross.Place.Lat, 64)
				d.Longitude, _ = strconv.ParseFloat(token.Cross.Place.Lng, 64)
				break
			}
		}
		if m.Request().URL.Query().Get("coordinate") == "mars" {
			d.ToMars(m.conversion)
		}
		data[i] = d
	}
	return data
}

func (m RouteMap) HandleSetGeomark(mark Geomark) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth(true)
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	mark.Type, mark.Id, mark.UpdatedAt, mark.Action = m.Vars()["mark_type"], m.Vars()["mark_id"], time.Now().Unix(), ""
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		mark.ToEarth(m.conversion)
	}
	if err := m.geomarksRepo.Set(int64(token.Cross.ID), mark); err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	go func() {
		mark.Action = "update"
		m.castLocker.RLock()
		broadcast := m.crossCast[int64(token.Cross.ID)]
		m.castLocker.RUnlock()
		if broadcast != nil {
			broadcast.Send(mark)
		}
	}()

	return
}

func (m RouteMap) HandleDeleteGeomark() {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth(true)
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	var mark Geomark
	mark.Type, mark.Id = m.Vars()["mark_type"], m.Vars()["mark_id"]
	if err := m.geomarksRepo.Delete(int64(token.Cross.ID), mark.Type, mark.Id); err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	go func() {
		mark.Action = "delete"
		m.castLocker.RLock()
		broadcast := m.crossCast[int64(token.Cross.ID)]
		m.castLocker.RUnlock()
		if broadcast != nil {
			broadcast.Send(mark)
		}
	}()

	return
}
