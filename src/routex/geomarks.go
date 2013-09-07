package routex

import (
	"fmt"
	"logger"
	"model"
	"net/http"
	"net/url"
	"routex/model"
	"strconv"
	"strings"
	"time"
)

const (
	XPlaceTag      = "xplace"
	DestinationTag = "destination"
)

func (m RouteMap) setTutorial(lat, lng float64, userId, crossId int64, locale, by string) (rmodel.Geomark, error) {
	var ret rmodel.Geomark
	query := make(url.Values)
	query.Set("keyword", "attractions")
	places, err := m.platform.GetPlace(lat, lng, locale, 10000, query)
	if err != nil {
		return ret, err
	}
	if len(places) == 0 {
		places, err = m.platform.GetPlace(lat, lng, locale, 50000, nil)
		if err != nil {
			return ret, err
		}
	}
	if len(places) == 0 {
		return ret, fmt.Errorf("can't find attraction place near %.7f,%.7f", lat, lng)
	}
	place := places[0]
	if lng, err = strconv.ParseFloat(place.Lng, 64); err != nil {
		return ret, err
	}
	if lat, err = strconv.ParseFloat(place.Lat, 64); err != nil {
		return ret, err
	}
	now := time.Now().Unix()
	ret = rmodel.Geomark{
		Id:          fmt.Sprintf("location.%04d", m.rand.Intn(1e4)),
		Type:        "location",
		CreatedAt:   now,
		CreatedBy:   by,
		UpdatedAt:   now,
		UpdatedBy:   by,
		Tags:        []string{"destination"},
		Icon:        "",
		Title:       place.Title,
		Description: place.Description,
		Longitude:   lng,
		Latitude:    lat,
	}
	if err := m.geomarksRepo.Set(crossId, ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func (m RouteMap) HandleSearchGeomarks() []rmodel.Geomark {
	crossIdStr := m.Vars()["cross_id"]
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return nil
	}
	toMars := false
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		toMars = true
	}
	data, err := m.geomarksRepo.Get(crossId)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", crossId, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		return []rmodel.Geomark{}
	}

	tag := m.Request().URL.Query().Get("tags")

	var idMap map[string]bool
	if ids, ok := m.Request().URL.Query()["id"]; ok {
		idMap = make(map[string]bool)
		for _, id := range ids {
			idMap[id] = true
		}
	}

	ret := []rmodel.Geomark{}
	for _, geomark := range data {
		ok := true
		switch {
		case tag != "" && !geomark.HasTag(tag):
			ok = false
		case idMap != nil && !idMap[geomark.Id]:
			ok = false
		}
		if ok {
			if toMars {
				geomark.ToMars(m.conversion)
			}
			ret = append(ret, geomark)
		}
	}
	return ret
}

func (m RouteMap) HandleGetGeomarks() []rmodel.Geomark {
	_, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	return m.HandleSearchGeomarks()
}

func (m RouteMap) getGeomarks(cross model.Cross, toMars bool) ([]rmodel.Geomark, error) {
	data, err := m.geomarksRepo.Get(int64(cross.ID))
	if err != nil {
		return nil, err
	}

	hasDestination := false
	for i, d := range data {
		if d.HasTag(DestinationTag) {
			hasDestination = true
		}
		if toMars {
			d.ToMars(m.conversion)
			data[i] = d
		}
	}

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
		xplace := rmodel.Geomark{
			Id:          m.xplaceId(int64(cross.ID)),
			Type:        "location",
			CreatedAt:   createdAt.Unix(),
			CreatedBy:   cross.By.Id(),
			UpdatedAt:   updatedAt.Unix(),
			UpdatedBy:   cross.By.Id(),
			Tags:        []string{XPlaceTag},
			Icon:        "",
			Title:       cross.Place.Title,
			Description: cross.Place.Description,
			Longitude:   lng,
			Latitude:    lat,
		}
		if !hasDestination {
			xplace.Tags = append(xplace.Tags, "destination")
		}
		if !toMars {
			xplace.ToEarth(m.conversion)
		}
		data = append(data, xplace)
	}

	return data, nil
}

