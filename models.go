package tempodb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type tempoTime struct {
	Time time.Time
}

type dataPoint struct {
	Ts tempoTime `json:"t"`
	V  float64   `json:"v"`
}

type DataPoint struct {
	Ts time.Time
	V  float64
}

type bulkDataSet struct {
	Ts   tempoTime   `json:"t"`
	Data []BulkPoint `json:"data"`
}

type BulkDataSet struct {
	Ts   time.Time
	Data []BulkPoint
}

type BulkPoint interface {
	GetValue() float64
}

type BulkKeyPoint struct {
	Key string  `json:"key"`
	V   float64 `json:"v"`
}

type BulkIdPoint struct {
	Id string  `json:"id"`
	V  float64 `json:"v"`
}

type createSeriesRequest struct {
	Key string
}

type dataSet struct {
	Series  Series             `json:"series"`
	Start   tempoTime          `json:"start"`
	End     tempoTime          `json:"end"`
	Data    []*DataPoint       `json:"data"`
	Summary map[string]float64 `json:"summary"`
}

type DataSet struct {
	Series  Series
	Start   time.Time
	End     time.Time
	Data    []*DataPoint
	Summary map[string]float64
}

type Series struct {
	Id         string            `json:"id"`
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Tags       []string          `json:"tags"`
}

type Filter struct {
	Ids        []string
	Keys       []string
	Tags       []string
	Attributes map[string]string
}

func NewFilter() *Filter {
	return &Filter{
		Ids:        make([]string, 0),
		Keys:       make([]string, 0),
		Tags:       make([]string, 0),
		Attributes: make(map[string]string),
	}
}

func (tt tempoTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", tt.Time.Format(ISO8601_FMT))
	return []byte(formatted), nil
}

func (tt tempoTime) UnmarshalJSON(data []byte) error {
	b := bytes.NewBuffer(data)
	decoded := json.NewDecoder(b)
	var s string
	if err := decoded.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(ISO8601_FMT, s)
	if err != nil {
		return err
	}
	tt.Time = t

	return nil
}

func (dp *DataPoint) MarshalJSON() ([]byte, error) {
	ts := tempoTime{Time: dp.Ts}
	pdp := &dataPoint{Ts: ts, V: dp.V}
	return json.Marshal(pdp)
}

func (dp *DataPoint) UnmarshalJSON(data []byte) error {
	pdp := new(dataPoint)
	err := json.Unmarshal(data, pdp)
	if err != nil {
		return err
	}
	dp.Ts = pdp.Ts.Time
	dp.V = pdp.V

	return nil
}

func (bds *BulkDataSet) MarshalJSON() ([]byte, error) {
	ts := tempoTime{Time: bds.Ts}
	pbds := &bulkDataSet{Ts: ts, Data: bds.Data}
	return json.Marshal(pbds)
}

func (bds *BulkDataSet) UnmarshalJSON(data []byte) error {
	pbds := new(bulkDataSet)
	err := json.Unmarshal(data, pbds)
	if err != nil {
		return err
	}
	bds.Ts = pbds.Ts.Time
	bds.Data = pbds.Data

	return nil
}

func (ds *DataSet) MarshalJSON() ([]byte, error) {
	start := tempoTime{Time: ds.Start}
	end := tempoTime{Time: ds.End}
	pds := &dataSet{Start: start, End: end, Data: ds.Data, Series: ds.Series, Summary: ds.Summary}
	return json.Marshal(pds)
}

func (ds *DataSet) UnmarshalJSON(data []byte) error {
	pds := new(dataSet)
	err := json.Unmarshal(data, pds)
	if err != nil {
		return err
	}
	ds.Start = pds.Start.Time
	ds.End = pds.End.Time
	ds.Series = pds.Series
	ds.Data = pds.Data
	ds.Summary = pds.Summary

	return nil
}

func (filter *Filter) AddId(id string) {
	filter.Ids = append(filter.Ids, id)
}

func (filter *Filter) AddKey(key string) {
	filter.Keys = append(filter.Keys, key)
}

func (filter *Filter) AddTag(tag string) {
	filter.Tags = append(filter.Tags, tag)
}

func (filter *Filter) AddAttribute(key string, value string) {
	filter.Attributes[key] = value
}

func (bp *BulkKeyPoint) GetValue() float64 {
	return bp.V
}

func (bp *BulkIdPoint) GetValue() float64 {
	return bp.V
}

