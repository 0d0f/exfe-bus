package routex

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-rest"
	"logger"
	"model"
	"net/http"
	"routex/model"
	"sort"
	"time"
)

func (m *RouteMap) getTutorialData(currentTime time.Time, userId int64, number int) []rmodel.SimpleLocation {
	data, ok := m.tutorialDatas[userId]
	if !ok {
		return nil
	}
	zhCN := time.FixedZone("zh-CN", int(8*time.Hour/time.Second))
	currentTime = currentTime.In(zhCN)
	now := currentTime.Unix()
	todayTime, _ := time.Parse("2006-01-02 15:04:05 -0700", currentTime.Format("2006-01-02 00:00:00 -0700"))
	today := todayTime.Unix()
	offset := (now - today) / 10 * 10

	current := sort.Search(len(data), func(i int) bool {
		return data[i].Offset >= offset
	})
	if data[current].Offset != offset && number == 1 {
		return nil
	}

	var ret []rmodel.SimpleLocation
	for ; number > 0; number-- {
		l := rmodel.SimpleLocation{
			Timestamp: today + data[current].Offset,
			GPS:       [3]float64{data[current].Latitude, data[current].Longitude, data[current].Accuracy},
		}
		l.ToEarth(m.conversion)
		ret = append(ret, l)
		if current > 0 {
			current--
		} else {
			break
		}
	}
	return ret
}

type BreadcrumbOffset struct {
	Latitude  float64 `json:"earth_to_mars_latitude"`
	Longitude float64 `json:"earth_to_mars_longitude"`
}

