package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"

	"github.com/oursky/ourd/oddb"
	"github.com/oursky/ourd/oderr"
	"github.com/oursky/ourd/router"
)

// recordPayload is the input parameter in RecordHandler
type recordPayload router.Payload

func (p *recordPayload) IsWriteAllowed() bool {
	return p.UserInfo == nil
}

// transportRecord override JSON serialization and deserialization of
// oddb.Record
type transportRecord oddb.Record

func (r transportRecord) MarshalJSON() ([]byte, error) {
	// NOTE(limouren): if there is a better way to shallow copy a map,
	// do let me know
	object := map[string]interface{}{}
	for k, v := range r.Data {
		object[k] = v
	}
	object["_id"] = r.Type + "/" + r.Key
	object["_type"] = r.Type

	return json.Marshal(object)
}

func (r *transportRecord) UnmarshalJSON(data []byte) error {
	object := map[string]interface{}{}
	err := json.Unmarshal(data, &object)

	if err != nil {
		return err
	}

	return r.InitFromMap(object)
}

func (r *transportRecord) InitFromMap(m map[string]interface{}) error {
	rawID, ok := m["_id"].(string)
	if !ok {
		return errors.New(`record/json: required field "_id" not found`)
	}
	delete(m, "_id")

	ss := strings.SplitN(rawID, "/", 2)
	if len(ss) == 1 {
		return fmt.Errorf(`record/json: "_id" should be of format '{type}/{id}', got %#v`, rawID)
	}

	recordType, id := ss[0], ss[1]

	r.Key = id
	r.Type = recordType
	r.Data = m

	return nil
}

// idResponseItem encapsulates an item in a list of Record ID as response.
type idResponseItem struct {
	ID   string `json:"_id,omitempty"`
	Type string `json:"_type,omitempty"`
	Code string `json:"_code,omitempty"`
}

/*
RecordSaveHandler is dummy implementation on save/modify Records
curl -X POST -H "Content-Type: application/json" \
  -d @- http://192.168.1.89/ <<EOF
{
    "action": "record:save",
    "access_token": "validToken",
    "database_id": "private",
    "records": [{
        "_id": "EA6A3E68-90F3-49B5-B470-5FFDB7A0D4E8",
        "_type": "note",
        "content": "ewdsa"
    }]
}
EOF
*/
func RecordSaveHandler(payload *router.Payload, response *router.Response) {
	if (*recordPayload)(payload).IsWriteAllowed() {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: write is not allowed")
		return
	}

	db := payload.Database
	recordMaps, ok := payload.Data["records"].([]interface{})
	if !ok {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: expected list of records")
		return
	}

	length := len(recordMaps)

	records := make([]transportRecord, length, length)
	results := make([]interface{}, length, length)
	for i := range records {
		r := recordMaps[i].(map[string]interface{})
		if err := records[i].InitFromMap(r); err != nil {
			results[i] = oderr.New(oderr.RequestInvalidErr, "invalid request: "+err.Error())
		}
	}

	for i := range records {
		_, fail := results[i].(error)
		if !fail {
			if err := db.Save((*oddb.Record)(&records[i])); err != nil {
				results[i] = oderr.New(oderr.PersistentStorageErr, "persistent error: failed to save record")
			} else {
				results[i] = records[i]
			}
		}
	}

	response.Result = results
}

/*
RecordFetchHandler is dummy implementation on fetching Records
curl -X POST -H "Content-Type: application/json" \
  -d @- http://localhost:3000/ <<EOF
{
    "action": "record:fetch",
    "access_token": "validToken",
    "database_id": "private",
    "ids": ["1004", "1005"]
}
EOF
*/
func RecordFetchHandler(payload *router.Payload, response *router.Response) {
	interfaces, ok := payload.Data["ids"].([]interface{})
	if !ok {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: expect list of ids")
		return
	}

	length := len(interfaces)
	recordIDs := make([]string, length, length)
	for i, it := range interfaces {
		rawID, ok := it.(string)
		if !ok {
			response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: expect list of ids")
			return
		}

		ss := strings.SplitN(rawID, "/", 2)
		if len(ss) == 1 {
			response.Err = oderr.NewFmt(oderr.RequestInvalidErr, "invalid id format: %v", rawID)
			return
		}

		recordIDs[i] = ss[1]
	}

	db := payload.Database

	results := make([]interface{}, length, length)
	for i, recordID := range recordIDs {
		record := oddb.Record{}
		if err := db.Get(recordID, &record); err != nil {
			if err == oddb.ErrRecordNotFound {
				results[i] = idResponseItem{
					ID:   recordID,
					Type: "_error",
					Code: "NOT_FOUND",
				}
			} else {
				results[i] = idResponseItem{
					ID:   recordID,
					Type: "_error",
					Code: "UNKNOWN_ERR",
				}
			}
		} else {
			results[i] = record
		}
	}

	response.Result = results
}

