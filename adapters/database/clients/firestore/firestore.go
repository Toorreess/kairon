package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/utils"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Client struct {
	Project string
	Storage *firestore.Client
	Ctx     context.Context
}

func NewFirestoreClient(projectID string) (*Client, error) {
	var client Client
	ctx := context.Background()

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	client.Project = projectID
	client.Storage = fsClient
	client.Ctx = ctx
	return &client, nil
}

func (client *Client) Close() error {
	if client.Storage == nil {
		return fmt.Errorf("No client found")
	}
	return client.Storage.Close()
}

// RunTransaction implements the DBRepository interface for Firestore.
func (c *Client) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	if c.Storage == nil {
		return fmt.Errorf("No client found")
	}

	return c.Storage.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ft := &firestoreTransaction{tx: tx, client: c}
		return f(ft)
	})
}

// NOTE: firestoreTransaction is a concrete type that adapts *firestore.Transaction to
// your DBTransaction interface.  It holds the *firestore.Transaction
// and any other context it needs to satisfy the DBTransaction methods.
type firestoreTransaction struct {
	tx     *firestore.Transaction
	client *Client
}

func (ft *firestoreTransaction) List(index string, entity any, queryOpts infrastructure.QueryOpts) ([]map[string]any, error) {
	if ft.client.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}
	collection := ft.client.Storage.Collection(index)
	q, err := processTxQuery(collection, queryOpts.QueryString, queryOpts.Offset, queryOpts.Limit, queryOpts.OrderBy, queryOpts.Order, queryOpts.RangeBy, queryOpts.RangeSlice, entity, utils.EntityHasDeleted(entity))
	if err != nil {
		return nil, err
	}

	iter := ft.tx.Documents(q)
	defer iter.Stop()
	var results []map[string]any
	for {
		docsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Print(err)
			return nil, err
		}
		result := make(map[string]any)
		if err := docsnap.DataTo(&result); err != nil {
			log.Print(err)
			return nil, err
		}

		result["id"] = docsnap.Ref.ID
		if _, ok := result["modification_date"]; !ok {
			result["modification_date"] = docsnap.UpdateTime
		}
		if _, ok := result["creation_date"]; !ok {
			result["creation_date"] = docsnap.CreateTime
		}
		results = append(results, result)
	}
	return results, err
}

func (ft *firestoreTransaction) Get(index string, id string, entity any) (map[string]any, error) {
	if ft.client.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}

	doc := ft.client.Storage.Collection(index).Doc(id)
	docsnap, err := ft.tx.Get(doc)
	if err != nil {
		return nil, err
	}

	data := make(map[string]any)
	if err := docsnap.DataTo(&data); err != nil {
		return nil, err
	}

	data["id"] = id

	if _, ok := data["modification_date"]; !ok {
		data["modification_date"] = docsnap.UpdateTime
	}
	if _, ok := data["creation_date"]; !ok {
		data["creation_date"] = docsnap.CreateTime
	}
	if utils.EntityHasDeleted(entity) && data["deleted"] == true {
		return nil, fmt.Errorf("Not found")
	}

	return data, err
}

func (ft *firestoreTransaction) Create(index string, entity any) (map[string]any, error) {
	if ft.client.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}

	docRef := ft.client.Storage.Collection(index).NewDoc()
	err := ft.tx.Create(docRef, entity)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &result)
	result["id"] = docRef.ID

	return result, err
}

func (ft *firestoreTransaction) CreateWithID(index, id string, entity any) (map[string]any, error) {
	if ft.client.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}

	docRef := ft.client.Storage.Collection(index).Doc(id)
	err := ft.tx.Create(docRef, entity)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &result)
	result["id"] = docRef.ID

	return result, err
}

func (ft *firestoreTransaction) Update(index string, id string, entity any, changes map[string]any) error {
	if ft.client.Storage == nil {
		return fmt.Errorf("No client found")
	}

	doc := ft.client.Storage.Collection(index).Doc(id)

	updates := []firestore.Update{}
	for k, v := range changes {
		updates = append(updates, firestore.Update{
			Path:  k,
			Value: v,
		})
	}

	err := ft.tx.Update(doc, updates)
	if err != nil {
		return err
	}

	return err
}

