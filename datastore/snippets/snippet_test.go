// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package datastore_snippets contains snippet code for the Cloud Datastore API.
// The code is not runnable.
package datastore_snippets

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

type Task struct {
	Category        string
	Done            bool
	Priority        int
	Description     string `datastore:",noindex"`
	PercentComplete float64
	Created         time.Time
	Tags            []string
	Collaborators   []string
}

func SnippetNewIncompleteKey() {
	// [START datastore_incomplete_key]
	// A complete key is assigned to the entity when it is Put.
	taskKey := datastore.IncompleteKey("Task", nil)
	// [END datastore_incomplete_key]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey() {
	// [START datastore_named_key]
	taskKey := datastore.NameKey("Task", "sampletask", nil)
	// [END datastore_named_key]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey_withParent() {
	// [START datastore_key_with_parent]
	parentKey := datastore.NameKey("TaskList", "default", nil)
	taskKey := datastore.NameKey("Task", "sampleTask", parentKey)
	// [END datastore_key_with_parent]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey_withMultipleParents() {
	// [START datastore_key_with_multilevel_parent]
	userKey := datastore.NameKey("User", "alice", nil)
	parentKey := datastore.NameKey("TaskList", "default", userKey)
	taskKey := datastore.NameKey("Task", "sampleTask", parentKey)
	// [END datastore_key_with_multilevel_parent]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetClient_Put() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_entity_with_parent]
	parentKey := datastore.NameKey("TaskList", "default", nil)
	key := datastore.IncompleteKey("Task", parentKey)

	task := Task{
		Category:    "Personal",
		Done:        false,
		Priority:    4,
		Description: "Learn Cloud Datastore",
	}

	// A complete key is assigned to the entity when it is Put.
	var err error
	key, err = client.Put(ctx, key, &task)
	// [END datastore_entity_with_parent]
	_ = err // Make sure you check err.
}

func Snippet_properties() {
	// [START datastore_properties]
	type Task struct {
		Category        string
		Done            bool
		Priority        int
		Description     string `datastore:",noindex"`
		PercentComplete float64
		Created         time.Time
	}
	task := &Task{
		Category:        "Personal",
		Done:            false,
		Priority:        4,
		Description:     "Learn Cloud Datastore",
		PercentComplete: 10.0,
		Created:         time.Now(),
	}
	// [END datastore_properties]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_sliceProperties() {
	// [START datastore_array_value]
	type Task struct {
		Tags          []string
		Collaborators []string
	}
	task := &Task{
		Tags:          []string{"fun", "programming"},
		Collaborators: []string{"alice", "bob"},
	}
	// [END datastore_array_value]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_basicEntity() {
	// [START datastore_basic_entity]
	type Task struct {
		Category        string
		Done            bool
		Priority        float64
		Description     string `datastore:",noindex"`
		PercentComplete float64
		Created         time.Time
	}
	task := &Task{
		Category:        "Personal",
		Done:            false,
		Priority:        4,
		Description:     "Learn Cloud Datastore",
		PercentComplete: 10.0,
		Created:         time.Now(),
	}
	// [END datastore_basic_entity]
	_ = task // Use the task in a datastore Put operation.
}

func SnippetClient_Put_upsert() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	task := &Task{} // Populated with appropriate data.
	// [START datastore_upsert]
	key := datastore.IncompleteKey("Task", nil)
	key, err := client.Put(ctx, key, task)
	// [END datastore_upsert]
	_ = err // Make sure you check err.
	_ = key // key is the complete key for the newly stored task
}

func SnippetTransaction_insert() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	task := Task{} // Populated with appropriate data.
	// [START datastore_insert]
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// We first check that there is no entity stored with the given key.
		var empty Task
		if err := tx.Get(taskKey, &empty); err != datastore.ErrNoSuchEntity {
			return err
		}
		// If there was no matching entity, store it now.
		_, err := tx.Put(taskKey, &task)
		return err
	})
	// [END datastore_insert]
	_ = err // Make sure you check err.
}

func SnippetClient_Get() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_lookup]
	var task Task
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	err := client.Get(ctx, taskKey, &task)
	// [END datastore_lookup]
	_ = err // Make sure you check err.
}

