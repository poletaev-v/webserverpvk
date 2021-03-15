package tcpserver

import (
	"encoding/xml"
	"reflect"
	"strings"
)

// XMLParser for parsing xml file
type XMLParser struct{}

type Record struct {
	RecordInfo            RecordInfo            `xml:"RecordInfo"`
	VehicleInfo           VehicleInfo           `xml:"VehicleInfo"`
	WheelBaseMeasurements WheelBaseMeasurements `xml:"WheelBaseMeasurements"`
	TrackBaseMeasurements TrackBaseMeasurements `xml:"TrackBaseMeasurements"`
	LengthMeasurements    LengthMeasurements    `xml:"LengthMeasurements"`
	WidthMeasurements     WidthMeasurements     `xml:"WidthMeasurements"`
	HeightMeasurements    HeightMeasurements    `xml:"HeightMeasurements"`
	Speed
}

type RecordInfo struct {
	PvkCode                                               uint64         `xml:"PvkCode"`
	EquipmentName                                         string         `xml:"EquipmentName"`
	Place                                                 string         `xml:"Place"`
	HighwayName                                           string         `xml:"HighwayName"`
	PvkCoordinates                                        PvkCoordinates `xml:"PvkCoordinates"`
	CertificateStatementSuchMeasurementNumber             string         `xml:"CertificateStatementSuchMeasurementNumber"`
	CertificateStatementSuchMeasurementDate               string         `xml:"CertificateStatementSuchMeasurementDate"`
	CertificateStatementSuchMeasurementRegistrationNumber string         `xml:"CertificateStatementSuchMeasurementRegistrationNumber"`
	CheckingDocNumber                                     string         `xml:"CheckingDocNumber"`
	CheckingCertificateDate                               string         `xml:"CheckingCertificateDate"`
	CheckingValid                                         string         `xml:"CheckingValid"`
	IDBetamount                                           uint64         `xml:"IDBetamount"`
	ExcessFactDate                                        string         `xml:"ExcessFactDate"`
	ElectronicStamp                                       string         `xml:"ElectronicStamp"`
}

type PvkCoordinates struct {
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
}

type VehicleInfo struct {
	PlatformID            uint64 `xml:"PlatformId"`
	TrackStateNumber      string `xml:"TrackStateNumber"`
	RearStateNumber       string `xml:"RearStateNumber"`
	TrackCategory         uint64 `xml:"TrackCategory"`
	TrackCategoryRus12    uint64 `xml:"TrackCategoryRus12"`
	TrackSubCategoryRus12 uint64 `xml:"TrackSubCategoryRus12"`
	TrackAxes             uint64 `xml:"TrackAxes"`
}

type WheelBaseMeasurements struct {
	MeasuredWheelBase12 uint64 `xml:"MeasuredWheelBase12"`
	TrackWheelBase12    uint64 `xml:"TrackWheelBase12"`
	NormWheelBase12     uint64 `xml:"NormWheelBase12"`
	MeasuredWheelBase23 uint64 `xml:"MeasuredWheelBase23"`
	TrackWheelBase23    uint64 `xml:"TrackWheelBase23"`
	NormWheelBase23     uint64 `xml:"NormWheelBase23"`
	MeasuredWheelBase34 uint64 `xml:"MeasuredWheelBase34"`
	TrackWheelBase34    uint64 `xml:"TrackWheelBase34"`
	NormWheelBase34     uint64 `xml:"NormWheelBase34"`
	MeasuredWheelBase45 uint64 `xml:"MeasuredWheelBase45"`
	TrackWheelBase45    uint64 `xml:"TrackWheelBase45"`
	NormWheelBase45     uint64 `xml:"NormWheelBase45"`
}