func (ft *firestoreTransaction) GetAll(index string, ids []string, entity any) ([]map[string]any, []string, error) {
	if ft.client.Storage == nil {
		return nil, nil, fmt.Errorf("No client found")
	}

	collection := ft.client.Storage.Collection(index)

	var docRefs []*firestore.DocumentRef
	for _, id := range ids {
		docRefs = append(docRefs, collection.Doc(id))
	}

	docsnapList, err := ft.tx.GetAll(docRefs)
	if err != nil {
		return nil, nil, err
	}

	var results []map[string]any
	var errors []string
	for _, docsnap := range docsnapList {
		if docsnap.Exists() {
			result := make(map[string]any)
			if err := docsnap.DataTo(&result); err != nil {
				log.Print(err)
				return nil, nil, err
			}
			result["id"] = docsnap.Ref.ID
			if _, ok := result["modification_date"]; !ok {
				result["modification_date"] = docsnap.UpdateTime
			}
			if _, ok := result["creation_date"]; !ok {
				result["creation_date"] = docsnap.CreateTime
			}
			results = append(results, result)
		} else {
			errors = append(errors, docsnap.Ref.ID)
		}
	}

	return results, errors, nil
}

type DynamicQuery interface {
	Documents(ctx context.Context) *firestore.DocumentIterator
	Where(path string, op string, value any) firestore.Query
	Limit(n int) firestore.Query
	Offset(n int) firestore.Query
	OrderBy(path string, dir firestore.Direction) firestore.Query
}

func convertStringToReflectedType(s string, rt reflect.Type) (any, error) {
	switch rt.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Bool:
		v, err := strconv.ParseBool(s)
		return v, err
	case reflect.Slice:
		return s, nil
	}
	return nil, fmt.Errorf("Not accepted types")
}

func processTxQuery(collection *firestore.CollectionRef, query string, offset, limit int, orderBy, order string, rangeBy string, rangeSlice []any, entity any, hasDeleted bool) (*firestore.Query, error) {
	var q = collection.Query
	andSlice := strings.Split(query, " AND ")
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	e := val.Type()
	typesMap := make(map[string]reflect.Type)
	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Field(i).Tag.Get("firestore")
		if fieldName != "" && fieldName != "-" {
			typesMap[fieldName] = e.Field(i).Type
		}
	}
	for _, s := range andSlice {
		if s != "" {
			replaced := strings.TrimSpace(s)
			r, _ := regexp.Compile("[a-z_]*:")
			r2, _ := regexp.Compile("\".*\"")
			match := r.FindString(replaced)
			var queryField string
			var queryRawValue string
			if match != "" {
				queryField = strings.ReplaceAll(match, ":", "")
				definitionSlice := strings.Split(s, ":")[1]

				match = r2.FindString(definitionSlice)
				if match != "" {
					queryRawValue = strings.TrimSpace(match)
				} else {
					queryRawValue = strings.TrimSpace(fmt.Sprintf("%v", definitionSlice))
				}
			} else {
				log.Print("free search")
				queryField = "name"
				match = r2.FindString(replaced)
				if match != "" {
					queryRawValue = strings.ReplaceAll(match, "\"", "")
				} else {
					queryRawValue = strings.ReplaceAll(replaced, " ", "")
				}
			}
			if tm, ok := typesMap[queryField]; ok {
				if queryField == "terms" {
					queryRawValue = strings.ToUpper(queryRawValue)
				}
				cv, err := convertStringToReflectedType(queryRawValue, tm)
				if err != nil {
					log.Print(err)
					return nil, err
				}

				if tm.Kind() == reflect.Slice {
					q = q.Where(queryField, "array-contains", cv)
				} else {
					q = q.Where(queryField, "==", cv)
				}
			}

		}
	}
	if hasDeleted {
		q = q.Where("deleted", "==", false)
	}
	q = q.Offset(offset).Limit(limit)
	if orderBy != "" {
		if order == "" {
			order = "DESC"
		}
		effectiveOrder := firestore.Desc
		if order == "ASC" {
			effectiveOrder = firestore.Asc
		}
		q = q.OrderBy(orderBy, effectiveOrder)
	}

	if rangeBy != "" && rangeSlice != nil && len(rangeSlice) == 2 {
		if len(rangeSlice) == 2 {
			q = q.Where(rangeBy, ">=", rangeSlice[0])
			q = q.Where(rangeBy, "<", rangeSlice[1])
		} else if len(rangeSlice) == 1 {
			q = q.Where(rangeBy, ">=", rangeSlice[0])
		}
	}
	return &q, nil
}

