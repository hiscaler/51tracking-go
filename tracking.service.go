package tracking51

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

type trackingService service

type CreateTrackingRequest struct {
	TrackingNumber          string `json:"tracking_number"`                     //	包裹物流单号
	CourierCode             string `json:"courier_code"`                        //	物流商对应的唯一简码
	OrderNumber             string `json:"order_number,omitempty"`              // 包裹的订单号，由商家/平台所产生的订单编号
	Title                   string `json:"title,omitempty"`                     // 包裹名称
	DestinationCode         string `json:"destination_code,omitempty"`          // 目的国的二字简码
	LogisticsChannel        string `json:"logistics_channel,omitempty"`         // 自定义字段，用于填写物流渠道（比如某货代商）
	Note                    string `json:"note,omitempty"`                      // 备注
	CustomerName            string `json:"customer_name,omitempty"`             // 客户姓名
	CustomerEmail           string `json:"customer_email,omitempty"`            // 客户邮箱
	CustomerPhone           string `json:"customer_phone,omitempty"`            // 顾客接收短信的手机号码。手机号码的格式应该为：“+区号手机号码”（例子：+8612345678910）
	ShippingDate            string `json:"shipping_date,omitempty"`             // 包裹发货时间（例子：2020-09-17 16:51）
	TrackingShippingDate    string `json:"tracking_shipping_date,omitempty"`    // 包裹的发货时间，其格式为：YYYYMMDD，有部分的物流商（如 deutsch-post）需要这个参数（例子：20200102）
	TrackingPostalCode      string `json:"tracking_postal_code,omitempty"`      // 收件人所在地邮编，仅有部分的物流商（如 postnl-3s）需要这个参数
	TrackingDestinationCode string `json:"tracking_destination_code,omitempty"` // 目的国对应的二字简码，部分物流商（如postnl-3s）需要这个参数
	TrackingCourierAccount  string `json:"tracking_courier_account,omitempty"`  // 物流商的官方账号，仅有部分的物流商（如 dynamic-logistics）需要这个参数
}

func (m CreateTrackingRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
		validation.Field(&m.CourierCode, validation.Required.Error("物流商简码不能为空")),
		validation.Field(&m.CustomerEmail, validation.When(m.CustomerEmail != "", is.EmailFormat.Error("客户邮箱地址格式错误"))),
		validation.Field(&m.CustomerPhone, validation.When(m.CustomerPhone != "", validation.Match(regexp.MustCompile(`^+\d{2}\d{11}$`)).Error("客户手机号码格式错误"))),
		validation.Field(&m.ShippingDate, validation.When(m.ShippingDate != "", validation.Date("2006-01-02 15:04").Error("包裹发货时间格式错误"))),
		validation.Field(&m.TrackingShippingDate, validation.When(m.TrackingShippingDate != "", validation.Date("20060102").Error("包裹发货时间格式错误"))),
	)
}

type Result struct {
	TrackingNumber string `json:"tracking_number"` //	包裹物流单号
	CourierCode    string `json:"courier_code"`    //	物流商对应的唯一简码
	OrderNumber    string `json:"order_number"`    //	包裹的订单号，由商家/平台所产生的订单编号
}

type CreateResult struct {
	Success []Result `json:"success"`
	Error   []Result `json:"error"`
}

