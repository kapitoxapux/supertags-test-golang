package service

type Fact struct {
	PeriodStart       string     `json:"period_start"`
	PeriodEnd         string     `json:"period_end"`
	PeriodKey         string     `json:"period_key"`
	IndicatorToMoId   uint32     `json:"indicator_to_mo_id"`
	IndicatorToFactId uint32     `json:"indicator_to_fact_id"`
	Value             uint8      `json:"value"`
	FactTime          string     `json:"fact_time"`
	IsPlan            uint8      `json:"is_plan"`
	Supertags         []Supertag `json:"supertags"`
	AuthUserId        uint32     `json:"auth_user_id"`
	Comments          []Row      `json:"comment"`
}

type Supertag struct {
	Tag   Tag    `json:"tag"`
	Value string `json:"value"`
}

type Tag struct {
	Id           uint8  `json:"id"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	ValuesSource uint8  `json:"values_source"`
}

type Row struct {
	Id      string `json:"_id"`
	Key     string `json:"_key"`
	Rev     string `json:"_rev"`
	Authors Author `json:"author"`
	Group   string `json:"group"`
	Msg     string `json:"msg"`
	Params  Param  `json:"params"`
	Time    string `json:"time"`
	Type    string `json:"type"`
}

type Author struct {
	MoId     uint32 `json:"mo_id"`
	UserId   uint32 `json:"user_id"`
	UserName string `json:"user_name"`
}

type Param struct {
	IndicatorToMoId uint32 `json:"indicator_to_mo_id"`
	Periods         Period `json:"period"`
	Platform        string `json:"platform"`
}

type Period struct {
	End     string `json:"end"`
	Start   string `json:"start"`
	TypeId  uint32 `json:"type_id"`
	TypeKey string `json:"type_key"`
}

type Storage interface {
	InsertFact(fact *Fact) (*Fact, error)
}

type Service struct {
	Storage Storage
}

func NewService(storage Storage) *Service {

	return &Service{
		Storage: storage,
	}
}
