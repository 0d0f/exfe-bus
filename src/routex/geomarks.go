package routex

import (
	"fmt"
	"logger"
	"model"
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
	toMars := false
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		toMars = true
	}
	ret, err := m.getGeomarks(token.Cross, toMars)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	return ret
}

func (m RouteMap) getGeomarks(cross model.Cross, toMars bool) ([]Geomark, error) {
	data, err := m.geomarksRepo.Get(int64(cross.ID))
	if err != nil {
		return nil, err
	}

	needCrossPlace := true
	hasDestination := false
	for i, d := range data {
		for _, t := range d.Tags {
			if t == CrossPlaceTag {
				needCrossPlace = false
			}
			if t == "destination" {
				hasDestination = true
			}
		}
		if toMars {
			d.ToMars(m.conversion)
			data[i] = d
		}
	}

	if needCrossPlace {
		var lat, lng float64
		if cross.Place != nil {
			if lng, err = strconv.ParseFloat(cross.Place.Lng, 64); err != nil {
				cross.Place = nil
			} else if lat, err = strconv.ParseFloat(cross.Place.Lat, 64); err != nil {
				cross.Place = nil
			}
		}
		if cross.Place != nil {
			createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", cross.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", cross.CreatedAt)
			if err != nil {
				updatedAt = time.Now()
			}
			xplace := Geomark{
				Id:          fmt.Sprintf("%08d@location", m.rand.Intn(1e8)),
				Type:        "location",
				CreatedAt:   createdAt.Unix(),
				CreatedBy:   cross.By.Id(),
				UpdatedAt:   updatedAt.Unix(),
				UpdatedBy:   cross.By.Id(),
				Tags:        []string{CrossPlaceTag},
				Icon:        "http://panda.0d0f.com/static/img/map_pin_blue@2x.png",
				Title:       cross.Place.Title,
				Description: cross.Place.Description,
				Longitude:   lng,
				Latitude:    lat,
			}
			if !hasDestination {
				xplace.Tags = append(xplace.Tags, "destination")
			}
			data = append(data, xplace)
		}
	}

	return data, nil
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

	for _, tag := range mark.Tags {
		if tag == CrossPlaceTag {
			m.Error(http.StatusBadRequest, fmt.Errorf("mark should not include %s tag", CrossPlaceTag))
			return
		}
	}

	mark.Type, mark.Id, mark.UpdatedAt, mark.Action = m.Vars()["mark_type"], m.Vars()["mark_id"], time.Now().Unix(), ""
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		mark.ToEarth(m.conversion)
	}
	if err := m.geomarksRepo.Set(int64(token.Cross.ID), mark); err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}
	m.routexRepo.Update(token.UserId, int64(token.Cross.ID))

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
	m.routexRepo.Update(token.UserId, int64(token.Cross.ID))

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