func (m RouteMap) HandleSetGeomark(mark rmodel.Geomark) {
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	mark.Type = m.Vars()["mark_type"]
	kind := m.Vars()["kind"]
	mark.Id = fmt.Sprintf("%s.%s", kind, m.Vars()["mark_id"])
	mark.UpdatedBy, mark.UpdatedAt, mark.Action = token.Identity.Id(), time.Now().Unix(), ""
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		mark.ToEarth(m.conversion)
	}

	if mark.HasTag(XPlaceTag) {
		if err := m.syncCrossPlace(&mark, token.Cross, mark.UpdatedBy); err != nil {
			logger.ERROR("can't set cross %d place: %s", token.Cross.ID, err)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	} else if kind == "location" || kind == "route" {
		if err := m.geomarksRepo.Set(int64(token.Cross.ID), mark); err != nil {
			logger.ERROR("save geomark %s %s error: %s", mark.Type, mark.Id, err)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	}

	mark.Action = "update"
	m.pubsub.Publish(m.publicName(int64(token.Cross.ID)), mark)
	m.checkGeomarks(token.Cross, mark)
	m.update(token.UserId, int64(token.Cross.ID))

	return
}

func (m RouteMap) HandleDeleteGeomark() {
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	var mark rmodel.Geomark
	mark.Type = m.Vars()["mark_type"]
	kind := m.Vars()["kind"]
	mark.Id = fmt.Sprintf("%s.%s", kind, m.Vars()["mark_id"])

	if mark.HasTag(XPlaceTag) {
		if err := m.syncCrossPlace(nil, token.Cross, token.Identity.Id()); err != nil {
			logger.ERROR("remove cross %d place error: %s", token.Cross.ID, err)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	}
	if kind == "location" || kind == "route" {
		if err := m.geomarksRepo.Delete(int64(token.Cross.ID), mark.Type, mark.Id); err != nil {
			logger.ERROR("delete geromark %s %s error: %s", mark.Type, mark.Id, err)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	}

	mark.Action = "delete"
	m.pubsub.Publish(m.publicName(int64(token.Cross.ID)), mark)
	m.checkGeomarks(token.Cross, mark)
	m.update(token.UserId, int64(token.Cross.ID))

	return
}

func (m RouteMap) checkGeomarks(cross model.Cross, mark rmodel.Geomark) {
	marks, _ := m.getGeomarks(cross, false)

	if mark.HasTag(DestinationTag) {
		if mark.Action == "update" {
			for _, mk := range marks {
				if mk.Id != mark.Id && mk.RemoveTag(DestinationTag) {
					m.geomarksRepo.Set(int64(cross.ID), mk)
					mk.Action = "update"
					m.pubsub.Publish(m.publicName(int64(cross.ID)), mk)
				}
			}
		}
		if mark.Action == "delete" {
			for _, mk := range marks {
				if mk.HasTag(XPlaceTag) {
					mk.Tags = append(mk.Tags, DestinationTag)
					mk.Action = "update"
					m.pubsub.Publish(m.publicName(int64(cross.ID)), mk)
				}
			}
		}
	}
	if mark.HasTag(XPlaceTag) && !strings.HasPrefix(mark.Id, XPlaceTag) {
		m.geomarksRepo.Delete(int64(cross.ID), mark.Type, mark.Id)
		mark.Action = "delete"
		m.pubsub.Publish(m.publicName(int64(cross.ID)), mark)
	}
}

func (m RouteMap) xplaceId(crossId int64) string {
	return fmt.Sprintf(XPlaceTag+".%d", crossId)
}

func (m RouteMap) syncCrossPlace(geomark *rmodel.Geomark, cross model.Cross, by string) error {
	updatedBy := model.FromIdentityId(by)
	place := model.Place{}

	if cross.Place != nil {
		place.ID = cross.Place.ID
	}
	if geomark != nil {
		p := *geomark
		p.ToMars(m.conversion)
		place.Title = p.Title
		place.Description = p.Description
		place.Lng = fmt.Sprintf("%.7f", p.Longitude)
		place.Lat = fmt.Sprintf("%.7f", p.Latitude)
		place.Provider = "routex"
		place.ExternalID = fmt.Sprintf("%d", cross.ID)
	}
	updateCross := map[string]interface{}{"place": place}
	return m.platform.BotCrossUpdate("cross_id", fmt.Sprintf("%d", cross.ID), updateCross, updatedBy)
}