type TrackBaseMeasurements struct {
	WeightSign                 bool   `xml:"WeightSign"`
	MeasuredTrackGrossWeight   uint64 `xml:"MeasuredTrackGrossWeight"`
	TrackGrossWeight           uint64 `xml:"TrackGrossWeight"`
	NormGrossWeight            uint64 `xml:"NormGrossWeight"`
	ThrustSign                 bool   `xml:"ThrustSign"`
	MeasuredThrust1            uint64 `xml:"MeasuredThrust1"`
	TrackThrust1               uint64 `xml:"TrackThrust1"`
	NormThrust1                uint64 `xml:"NormThrust1"`
	DifferenceNormTrackThrust1 uint64 `xml:"DifferenceNormTrackThrust1"`
	TrackWheels1               uint64 `xml:"TrackWheels1"`
	TrackWheelsEx1             uint64 `xml:"TrackWheelsEx1"`
	MeasuredThrust2            uint64 `xml:"MeasuredThrust2"`
	TrackThrust2               uint64 `xml:"TrackThrust2"`
	NormThrust2                uint64 `xml:"NormThrust2"`
	DifferenceNormTrackThrust2 uint64 `xml:"DifferenceNormTrackThrust2"`
	TrackWheels2               uint64 `xml:"TrackWheels2"`
	TrackWheelsEx2             uint64 `xml:"TrackWheelsEx2"`
	GroupAxlesNumber2          uint64 `xml:"GroupAxlesNumber2"`
	MeasuredThrust3            uint64 `xml:"MeasuredThrust3"`
	TrackThrust3               uint64 `xml:"TrackThrust3"`
	NormThrust3                uint64 `xml:"NormThrust3"`
	DifferenceNormTrackThrust3 uint64 `xml:"DifferenceNormTrackThrust3"`
	TrackWheels3               uint64 `xml:"TrackWheels3"`
	TrackWheelsEx3             uint64 `xml:"TrackWheelsEx3"`
	GroupAxlesNumber3          uint64 `xml:"GroupAxlesNumber3"`
	MeasuredThrust4            uint64 `xml:"MeasuredThrust4"`
	TrackThrust4               uint64 `xml:"TrackThrust4"`
	NormThrust4                uint64 `xml:"NormThrust4"`
	DifferenceNormTrackThrust4 uint64 `xml:"DifferenceNormTrackThrust4"`
	TrackWheels4               uint64 `xml:"TrackWheels4"`
	TrackWheelsEx4             uint64 `xml:"TrackWheelsEx4"`
	MeasuredThrust5            uint64 `xml:"MeasuredThrust5"`
	TrackThrust5               uint64 `xml:"TrackThrust5"`
	NormThrust5                uint64 `xml:"NormThrust5"`
	DifferenceNormTrackThrust5 uint64 `xml:"DifferenceNormTrackThrust5"`
	TrackWheels5               uint64 `xml:"TrackWheels5"`
	TrackWheelsEx5             uint64 `xml:"TrackWheelsEx5"`
	GroupAxlesCount1           uint64 `xml:"GroupAxlesCount1"`
	GroupAxlesMeasuredWeight1  uint64 `xml:"GroupAxlesMeasuredWeight1"`
	GroupAxlesWeight1          uint64 `xml:"GroupAxlesWeight1"`
	GroupAxlesNorm1            uint64 `xml:"GroupAxlesNorm1"`
	GroupAxlesSign1            bool   `xml:"GroupAxlesSign1"`
}

type LengthMeasurements struct {
	LengthSign                bool   `xml:"LengthSign"`
	MeasuredTrackLength       uint64 `xml:"MeasuredTrackLength"`
	TrackLength               uint64 `xml:"TrackLength"`
	NormLength                uint64 `xml:"NormLength"`
	DifferenceTrackNormLength uint64 `xml:"DifferenceTrackNormLength"`
}

type WidthMeasurements struct {
	WidthSign                bool   `xml:"WidthSign"`
	MeasuredTrackWidth       uint64 `xml:"MeasuredTrackWidth"`
	TrackWidth               uint64 `xml:"TrackWidth"`
	NormWidth                uint64 `xml:"NormWidth"`
	DifferenceTrackNormWidth uint64 `xml:"DifferenceTrackNormWidth"`
}

type HeightMeasurements struct {
	HeightSign                bool   `xml:"HeightSign"`
	MeasuredTrackHeight       uint64 `xml:"MeasuredTrackHeight"`
	TrackWHeight              uint64 `xml:"TrackWHeight"`
	NormHeight                uint64 `xml:"NormHeight"`
	DifferenceTrackNormHeight uint64 `xml:"DifferenceTrackNormHeight"`
}

type Speed struct {
	Speed uint64
}

func (xp *XMLParser) parse(data []byte) (Record, error) {
	var values Record

	if err := xml.Unmarshal(data, &values); err != nil {
		return values, err
	}
	return values, nil
}

func (xp *XMLParser) getValues(data Record) map[string]interface{} {
	values := make(map[string]interface{})
	e := reflect.ValueOf(&data).Elem()
	xp.getRecValue(e, e.Type().Field(0), values)

	return values
}

func (xp XMLParser) getRecValue(e reflect.Value, f reflect.StructField, values map[string]interface{}) {
	if e.Kind() == reflect.Struct {
		for i := 0; i < e.NumField(); i++ {
			xp.getRecValue(e.Field(i), e.Type().Field(i), values)
		}
	} else {
		values[strings.ToLower(f.Name)] = e.Interface()
	}
}