func (o BreadcrumbOffset) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"earth_to_mars_latitude":%.6f,"earth_to_mars_longitude":%.6f}`, o.Latitude, o.Longitude)), nil
}

func (m RouteMap) UpdateBreadcrums(ctx rest.Context, breadcrumbs []rmodel.SimpleLocation) {
	var token rmodel.Token
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}

	ctx.Request().URL.RawQuery += fmt.Sprintf("&user_id=%d", token.UserId)
	m.UpdateBreadcrumsInner(ctx, breadcrumbs)
}

func (m RouteMap) UpdateBreadcrumsInner(ctx rest.Context, breadcrumbs []rmodel.SimpleLocation) {
	var ret BreadcrumbOffset
	var userId int64
	var coordinate string

	fmt.Println("url:", ctx.Request().URL.String())
	ctx.Bind("user_id", &userId)
	ctx.Bind("coordinate", &coordinate)
	fmt.Println(ctx.BindError())
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	breadcrumb := breadcrumbs[0]
	mars, earth := breadcrumb, breadcrumb

	if coordinate == "mars" {
		breadcrumb.ToEarth(m.conversion)
		earth = breadcrumb
	} else {
		mars.ToMars(m.conversion)
	}
	lat, lng, acc := breadcrumb.GPS[0], breadcrumb.GPS[1], breadcrumb.GPS[2]
	if acc <= 0 {
		ctx.Return(http.StatusBadRequest, "invalid accuracy: %f", acc)
		return
	}

	breadcrumb.Timestamp = time.Now().Unix()
	distance := float64(-1)
	if acc <= 70 {
		distance = 100
		if last, err := m.breadcrumbCache.Load(userId); err == nil {
			lastLat, lastLng := last.GPS[0], last.GPS[1]
			distance = Distance(lat, lng, lastLat, lastLng)
		}
	}
	var crossIds []int64
	var err error
	action := ""
	if distance > 30 {
		action = "save_to_history"
		logger.INFO("routex", "user", userId, "breadcrumb", fmt.Sprintf("%.7f", lat), fmt.Sprintf("%.7f", lng), acc, "distance", fmt.Sprintf("%.2f", distance), "save")
		if crossIds, err = m.breadcrumbCache.SaveCross(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
			ctx.Return(http.StatusInternalServerError, err)
			return
		}
		if err := m.breadcrumbCache.Save(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
		}
		if err := m.breadcrumbsRepo.Save(userId, breadcrumb); err != nil {
			logger.ERROR("can't save user %d breadcrumb: %s with %+v", userId, err, breadcrumb)
		}
	} else {
		logger.INFO("routex", "user", userId, "breadcrumb", fmt.Sprintf("%.7f", lat), fmt.Sprintf("%.7f", lng), acc, "distance", fmt.Sprintf("%.2f", distance), "nosave")
		if crossIds, err = m.breadcrumbCache.SaveCross(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
			ctx.Return(http.StatusInternalServerError, err)
			return
		}
		if acc <= 70 {
			if err := m.breadcrumbsRepo.UpdateLast(userId, breadcrumb); err != nil {
				logger.ERROR("can't update user %d breadcrumb: %s with %+v", userId, err, breadcrumb)
			}
		}
	}

	ret = BreadcrumbOffset{
		Latitude:  mars.GPS[0] - earth.GPS[0],
		Longitude: mars.GPS[1] - earth.GPS[1],
	}

	route := m.breadcrumbsToGeomark(userId, 1, []rmodel.SimpleLocation{breadcrumb})
	route.Action = action
	for _, cross := range crossIds {
		m.pubsub.Publish(m.publicName(cross), route)
	}
	ctx.Render(ret)
}

func (m RouteMap) GetBreadcrums(ctx rest.Context) {
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}
	var coordinate string
	ctx.Bind("coordinate", &coordinate)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	toMars := coordinate == "mars"
	breadcrumbs := m.getBreadcrumbs(token.Cross, toMars)
	ctx.Render(breadcrumbs)
}

func (m RouteMap) getBreadcrumbs(cross model.Cross, toMars bool) []rmodel.Geomark {
	var ret []rmodel.Geomark
	for _, invitation := range cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		marks := m.getUserBreadcrumbs(cross, userId, time.Now(), toMars)
		if len(marks) > 0 {
			ret = append(ret, marks...)
		}
	}
	return ret
}

func (m RouteMap) getUserBreadcrumbs(cross model.Cross, userId int64, after time.Time, toMars bool) []rmodel.Geomark {
	fmt.Println("cross", cross.ID, "user", userId, "after", after.Unix(), "tomars", toMars)
	var locations []rmodel.SimpleLocation
	if locations = m.getTutorialData(after, userId, 360); locations == nil {
		var err error
		if locations, err = m.breadcrumbsRepo.Load(userId, int64(cross.ID), after.Unix()); err != nil {
			logger.ERROR("can't get user %d breadcrumbs of cross %d: %s", userId, cross.ID, err)
			return nil
		}
	}
	fmt.Println(len(locations))
	if len(locations) == 0 {
		return nil
	}
	mark := m.breadcrumbsToGeomark(userId, 1, locations)
	if toMars {
		mark.ToMars(m.conversion)
	}
	return []rmodel.Geomark{mark}
}

func (m RouteMap) GetUserBreadcrums(ctx rest.Context) {
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}

	var cooridnate string
	var userId int64
	var afterFlag bool
	ctx.Bind("coordinate", &cooridnate)
	ctx.Bind("user_id", &userId)
	ctx.Bind("after_timestamp", &afterFlag)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	toMars := cooridnate == "mars"
	inCross := false
	for _, invitation := range token.Cross.Exfee.Invitations {
		if invitation.Identity.UserID == userId {
			inCross = true
			break
		}
	}
	if !inCross {
		ctx.Return(http.StatusUnauthorized, "user %d not in cross %d", userId, token.Cross.ID)
		return
	}

	after := time.Now().UTC()
	if afterFlag {
		var afterTimestamp int64
		ctx.BindReset()
		ctx.Bind("after_timestamp", &afterTimestamp)
		if err := ctx.BindError(); err != nil {
			ctx.Return(http.StatusBadRequest, err)
			return
		}
		after = time.Unix(afterTimestamp, 0)
	}
	breadcrumbs := m.getUserBreadcrumbs(token.Cross, userId, after, toMars)
	ctx.Render(breadcrumbs)
}

func (m RouteMap) GetUserBreadcrumsInner(ctx rest.Context) {
	var coordinate string
	var userId int64
	ctx.Bind("coordinate", &coordinate)
	ctx.Bind("user_id", &userId)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	toMars := coordinate == "mars"

	l, err := m.breadcrumbCache.Load(userId)
	if err != nil {
		if err == redis.ErrNil {
			ctx.Return(http.StatusNotFound, "can't find any breadcrumbs")
		} else {
			logger.ERROR("can't get user %d breadcrumbs: %s", userId, err)
			ctx.Return(http.StatusInternalServerError, err)
		}
		return
	}
	ret := m.breadcrumbsToGeomark(userId, 1, []rmodel.SimpleLocation{l})
	if toMars {
		ret.ToMars(m.conversion)
	}
	ctx.Render(ret)
}

func (m RouteMap) breadcrumbsId(userId int64) string {
	return fmt.Sprintf("breadcrumbs.%d", userId)
}

func (m RouteMap) breadcrumbsToGeomark(userId int64, windowId int, positions []rmodel.SimpleLocation) rmodel.Geomark {
	oldest := positions[len(positions)-1]
	r := rmodel.Geomark{
		Id:        fmt.Sprintf("breadcrumbs.%d", userId),
		Type:      "route",
		CreatedBy: fmt.Sprintf("%d@exfe", userId),
		CreatedAt: oldest.Timestamp,
		UpdatedBy: fmt.Sprintf("%d@exfe", userId),
		UpdatedAt: positions[0].Timestamp,
		Tags:      []string{"breadcrumbs"},
		Positions: positions,
	}
	return r
}
