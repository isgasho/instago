package instago

import (
	"encoding/json"
	"errors"
	"strings"
)

const urlReelsMedia = `https://i.instagram.com/api/v1/feed/reels_media/`

type highlightMedia struct {
	Reels struct {
		ReelsMedia IGStoryHighlightsTray `json:"reels_media"`
	} `json:"reels"`
	Status string `json:"status"`
}

func (m *IGApiManager) GetHighlightsReelsMedia(id string) (tray IGStoryHighlightsTray, err error) {
	url := urlReelsMedia + "?user_ids=" + id
	b, err := getHTTPResponse(url, m.dsUserId, m.sessionid, m.csrftoken)
	if err != nil {
		return
	}

	// The name of json field is the id of the highlight tray, which is only
	// known in run-time, not compile-time. So we need to replace the id of
	// the highlight tray with *reels_media*, which can be decoded by Go
	// standard encoding/json package.
	bb := []byte(strings.Replace(string(b), id, "reels_media", 1))

	h := highlightMedia{}
	err = json.Unmarshal(bb, &h)
	if err != nil {
		return
	}

	// Check validity
	if h.Reels.ReelsMedia.Id != id {
		err = errors.New("Returned highlight tray seems weird")
		return
	}
	tray = h.Reels.ReelsMedia
	return
}