func SnippetTransaction_update() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_update]
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	var task Task
	if err := tx.Get(taskKey, &task); err != nil {
		log.Fatalf("tx.Get: %v", err)
	}
	task.Priority = 5
	if _, err := tx.Put(taskKey, &task); err != nil {
		log.Fatalf("tx.Put: %v", err)
	}
	if _, err := tx.Commit(); err != nil {
		log.Fatalf("tx.Commit: %v", err)
	}
	// [END datastore_update]
}

func SnippetClient_Delete() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_delete]
	key := datastore.NameKey("Task", "sampletask", nil)
	err := client.Delete(ctx, key)
	// [END datastore_delete]
	_ = err // Make sure you check err.
}

func SnippetClient_PutMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_batch_upsert]
	tasks := []*Task{
		{
			Category:    "Personal",
			Done:        false,
			Priority:    4,
			Description: "Learn Cloud Datastore",
		},
		{
			Category:    "Personal",
			Done:        false,
			Priority:    5,
			Description: "Integrate Cloud Datastore",
		},
	}
	keys := []*datastore.Key{
		datastore.IncompleteKey("Task", nil),
		datastore.IncompleteKey("Task", nil),
	}

	keys, err := client.PutMulti(ctx, keys, tasks)
	// [END datastore_batch_upsert]
	_ = err  // Make sure you check err.
	_ = keys // keys now has the complete keys for the newly stored tasks.
}

func SnippetClient_GetMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_batch_lookup]
	var taskKeys []*datastore.Key // Populated with incomplete keys.
	tasks := make([]*Task, len(taskKeys))
	err := client.GetMulti(ctx, taskKeys, &tasks)
	// [END datastore_batch_lookup]
	_ = err // Make sure you check err.
}

func SnippetClient_DeleteMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	var taskKeys []*datastore.Key // Populated with incomplete keys.
	// [START datastore_batch_delete]
	err := client.DeleteMulti(ctx, taskKeys)
	// [END datastore_batch_delete]
	_ = err // Make sure you check err.
}

func SnippetQuery_basic() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_basic_query]
	query := datastore.NewQuery("Task").
		FilterField("Done", "=", false).
		FilterField("Priority", ">=", 4).
		Order("-Priority")
	// [END datastore_basic_query]
	// [START datastore_run_query]
	it := client.Run(ctx, query)
	for {
		var task Task
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next task: %v", err)
		}
		fmt.Printf("Task %q, Priority %d\n", task.Description, task.Priority)
	}
	// [END datastore_run_query]
}

func SnippetQuery_propertyFilter() {
	// [START datastore_property_filter]
	query := datastore.NewQuery("Task").FilterField("Done", "=", false)
	// [END datastore_property_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_compositeFilter() {
	// [START datastore_composite_filter]
	query := datastore.NewQuery("Task").
		FilterField("Done", "=", false).
		FilterField("Priority", "=", 4)
	// [END datastore_composite_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_keyFilter() {
	// [START datastore_key_filter]
	key := datastore.NameKey("Task", "someTask", nil)
	query := datastore.NewQuery("Task").FilterField("__key__", ">", key)
	// [END datastore_key_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortAscending() {
	// [START datastore_ascending_sort]
	query := datastore.NewQuery("Task").Order("created")
	// [END datastore_ascending_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortDescending() {
	// [START datastore_descending_sort]
	query := datastore.NewQuery("Task").Order("-created")
	// [END datastore_descending_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortMulti() {
	// [START datastore_multi_sort]
	query := datastore.NewQuery("Task").Order("-priority").Order("created")
	// [END datastore_multi_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_kindless() {
	var lastSeenKey *datastore.Key
	// [START datastore_kindless_query]
	query := datastore.NewQuery("").FilterField("__key__", ">", lastSeenKey)
	// [END datastore_kindless_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Ancestor() {
	// [START datastore_ancestor_query]
	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor)
	// [END datastore_ancestor_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Project() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_projection_query]
	query := datastore.NewQuery("Task").Project("Priority", "PercentComplete")
	// [END datastore_projection_query]
	// [START datastore_run_query_projection]
	var priorities []int
	var percents []float64
	it := client.Run(ctx, query)
	for {
		var task Task
		if _, err := it.Next(&task); err == iterator.Done {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		priorities = append(priorities, task.Priority)
		percents = append(percents, task.PercentComplete)
	}
	// [END datastore_run_query_projection]
}

