package sport_api

type SportApi struct {
}

func NewSportApi() *SportApi {
	p := &SportApi{}
	p.Init()
	return p
}
func (self *SportApi) Init() {
}