func keyPathFromRaw(rawKeyPath map[string]interface{}) string {
	mapType := rawKeyPath["$type"]
	if mapType != "keypath" {
		panic(fmt.Errorf("got key path's type %v, want \"keypath\"", mapType))
	}

	keypath := rawKeyPath["$val"].(string)
	if keypath == "" {
		panic("empty key path value")
	}

	return keypath
}

func sortFromRaw(rawSort []interface{}, sort *oddb.Sort) {
	keyPathMap, _ := rawSort[0].(map[string]interface{})
	if len(keyPathMap) == 0 {
		panic(errors.New("empty key path in sort descriptor"))
	}
	keyPath := keyPathFromRaw(keyPathMap)

	orderStr, _ := rawSort[1].(string)
	if orderStr == "" {
		panic(errors.New("empty sort order in sort descriptor"))
	}

	var sortOrder oddb.SortOrder
	switch orderStr {
	case "asc":
		sortOrder = oddb.Asc
	case "desc":
		sortOrder = oddb.Desc
	default:
		panic(fmt.Errorf("unknown sort order: %v", orderStr))
	}

	sort.KeyPath = keyPath
	sort.Order = sortOrder
}

func sortsFromRaw(rawSorts []interface{}) []oddb.Sort {
	length := len(rawSorts)
	sorts := make([]oddb.Sort, length, length)

	for i := range rawSorts {
		sortSlice, _ := rawSorts[i].([]interface{})
		if len(sortSlice) != 2 {
			panic(fmt.Errorf("got len(sort descriptor) = %v, want 2", len(sortSlice)))
		}
		sortFromRaw(sortSlice, &sorts[i])
	}

	return sorts
}

func queryFromPayload(payload *router.Payload, query *oddb.Query) (err oderr.Error) {
	defer func() {
		// use panic to escape from inner error
		if r := recover(); r != nil {
			if queryErr, ok := r.(error); ok {
				log.WithField("payload", payload).Debugln("failed to construct query")
				err = oderr.NewFmt(oderr.RequestInvalidErr, "failed to construct query: %v", queryErr.Error())
			} else {
				log.WithField("recovered", r).Errorln("panic recovered while constructing query")
				err = oderr.New(oderr.RequestInvalidErr, "error occurred while constructing query")
			}
		}
	}()

	recordType, _ := payload.Data["record_type"].(string)
	if recordType == "" {
		return oderr.New(oderr.RequestInvalidErr, "recordType cannot be empty")
	}
	query.Type = recordType

	if rawSorts, ok := payload.Data["sort"]; ok {
		if rawSorts, ok := rawSorts.([]interface{}); ok {
			query.Sorts = sortsFromRaw(rawSorts)
		} else {
			return oderr.New(oderr.RequestInvalidErr, "order has to be an array")
		}
	}

	return nil
}

type tQuery struct {
	Type string ``
}

/*
RecordQueryHandler is dummy implementation on fetching Records
curl -X POST -H "Content-Type: application/json" \
  -d @- http://localhost:3000/ <<EOF
{
    "action": "record:query",
    "access_token": "validToken",
    "database_id": "private",
    "record_type": "note",
    "sort": [
        [{"$val": "noteOrder", "$type": "desc"}, "asc"]
    ]
}
EOF
*/
func RecordQueryHandler(payload *router.Payload, response *router.Response) {
	db := payload.Database

	query := oddb.Query{}
	if err := queryFromPayload(payload, &query); err != nil {
		response.Err = err
		return
	}

	results, err := db.Query(&query)
	if err != nil {
		response.Err = oderr.New(oderr.UnknownErr, "failed to open database")
		return
	}
	defer results.Close()

	records := []transportRecord{}
	for results.Scan() {
		records = append(records, transportRecord(results.Record()))
	}

	if err != nil {
		response.Err = oderr.New(oderr.UnknownErr, "failed to query records")
		return
	}

	response.Result = records
}

/*
RecordDeleteHandler is dummy implementation on delete Records
curl -X POST -H "Content-Type: application/json" \
  -d @- http://localhost:3000/ <<EOF
{
    "action": "record:delete",
    "access_token": "validToken",
    "database_id": "_private",
    "ids": ["EA6A3E68-90F3-49B5-B470-5FFDB7A0D4E8"]
}
EOF
*/
func RecordDeleteHandler(payload *router.Payload, response *router.Response) {
	if (*recordPayload)(payload).IsWriteAllowed() {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: write is not allowed")
		return
	}

	db := payload.Database

	recordIDs, ok := payload.Data["ids"].([]interface{})
	if !ok {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: expect list of ids")
		return
	}
	results := []idResponseItem{}
	for i := range recordIDs {
		ID, ok := recordIDs[i].(string)
		if !ok {
			results = append(results, idResponseItem{
				ID,
				"_error",
				"ID_FORMAT_ERROR", // FIXME: Dummy
			})
		} else if err := db.Delete(ID); err != nil {
			results = append(results, idResponseItem{
				ID,
				"_error",
				"NOT_FOUND", // FIXME: Dummy
			})
		}
	}
	response.Result = results
	log.Println("RecordDeleteHandler")
	return
}