func SnippetQuery_KeysOnly() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_keys_only_query]
	query := datastore.NewQuery("Task").KeysOnly()
	// [END datastore_keys_only_query]

	keys, err := client.GetAll(ctx, query, nil)
	_ = err  // Make sure you check err.
	_ = keys // Keys contains keys for all stored tasks.
}

func SnippetQuery_DistinctOn() {
	// [START datastore_distinct_on_query]
	query := datastore.NewQuery("Task").
		Project("Priority", "Category").
		DistinctOn("Category").
		Order("Category").Order("Priority")
	// [END datastore_distinct_on_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_arrayInequality() {
	// [START datastore_array_value_inequality_range]
	query := datastore.NewQuery("Task").
		FilterField("Tag", ">", "learn").
		FilterField("Tag", "<", "math")
	// [END datastore_array_value_inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_arrayEquality() {
	// [START datastore_array_value_equality]
	query := datastore.NewQuery("Task").
		FilterField("Tag", "=", "fun").
		FilterField("Tag", "=", "programming")
	// [END datastore_array_value_equality]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_inequality() {
	// [START datastore_inequality_range]
	query := datastore.NewQuery("Task").
		FilterField("Created", ">", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		FilterField("Created", "<", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// [END datastore_inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_invalidInequality() {
	// [START datastore_inequality_invalid]
	query := datastore.NewQuery("Task").
		FilterField("Created", ">", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		FilterField("Priority", ">", 3)
	// [END datastore_inequality_invalid]
	_ = query // The query is invalid.
}

func SnippetQuery_Filter_mixed() {
	// [START datastore_equal_and_inequality_range]
	query := datastore.NewQuery("Task").
		FilterField("Priority", "=", 4).
		FilterField("Done", "=", false).
		FilterField("Created", ">", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		FilterField("Created", "<", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// [END datastore_equal_and_inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_inequalitySort() {
	// [START datastore_inequality_sort]
	query := datastore.NewQuery("Task").
		FilterField("Priority", ">", 3).
		Order("Priority").
		Order("Created")
	// [END datastore_inequality_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_invalidInequalitySortA() {
	// [START datastore_inequality_sort_invalid_not_same]
	query := datastore.NewQuery("Task").
		FilterField("Priority", ">", 3).
		Order("Created")
	// [END datastore_inequality_sort_invalid_not_same]
	_ = query // The query is invalid.
}

func SnippetQuery_invalidInequalitySortB() {
	// [START datastore_inequality_sort_invalid_not_first]
	query := datastore.NewQuery("Task").
		FilterField("Priority", ">", 3).
		Order("Created").
		Order("Priority")
	// [END datastore_inequality_sort_invalid_not_first]
	_ = query // The query is invalid.
}

func SnippetQuery_Limit() {
	// [START datastore_limit]
	query := datastore.NewQuery("Task").Limit(5)
	// [END datastore_limit]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetIterator_Cursor() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_cursor_paging]
	// cursorStr is a cursor to start querying at.
	cursorStr := ""

	const pageSize = 5
	query := datastore.NewQuery("Tasks").Limit(pageSize)
	if cursorStr != "" {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			log.Fatalf("Bad cursor %q: %v", cursorStr, err)
		}
		query = query.Start(cursor)
	}

	// Read the tasks.
	it := client.Run(ctx, query)
	var tasks []Task
	for {
		var task Task
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed fetching results: %v", err)
		}
		tasks = append(tasks, task)
	}

	// Get the cursor for the next page of results.
	// nextCursor.String can be used as the next page's token.
	nextCursor, err := it.Cursor()
	// [END datastore_cursor_paging]
	_ = err        // Check the error.
	_ = nextCursor // Use nextCursor.String as the next page's token.
}

func SnippetQuery_EventualConsistency() {
	// [START datastore_eventual_consistent_query]
	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor).EventualConsistency()
	// [END datastore_eventual_consistent_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_unindexed() {
	// [START datastore_unindexed_property_query]
	query := datastore.NewQuery("Tasks").
		FilterField("Description", "=", "A task description")
	// [END datastore_unindexed_property_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func Snippet_explodingProperties() {
	// [START datastore_exploding_properties]
	task := &Task{
		Tags:          []string{"fun", "programming", "learn"},
		Collaborators: []string{"alice", "bob", "charlie"},
		Created:       time.Now(),
	}
	// [END datastore_exploding_properties]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_Transaction() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	var to, from *datastore.Key
	// [START datastore_transactional_update]
	type BankAccount struct {
		Balance int
	}

	const amount = 50
	keys := []*datastore.Key{to, from}
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	accs := make([]BankAccount, 2)
	if err := tx.GetMulti(keys, accs); err != nil {
		tx.Rollback()
		log.Fatalf("tx.GetMulti: %v", err)
	}
	accs[0].Balance += amount
	accs[1].Balance -= amount
	if _, err := tx.PutMulti(keys, accs); err != nil {
		tx.Rollback()
		log.Fatalf("tx.PutMulti: %v", err)
	}
	if _, err = tx.Commit(); err != nil {
		log.Fatalf("tx.Commit: %v", err)
	}
	// [END datastore_transactional_update]
}

func Snippet_Client_RunInTransaction() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	var to, from *datastore.Key
	// [START datastore_transactional_retry]
	type BankAccount struct {
		Balance int
	}

	const amount = 50
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		keys := []*datastore.Key{to, from}
		accs := make([]BankAccount, 2)
		if err := tx.GetMulti(keys, accs); err != nil {
			return err
		}
		accs[0].Balance += amount
		accs[1].Balance -= amount
		_, err := tx.PutMulti(keys, accs)
		return err
	})
	// [END datastore_transactional_retry]
	_ = err // Check error.
}

