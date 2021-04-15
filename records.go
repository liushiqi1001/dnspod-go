package dnspod

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	methodRecordList   = "Record.List"
	methodRecordCreate = "Record.Create"
	methodRecordInfo   = "Record.Info"
	methodRecordRemove = "Record.Remove"
	methodRecordModify = "Record.Modify"
	methodRecordRemark = "Record.Remark"
)

// Record is the DNS record representation.
type Record struct {
	ID            json.Number `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Line          string      `json:"line,omitempty"`
	LineID        string      `json:"line_id,omitempty"`
	Type          string      `json:"type,omitempty"`
	TTL           string      `json:"ttl,omitempty"`
	Value         string      `json:"value,omitempty"`
	MX            string      `json:"mx,omitempty"`
	Enabled       string      `json:"enabled,omitempty"`
	Status        string      `json:"status,omitempty"`
	MonitorStatus string      `json:"monitor_status,omitempty"`
	Remark        string      `json:"remark,omitempty"`
	UpdateOn      string      `json:"updated_on,omitempty"`
	UseAQB        string      `json:"use_aqb,omitempty"`
}

type recordsWrapper struct {
	Status  Status     `json:"status"`
	Info    DomainInfo `json:"info"`
	Records []Record   `json:"records"`
}

type recordWrapper struct {
	Status Status     `json:"status"`
	Info   DomainInfo `json:"info"`
	Record Record     `json:"record"`
}

// RecordsService handles communication with the DNS records related methods of the dnspod API.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html
type RecordsService struct {
	client *Client
}

type RecordListRequest struct {
	CommonParams
	DomainId     string
	Domain       string
	Offset       string
	Length       string
	SubDomain    string
	RecordType   string
	RecordLine   string
	RecordLIneId string
	KeyWord      string
}

func (r *RecordListRequest) toPayLOad() url.Values {
	p := r.CommonParams.toPayLoad()
	if r.DomainId != "" {
		p.Add("domain_id",r.DomainId)
	}
	if r.Domain != "" {
		p.Add("domain",r.Domain)
	}
	if r.Offset != "" {
		p.Add("offset",r.Offset)
	}
	if r.Length != "" {
		p.Add("length",r.Length)
	}
	if r.SubDomain != "" {
		p.Add("sub_domain",r.SubDomain)
	}
	if r.RecordType != "" {
		p.Add("record_type",r.RecordType)
	}
	if r.RecordLine != "" {
		p.Add("record_line",r.RecordLine)
	}
	if r.RecordLIneId != "" {
		p.Add("record_line_id" ,r.RecordLIneId)
	}
	if r.KeyWord != "" {
		p.Add("keyword",r.KeyWord)
	}
	return p
}

// List List the domain records.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-list
func (s *RecordsService) List(request *RecordListRequest) ([]Record, *Response, error) {
	payload := request.toPayLOad()

	wrappedRecords := recordsWrapper{}

	res, err := s.client.post(methodRecordList, payload, &wrappedRecords)
	if err != nil {
		return nil, res, err
	}

	if wrappedRecords.Status.Code != "1" {
		return nil, nil, fmt.Errorf("could not get domains: %s", wrappedRecords.Status.Message)
	}

	return wrappedRecords.Records, res, nil
}

// Create Creates a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-create
func (s *RecordsService) Create(domain string, recordAttributes Record) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domain)

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordCreate, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Get Fetches the domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-info
func (s *RecordsService) Get(domain string, recordID string) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domain)
	payload.Add("record_id", recordID)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordInfo, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Update Updates a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-modify
func (s *RecordsService) Update(domain string, recordID string, recordAttributes Record) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domain)

	if recordAttributes.ID != "" {
		payload.Add("record_id", string(recordAttributes.ID))
	}

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.ID == "" && recordID != "" {
		payload.Add("record_id", recordID)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordModify, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Remark Updates a domain record Remark.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-remark
func (s *RecordsService) Remark(domain string, recordID string, recordAttributes Record) (Record, *Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domain)

	if recordAttributes.ID != "" {
		payload.Add("record_id", string(recordAttributes.ID))
	}

	if recordAttributes.ID == "" && recordID != "" {
		payload.Add("record_id", recordID)
	}

	if recordAttributes.Remark != "" {
		payload.Add("remark", recordAttributes.Remark)
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordRemark, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// Delete Deletes a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-remove
func (s *RecordsService) Delete(domain string, recordID string) (*Response, error) {
	payload := s.client.CommonParams.toPayLoad()
	payload.Add("domain_id", domain)
	payload.Add("record_id", recordID)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(methodRecordRemove, payload, &returnedRecord)
	if err != nil {
		return res, err
	}

	if returnedRecord.Status.Code != "1" {
		return nil, fmt.Errorf("could not get domains: %s", returnedRecord.Status.Message)
	}

	return res, nil
}
