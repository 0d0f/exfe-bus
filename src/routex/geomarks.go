package routex

import (
	"fmt"
	"logger"
	"model"
	"net/http"
	"strconv"
	"strings"
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
				Id:          m.xplaceId(int64(cross.ID)),
				Type:        "location",
				CreatedAt:   createdAt.Unix(),
				CreatedBy:   cross.By.Id(),
				UpdatedAt:   updatedAt.Unix(),
				UpdatedBy:   cross.By.Id(),
				Tags:        []string{CrossPlaceTag},
				Icon:        "",
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

	mark.Type, mark.Id, mark.UpdatedAt, mark.Action = m.Vars()["mark_type"], m.Vars()["mark_id"], time.Now().Unix(), ""
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		mark.ToEarth(m.conversion)
	}

	for i := len(mark.Tags) - 1; i >= 0; i-- {
		if mark.Tags[i] == CrossPlaceTag {
			go func() {
				if err := m.syncCrossPlace(&mark, int64(token.Cross.ID)); err != nil {
					logger.ERROR("can't set cross %d place: %s", token.Cross.ID, err)
				} else {
					if err := m.geomarksRepo.Delete(int64(token.Cross.ID), mark.Type, mark.Id); err != nil {
						logger.ERROR("can't delete cross %d geomark %s %s: %s", token.Cross.ID, mark.Type, mark.Id, err)
					}
					m.routexRepo.Update(token.UserId, int64(token.Cross.ID))

					mark.Action = "delete"
					m.castLocker.RLock()
					broadcast := m.crossCast[int64(token.Cross.ID)]
					m.castLocker.RUnlock()
					if broadcast != nil {
						broadcast.Send(mark)
						mark.Id, mark.Action = m.xplaceId(int64(token.Cross.ID)), ""
						broadcast.Send(mark)
					}
					return
				}
			}()
			return
		}
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
	if strings.HasSuffix(mark.Id, ".location") || strings.HasSuffix(mark.Id, ".route") {
		if err := m.geomarksRepo.Delete(int64(token.Cross.ID), mark.Type, mark.Id); err != nil {
			m.Error(http.StatusInternalServerError, err)
			return
		}
	}
	m.routexRepo.Update(token.UserId, int64(token.Cross.ID))

	go func() {
		if mark.Id == m.xplaceId(int64(token.Cross.ID)) {
			if err := m.syncCrossPlace(nil, int64(token.Cross.ID)); err != nil {
				logger.ERROR("remove cross %d place error: %s", token.Cross.ID, err)
			}
		}
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

func (m RouteMap) xplaceId(crossId int64) string {
	return fmt.Sprintf("%d."+CrossPlaceTag, crossId)
}

func (m RouteMap) syncCrossPlace(geomark *Geomark, crossId int64) error {
	updatedBy := model.FromIdentityId(geomark.UpdatedBy)
	cross := model.Cross{}
	if geomark != nil {
		cross.Place = &model.Place{
			Title:       geomark.Title,
			Description: geomark.Description,
			Lng:         fmt.Sprintf("%.7f", geomark.Longitude),
			Lat:         fmt.Sprintf("%.7f", geomark.Latitude),
			Provider:    "routex",
			ExternalID:  fmt.Sprintf("%d", crossId),
		}
	}
	return m.platform.BotCrossUpdate("cross_id", fmt.Sprintf("%d", crossId), cross, updatedBy)
}