func SnippetTransaction_getOrCreate() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	key := datastore.NameKey("Task", "sampletask", nil)
	// [START datastore_transactional_get_or_create]
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		if err := tx.Get(key, &task); err != datastore.ErrNoSuchEntity {
			return err
		}
		_, err := tx.Put(key, &Task{
			Category:    "Personal",
			Done:        false,
			Priority:    4,
			Description: "Learn Cloud Datastore",
		})
		return err
	})
	// [END datastore_transactional_get_or_create]
	_ = err // Check error.
}

func SnippetTransaction_runQuery() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_transactional_single_entity_group_read_only]
	tx, err := client.NewTransaction(ctx, datastore.ReadOnly)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	defer tx.Rollback() // Transaction only used for read.

	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor).Transaction(tx)
	var tasks []Task
	_, err = client.GetAll(ctx, query, &tasks)
	// [END datastore_transactional_single_entity_group_read_only]
	_ = err // Check error.
}

// [START datastore_namespace_run_query]

func metadataNamespaces(w io.Writer, projectID string) error {
	// projectID := "my-project"

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("datastore.NewClient: %w", err)
	}
	defer client.Close()

	start := datastore.NameKey("__namespace__", "g", nil)
	end := datastore.NameKey("__namespace__", "h", nil)
	query := datastore.NewQuery("__namespace__").
		FilterField("__key__", ">=", start).
		FilterField("__key__", "<", end).
		KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("client.GetAll: %w", err)
	}

	fmt.Fprintln(w, "Namespaces:")
	for _, k := range keys {
		fmt.Fprintf(w, "\t%v", k.Name)
	}
	return nil
}

// [END datastore_namespace_run_query]

func TestMetadaNamespaces(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := metadataNamespaces(buf, tc.ProjectID); err != nil {
		t.Errorf("metadataNamespaces got err: %v, want no error", err)
	}
	if got, want := buf.String(), "Namespaces"; !strings.Contains(got, want) {
		t.Errorf("metadataNamespaces got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}
}

func Snippet_metadataKinds() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_kind_run_query]
	query := datastore.NewQuery("__kind__").KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		log.Fatalf("client.GetAll: %v", err)
	}

	kinds := make([]string, 0, len(keys))
	for _, k := range keys {
		kinds = append(kinds, k.Name)
	}
	// [END datastore_kind_run_query]
}

func Snippet_metadataProperties() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_property_run_query]
	query := datastore.NewQuery("__property__").KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		log.Fatalf("client.GetAll: %v", err)
	}

	props := make(map[string][]string) // Map from kind to slice of properties.
	for _, k := range keys {
		prop := k.Name
		kind := k.Parent.Name
		props[kind] = append(props[kind], prop)
	}
	// [END datastore_property_run_query]
}

