package restagent

import "gitlab.utc.fr/milairhu/ia04-api-rest/restagent/comsoc"

type Request struct {
	Preferences []comsoc.Alternative `json:"pref"`
}

type Response struct {
	Result []comsoc.Alternative `json:"res"`
}