func (s trackingService) Create(req CreateTrackingRequest) (res CreateResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().
		SetBody(req).
		Put("/create")
	if err != nil {
		return
	}

	r := struct {
		NormalResponse
		Data CreateResult `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &r); err == nil {
		res = r.Data
	}
	return
}

// 获取查询结果

type Track struct {
	TrackingNumber     string          `json:"tracking_number"`     // 包裹物流单号
	CourierCode        string          `json:"courier_code"`        // 物流商对应的唯一简码
	LogisticsChannel   string          `json:"logistics_channel"`   // 自定义字段，用于填写物流渠道（比如某货代商）
	Destination        string          `json:"destination"`         // 目的国的二字简码
	TrackUpdate        bool            `json:"track_update"`        // 自动更新查询功能的状态，“true”代表系统会自动更新查询结果，“false”则反之
	Consignee          string          `json:"consignee"`           // 签收人
	Updating           bool            `json:"updating"`            // “true”表示该单号会被继续更新，“false”表示该单号已停止更新
	CreatedAt          string          `json:"created_at"`          // 创建查询的时间
	UpdateDate         string          `json:"update_date"`         // 系统最后更新查询的时间
	OrderCreateTime    string          `json:"order_create_time"`   // 包裹发货时间
	CustomerEmail      string          `json:"customer_email"`      // 客户邮箱
	CustomerPhone      string          `json:"customer_phone"`      // 顾客接收短信的手机号码
	Title              string          `json:"title"`               // 包裹名称
	OrderNumber        string          `json:"order_number"`        // 包裹的订单号，由商家/平台所产生的订单编号
	Note               string          `json:"note"`                // 备注，可自定义
	CustomerName       string          `json:"customer_name"`       // 客户姓名
	Archived           bool            `json:"archived"`            // “true”表示该单号已被归档，“false”表示该单号处于未归档状态
	Original           string          `json:"original"`            // 发件国的名称
	DestinationCountry string          `json:"destination_country"` // 目的国的名称
	TransitTime        int             `json:"transit_time"`        // 包裹的从被揽收至被送达的时长（天）
	StayTime           int             `json:"stay_time"`           // 物流信息未更新的时长（单位：天），由当前时间减去物流信息最近更新时间得到
	OriginInfo         TrackOriginInfo `json:"origin_info"`         // 发件国的物流信息
}

type TrackOriginInfo struct {
	DestinationTrackNumber string      `json:"destination_track_number"` // 该包裹在目的国的物流单号
	ReferenceNumber        string      `json:"reference_number"`         // 包裹对应的另一个单号，作用与当前单号相同（仅有少部分物流商提供）
	ExchangeNumber         string      `json:"exchangeNumber"`           // 该包裹在中转站的物流商单号
	ReceivedDate           string      `json:"received_date"`            // 物流商接收包裹的时间（也称为上网时间）
	DispatchedDate         string      `json:"dispatched_date"`          // 包裹封发时间，封发指将多个小包裹打包成一个货物（方便运输）
	DepartedAirportDate    string      `json:"departed_airport_date"`    // 包裹离开此出发机场的时间
	ArrivedAbroadDate      string      `json:"arrived_abroad_date"`      // 包裹达到目的国的时间
	CustomsReceivedDate    string      `json:"customs_received_date"`    // 包裹移交给海关的时间
	ArrivedDestinationDate string      `json:"arrived_destination_date"` // 包裹达到目的国、目的城市的时间
	Weblink                string      `json:"weblink"`                  // 物流商的官网的链接
	CourierPhone           string      `json:"courier_phone"`            // 物流商官网上的电话
	TrackInfo              []TrackInfo `json:"trackinfo"`                // 详细物流信息
	ServiceCode            string      `json:"service_code"`             // 快递服务类型，比如次日达（部分物流商返回）
	StatusInfo             string      `json:"status_info"`              // 最新的一条物流信息
	Weight                 string      `json:"weight"`                   // 该货物的重量（多个包裹会被打包成一个“货物”）
	DestinationInfo        string      `json:"destination_info"`         // 目的国的物流信息
	LatestEvent            string      `json:"latest_event"`             // 最新物流信息的梗概，包括以下信息：状态、地址、时间
	LastestCheckpointTime  string      `json:"lastest_checkpoint_time"`  // 最新物流信息的更新时间
}

// TrackInfo 详细物流信息
type TrackInfo struct {
	CheckpointDate              string `json:"checkpoint_date"`               // 本条物流信息的更新时间，由物流商提供（包裹被扫描时，物流信息会被更新）
	TrackingDetail              string `json:"tracking_detail"`               // 具体的物流情况
	Location                    string `json:"location"`                      // 物流信息更新的地址（该包裹被扫描时，所在的地址）
	CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`    // 根据具体物流情况所识别出来的物流状态
	CheckpointDeliverySubStatus string `json:"checkpoint_delivery_substatus"` // 物流状态的子状态（物流状态）
}

type TracksQueryParams struct {
	TrackingNumbers string `json:"tracking_numbers"`  // 查询单号，每次不得超过40个，单号间以逗号分隔
	OrderNumbers    string `json:"order_numbers"`     // 订单号，每次查询不得超过40个，订单号间以逗号分隔
	DeliveryStatus  string `json:"delivery_status"`   // 状态
	ArchivedStatus  string `json:"archived_status"`   // 指定该单号是否被归档
	ItemsAmount     int    `json:"items_amount"`      // 每页展示的单号个数
	PagesAmount     int    `json:"pages_amount"`      // 返回结果的页数
	CreatedDateMin  int    `json:"created_date_min"`  // 创建查询的起始时间，时间戳格式
	CreatedDateMax  int    `json:"created_date_max"`  // 创建查询的结束时间，时间戳格式
	ShippingDateMin int    `json:"shipping_date_min"` // 发货的起始时间，时间戳格式
	ShippingDateMax int    `json:"shipping_date_max"` // 发货的结束时间，时间戳格式
	UpdatedDateMin  int    `json:"updated_date_min"`  // 查询更新的起始时间，时间戳格式
	UpdatedDateMax  int    `json:"updated_date_max"`  // 查询更新的结束时间，时间戳格式
	Lang            string `json:"lang"`              // 查询结果的语言（例子：cn, en），若未指定该参数，结果会以英文或中文呈现。 注意：只有物流商支持多语言查询结果时，该指定才会生效
}

func (m TracksQueryParams) Validate() error {
	return nil
}

func (s trackingService) All(params TracksQueryParams) (items []Track, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().
		SetQueryParamsFromValues(toValues(params)).
		Get("/get")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data []Track `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		items = res.Data
		isLastPage = len(items) < params.ItemsAmount
	}
	return
}