func Snippet_metadataPropertiesForKind() {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()
	// [START datastore_property_by_kind_run_query]
	kindKey := datastore.NameKey("__kind__", "Task", nil)
	query := datastore.NewQuery("__property__").Ancestor(kindKey)

	type Prop struct {
		Repr []string `datastore:"property_representation"`
	}

	var props []Prop
	keys, err := client.GetAll(ctx, query, &props)
	// [END datastore_property_by_kind_run_query]
	_ = err  // Check error.
	_ = keys // Use keys to find property names, and props for their representations.
}

func SnippetQuery_RunQueryWithExplain(w io.Writer) {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()

	// [START datastore_query_explain_entity]
	// Build the query
	query := datastore.NewQuery("Task")

	// Set the explain options to get back *only* the plan summary
	it := client.RunWithOptions(ctx, query, datastore.ExplainOptions{})
	_, err := it.Next(nil)

	// Get the explain metrics
	explainMetrics := it.ExplainMetrics

	planSummary := explainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END datastore_query_explain_entity]
	_ = err // Check non-nil errors other than Iterator.Done
}

func SnippetQuery_RunQueryWithExplainAnalyze(w io.Writer) {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()

	// [START datastore_query_explain_analyze_entity]
	// Build the query
	query := datastore.NewQuery("Task")

	// Set explain options with analzye = true to get back the query stats, plan info, and query
	// results
	it := client.RunWithOptions(ctx, query, datastore.ExplainOptions{Analyze: true})

	// Get the query results
	fmt.Fprintln(w, "----- Query Results -----")
	for {
		var task Task
		_, err := it.Next(&task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Error fetching next task: %v", err)
			return
		}
		fmt.Fprintf(w, "Task %q, Priority %d\n", task.Description, task.Priority)
	}

	// Get plan summary
	planSummary := it.ExplainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	// Get the execution stats
	executionStats := it.ExplainMetrics.ExecutionStats
	fmt.Fprintln(w, "----- Execution Stats -----")
	fmt.Fprintf(w, "%+v\n", executionStats)
	fmt.Fprintln(w, "----- Debug Stats -----")
	for k, v := range *executionStats.DebugStats {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END datastore_query_explain_analyze_entity]
}

func SnippetQuery_RunAggregationQueryWithExplain(w io.Writer) {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()

	// [START datastore_query_explain_aggregation]
	// Build the query
	query := datastore.NewQuery("Task")

	// Set the explain options to get back *only* the plan summary
	ar, err := client.RunAggregationQueryWithOptions(ctx, query.NewAggregationQuery().WithCount("count"), datastore.ExplainOptions{})

	// Get the explain metrics
	explainMetrics := ar.ExplainMetrics

	planSummary := explainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END datastore_query_explain_aggregation]

	_ = err // Check non-nil errors
}

func SnippetQuery_RunAggregationQueryWithExplainAnalyze(w io.Writer) {
	ctx := context.Background()
	client, _ := datastore.NewClientWithDatabase(ctx, "my-proj", "my-database-id")
	defer client.Close()

	// [START datastore_query_explain_analyze_aggregation]
	// Build the query
	query := datastore.NewQuery("Task")

	// Set explain options with analzye = true to get back the query stats, plan info, and query
	// results
	countAlias := "count"
	ar, err := client.RunAggregationQueryWithOptions(ctx,
		query.NewAggregationQuery().WithCount(countAlias), datastore.ExplainOptions{Analyze: true})

	// Get the query results
	fmt.Fprintln(w, "----- Query Results -----")
	result := ar.Result[countAlias]
	fmt.Fprintf(w, "Count %v\n", result)

	// Get plan summary
	planSummary := ar.ExplainMetrics.PlanSummary
	fmt.Fprintln(w, "----- Indexes Used -----")
	for k, v := range planSummary.IndexesUsed {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}

	// Get the execution stats
	executionStats := ar.ExplainMetrics.ExecutionStats
	fmt.Fprintln(w, "----- Execution Stats -----")
	fmt.Fprintf(w, "%+v\n", executionStats)
	fmt.Fprintln(w, "----- Debug Stats -----")
	for k, v := range *executionStats.DebugStats {
		fmt.Fprintf(w, "%+v: %+v\n", k, v)
	}
	// [END datastore_query_explain_analyze_aggregation]

	_ = err // Check non-nil errors
}