func processQuery(collection *firestore.CollectionRef, query string, offset, limit int, orderBy, order string, rangeBy string, rangeSlice []any, entity any, hasDeleted bool) (DynamicQuery, error) {
	var q DynamicQuery = collection
	andSlice := strings.Split(query, " AND ")
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	e := val.Type()
	typesMap := make(map[string]reflect.Type)
	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Field(i).Tag.Get("firestore")
		if fieldName != "" && fieldName != "-" {
			typesMap[fieldName] = e.Field(i).Type
		}
	}
	for _, s := range andSlice {
		if s != "" {
			replaced := strings.TrimSpace(s)
			r, _ := regexp.Compile("[a-z_]*:")
			r2, _ := regexp.Compile("\".*\"")
			match := r.FindString(replaced)
			var queryField string
			var queryRawValue string
			if match != "" {
				queryField = strings.ReplaceAll(match, ":", "")
				definitionSlice := strings.Split(s, ":")[1]

				match = r2.FindString(definitionSlice)
				if match != "" {
					queryRawValue = strings.TrimSpace(match)
				} else {
					queryRawValue = strings.TrimSpace(fmt.Sprintf("%v", definitionSlice))
				}
			} else {
				log.Print("free search")
				queryField = "name"
				match = r2.FindString(replaced)
				if match != "" {
					queryRawValue = strings.ReplaceAll(match, "\"", "")
				} else {
					queryRawValue = strings.ReplaceAll(replaced, " ", "")
				}
			}
			if tm, ok := typesMap[queryField]; ok {
				if queryField == "terms" {
					queryRawValue = strings.ToUpper(queryRawValue)
				}
				cv, err := convertStringToReflectedType(queryRawValue, tm)
				if err != nil {
					log.Print(err)
					return nil, err
				}

				if tm.Kind() == reflect.Slice {
					q = q.Where(queryField, "array-contains", cv)
				} else {
					q = q.Where(queryField, "==", cv)
				}
			}

		}
	}
	if hasDeleted {
		q = q.Where("deleted", "==", false)
	}
	q = q.Offset(offset).Limit(limit)
	if orderBy != "" {
		if order == "" {
			order = "DESC"
		}
		effectiveOrder := firestore.Desc
		if order == "ASC" {
			effectiveOrder = firestore.Asc
		}
		q = q.OrderBy(orderBy, effectiveOrder)
	}

	if rangeBy != "" && rangeSlice != nil {
		if len(rangeSlice) == 2 {
			q = q.Where(rangeBy, ">=", rangeSlice[0])
			q = q.Where(rangeBy, "<=", rangeSlice[1])
		} else if len(rangeSlice) == 1 {
			q = q.Where(rangeBy, ">=", rangeSlice[0])
		}
	}
	return q, nil
}

func (c Client) Get(index, id string, hasDeleted bool) (map[string]any, error) {
	if c.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}

	collection := c.Storage.Collection(index)
	doc := collection.Doc(id)
	docsnap, err := doc.Get(c.Ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	if err := docsnap.DataTo(&result); err != nil {
		return nil, err
	}
	result["id"] = id

	if _, ok := result["modification_date"]; !ok {
		result["modification_date"] = docsnap.UpdateTime
	}
	if _, ok := result["creation_date"]; !ok {
		result["creation_date"] = docsnap.CreateTime
	}
	if hasDeleted && result["deleted"] == true {
		return nil, fmt.Errorf("Not found")
	}
	return result, nil
}

func (c Client) Create(index string, entity any) (map[string]any, error) {
	if c.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	docRef, wr, err := collection.Add(c.Ctx, entity)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &result)
	result["id"] = docRef.ID
	if _, ok := result["modification_date"]; !ok {
		result["modification_date"] = wr.UpdateTime
	}
	if _, ok := result["creation_date"]; !ok {
		result["creation_date"] = wr.UpdateTime
	}
	return result, nil
}

func (c Client) CreateWithID(index string, id string, entity any) (map[string]any, error) {
	if c.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	docRef := collection.Doc(id)
	wr, err := docRef.Set(c.Ctx, entity)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &result)
	result["id"] = docRef.ID
	if _, ok := result["modification_date"]; !ok {
		result["modification_date"] = wr.UpdateTime
	}
	if _, ok := result["creation_date"]; !ok {
		result["creation_date"] = wr.UpdateTime
	}
	return result, nil
}

func (c Client) Update(index, id string, hasDeleted bool, updates map[string]any) (map[string]any, error) {
	if c.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	doc := collection.Doc(id)
	var fsUpdates []firestore.Update
	for k, v := range updates {
		fsUpdates = append(fsUpdates, firestore.Update{Path: k, Value: v})
	}
	_, err := doc.Update(c.Ctx, fsUpdates)
	if err != nil {
		return nil, err
	}
	return c.Get(index, id, hasDeleted)

}

func (c Client) TUpdate(tx *firestore.Transaction, index, id string, hasDeleted bool, updates map[string]any) error {
	if c.Storage == nil {
		return fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	doc := collection.Doc(id)
	var fsUpdates []firestore.Update
	for k, v := range updates {
		fsUpdates = append(fsUpdates, firestore.Update{Path: k, Value: v})
	}
	err := tx.Update(doc, fsUpdates)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) Delete(index, id string, hasDeleted bool) error {
	if c.Storage == nil {
		return fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	doc := collection.Doc(id)
	if hasDeleted {
		if _, err := doc.Update(c.Ctx, []firestore.Update{{Path: "deleted", Value: true}}); err != nil {
			return err
		}
	} else {
		if _, err := doc.Delete(c.Ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c Client) List(index string, entity any, hasDeleted bool, queryOpts infrastructure.QueryOpts) ([]map[string]any, error) {
	if c.Storage == nil {
		return nil, fmt.Errorf("No client found")
	}
	collection := c.Storage.Collection(index)
	q, err := processQuery(collection, queryOpts.QueryString, queryOpts.Offset, queryOpts.Limit, queryOpts.OrderBy, queryOpts.Order, queryOpts.RangeBy, queryOpts.RangeSlice, entity, hasDeleted)
	if err != nil {
		return nil, err
	}
	// TODO Falta insertar queries
	iter := q.Documents(c.Ctx)
	defer iter.Stop()
	var results []map[string]any
	for {
		docsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Print(err)
			return nil, err
		}
		result := make(map[string]any)
		if err := docsnap.DataTo(&result); err != nil {
			log.Print(err)
			return nil, err
		}
		result["id"] = docsnap.Ref.ID
		if _, ok := result["modification_date"]; !ok {
			result["modification_date"] = docsnap.UpdateTime
		}
		if _, ok := result["creation_date"]; !ok {
			result["creation_date"] = docsnap.CreateTime
		}
		results = append(results, result)
	}
	return results, nil
}

func (c Client) GetAll(index string, ids []string, entity any, hasDeleted bool) ([]map[string]any, []string, error) {
	if c.Storage == nil {
		return nil, nil, fmt.Errorf("No client found")
	}

	collection := c.Storage.Collection(index)

	var docRefs []*firestore.DocumentRef
	for _, id := range ids {
		docRefs = append(docRefs, collection.Doc(id))
	}

	docsnapList, err := c.Storage.GetAll(c.Ctx, docRefs)
	if err != nil {
		return nil, nil, err
	}

	var results []map[string]any
	var errors []string
	for _, docsnap := range docsnapList {
		if docsnap.Exists() {
			result := make(map[string]any)
			if err := docsnap.DataTo(&result); err != nil {
				log.Print(err)
				return nil, nil, err
			}
			result["id"] = docsnap.Ref.ID
			if _, ok := result["modification_date"]; !ok {
				result["modification_date"] = docsnap.UpdateTime
			}
			if _, ok := result["creation_date"]; !ok {
				result["creation_date"] = docsnap.CreateTime
			}
			results = append(results, result)
		} else {
			errors = append(errors, docsnap.Ref.ID)
		}
	}

	return results, errors, nil
}
