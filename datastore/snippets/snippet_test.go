// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package datastore_snippets contains snippet code for the Cloud Datastore API.
// The code is not runnable.
package datastore_snippets

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
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
	// [START incomplete_key]
	taskKey := datastore.IncompleteKey("Task", nil)
	// [END incomplete_key]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey() {
	// [START named_key]
	taskKey := datastore.NameKey("Task", "sampletask", nil)
	// [END named_key]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey_withParent() {
	// [START key_with_parent]
	parentKey := datastore.NameKey("TaskList", "default", nil)
	taskKey := datastore.NameKey("Task", "sampleTask", parentKey)
	// [END key_with_parent]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetNewKey_withMultipleParents() {
	// [START key_with_multilevel_parent]
	userKey := datastore.NameKey("User", "alice", nil)
	parentKey := datastore.NameKey("TaskList", "default", userKey)
	taskKey := datastore.NameKey("Task", "sampleTask", parentKey)
	// [END key_with_multilevel_parent]
	_ = taskKey // Use the task key for datastore operations.
}

func SnippetClient_Put() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START entity_with_parent]
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
	// [END entity_with_parent]
	_ = err // Make sure you check err.
}

func Snippet_properties() {
	// [START properties]
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
	// [END properties]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_sliceProperties() {
	// [START array_value]
	type Task struct {
		Tags          []string
		Collaborators []string
	}
	task := &Task{
		Tags:          []string{"fun", "programming"},
		Collaborators: []string{"alice", "bob"},
	}
	// [END array_value]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_basicEntity() {
	// [START basic_entity]
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
	// [END basic_entity]
	_ = task // Use the task in a datastore Put operation.
}

func SnippetClient_Put_upsert() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	task := &Task{} // Populated with appropriate data.
	key := datastore.IncompleteKey("Task", nil)
	// [START upsert]
	key, err := client.Put(ctx, key, task)
	// [END upsert]
	_ = err // Make sure you check err.
	_ = key // key is the complete key for the newly stored task
}

func SnippetTransaction_insert() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	task := Task{} // Populated with appropriate data.
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	// [START insert]
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
	// [END insert]
	_ = err // Make sure you check err.
}

func SnippetClient_Get() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	// [START lookup]
	var task Task
	err := client.Get(ctx, taskKey, &task)
	// [END lookup]
	_ = err // Make sure you check err.
}

func SnippetTransaction_update() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	taskKey := datastore.NameKey("Task", "sampleTask", nil)
	// [START update]
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	var task Task
	if err := tx.Get(taskKey, &task); err != nil {
		log.Fatalf("tx.Get: %v", err)
	}
	task.Priority = 5
	if _, err := tx.Put(taskKey, task); err != nil {
		log.Fatalf("tx.Put: %v", err)
	}
	if _, err := tx.Commit(); err != nil {
		log.Fatalf("tx.Commit: %v", err)
	}
	// [END update]
}

func SnippetClient_Delete() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	key := datastore.NameKey("Task", "sampletask", nil)
	// [START delete]
	err := client.Delete(ctx, key)
	// [END delete]
	_ = err // Make sure you check err.
}

func SnippetClient_PutMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START batch_upsert]
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
	// [END batch_upsert]
	_ = err  // Make sure you check err.
	_ = keys // keys now has the complete keys for the newly stored tasks.
}

func SnippetClient_GetMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	var taskKeys []*datastore.Key // Populated with incomplete keys.
	// [START batch_lookup]
	var tasks []*Task
	err := client.GetMulti(ctx, taskKeys, &tasks)
	// [END batch_lookup]
	_ = err // Make sure you check err.
}

func SnippetClient_DeleteMulti() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	var taskKeys []*datastore.Key // Populated with incomplete keys.
	// [START batch_delete]
	err := client.DeleteMulti(ctx, taskKeys)
	// [END batch_delete]
	_ = err // Make sure you check err.
}

func SnippetQuery_basic() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START basic_query]
	query := datastore.NewQuery("Task").
		Filter("Done =", false).
		Filter("Priority >=", 4).
		Order("-Priority")
	// [END basic_query]
	// [START run_query]
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
	// [END run_query]
}

func SnippetQuery_propertyFilter() {
	// [START property_filter]
	query := datastore.NewQuery("Task").Filter("Done =", false)
	// [END property_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_compositeFilter() {
	// [START composite_filter]
	query := datastore.NewQuery("Task").Filter("Done =", false).Filter("Priority =", 4)
	// [END composite_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_keyFilter() {
	// [START key_filter]
	key := datastore.NameKey("Task", "someTask", nil)
	query := datastore.NewQuery("Task").Filter("__key__ >", key)
	// [END key_filter]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortAscending() {
	// [START ascending_sort]
	query := datastore.NewQuery("Task").Order("created")
	// [END ascending_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortDescending() {
	// [START descending_sort]
	query := datastore.NewQuery("Task").Order("-created")
	// [END descending_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_sortMulti() {
	// [START multi_sort]
	query := datastore.NewQuery("Task").Order("-priority").Order("created")
	// [END multi_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_kindless() {
	var lastSeenKey *datastore.Key
	// [START kindless_query]
	query := datastore.NewQuery("").Filter("__key__ >", lastSeenKey)
	// [END kindless_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Ancestor() {
	// [START ancestor_query]
	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor)
	// [END ancestor_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Project() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START projection_query]
	query := datastore.NewQuery("Task").Project("Priority", "PercentComplete")
	// [END projection_query]
	// [START run_query_projection]
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
	// [END run_query_projection]
}

func SnippetQuery_KeysOnly() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START keys_only_query]
	query := datastore.NewQuery("Task").KeysOnly()
	// [END keys_only_query]
	// [START run_keys_only_query]
	keys, err := client.GetAll(ctx, query, nil)
	// [END run_keys_only_query]
	_ = err  // Make sure you check err.
	_ = keys // Keys contains keys for all stored tasks.
}

func SnippetQuery_Distinct() {
	// [START distinct_query]
	query := datastore.NewQuery("Task").
		Project("Priority", "PercentComplete").
		Distinct().
		Order("Category").Order("Priority")
	// [END distinct_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_DistinctOn() {
	// [START distinct_on_query]
	query := datastore.NewQuery("Task").
		Project("Priority", "Category").
		DistinctOn("Category").
		Order("Category").Order("Priority")
	// [END distinct_on_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_arrayInequality() {
	// [START array_value_inequality_range]
	query := datastore.NewQuery("Task").
		Filter("Tag >", "learn").
		Filter("Tag <", "math")
	// [END array_value_inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_arrayEquality() {
	// [START array_value_equality]
	query := datastore.NewQuery("Task").
		Filter("Tag =", "fun").
		Filter("Tag =", "programming")
	// [END array_value_equality]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_inequality() {
	// [START inequality_range]
	query := datastore.NewQuery("Task").
		Filter("Created >", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		Filter("Created <", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// [END inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_Filter_invalidInequality() {
	// [START inequality_invalid]
	query := datastore.NewQuery("Task").
		Filter("Created >", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		Filter("Priority >", 3)
	// [END inequality_invalid]
	_ = query // The query is invalid.
}

func SnippetQuery_Filter_mixed() {
	// [START equal_and_inequality_range]
	query := datastore.NewQuery("Task").
		Filter("Priority =", 4).
		Filter("Done =", false).
		Filter("Created >", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)).
		Filter("Created <", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	// [END equal_and_inequality_range]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_inequalitySort() {
	// [START inequality_sort]
	query := datastore.NewQuery("Task").
		Filter("Priority >", 3).
		Order("Priority").
		Order("Created")
	// [END inequality_sort]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_invalidInequalitySortA() {
	// [START inequality_sort_invalid_not_same]
	query := datastore.NewQuery("Task").
		Filter("Priority >", 3).
		Order("Created")
	// [END inequality_sort_invalid_not_same]
	_ = query // The query is invalid.
}

func SnippetQuery_invalidInequalitySortB() {
	// [START inequality_sort_invalid_not_first]
	query := datastore.NewQuery("Task").
		Filter("Priority >", 3).
		Order("Created").
		Order("Priority")
	// [END inequality_sort_invalid_not_first]
	_ = query // The query is invalid.
}

func SnippetQuery_Limit() {
	// [START limit]
	query := datastore.NewQuery("Task").Limit(5)
	// [END limit]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetIterator_Cursor() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	cursorStr := ""
	// [START cursor_paging]
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
	var tasks []Task
	var task Task
	it := client.Run(ctx, query)
	_, err := it.Next(&task)
	for err == nil {
		tasks = append(tasks, task)
		_, err = it.Next(&task)
	}
	if err != iterator.Done {
		log.Fatalf("Failed fetching results: %v", err)
	}

	// Get the cursor for the next page of results.
	nextCursor, err := it.Cursor()
	// [END cursor_paging]
	_ = err        // Check the error.
	_ = nextCursor // Use nextCursor.String as the next page's token.
}

func SnippetQuery_EventualConsistency() {
	// [START eventual_consistent_query]
	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor).EventualConsistency()
	// [END eventual_consistent_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func SnippetQuery_unindexed() {
	// [START unindexed_property_query]
	query := datastore.NewQuery("Tasks").Filter("Description =", "A task description")
	// [END unindexed_property_query]
	_ = query // Use client.Run or client.GetAll to execute the query.
}

func Snippet_explodingProperties() {
	// [START exploding_properties]
	task := &Task{
		Tags:          []string{"fun", "programming", "learn"},
		Collaborators: []string{"alice", "bob", "charlie"},
		Created:       time.Now(),
	}
	// [END exploding_properties]
	_ = task // Use the task in a datastore Put operation.
}

func Snippet_Transaction() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	var to, from *datastore.Key
	// [START transactional_update]
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
	// [END transactional_update]
}

func Snippet_Client_RunInTransaction() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	var to, from *datastore.Key
	// [START transactional_retry]
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
	// [END transactional_retry]
	_ = err // Check error.
}

func SnippetTransaction_getOrCreate() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	key := datastore.NameKey("Task", "sampletask", nil)
	// [START transactional_get_or_create]
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
	// [END transactional_get_or_create]
	_ = err // Check error.
}

func SnippetTransaction_runQuery() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START transactional_single_entity_group_read_only]
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("client.NewTransaction: %v", err)
	}
	defer tx.Rollback() // Transaction only used for read.

	ancestor := datastore.NameKey("TaskList", "default", nil)
	query := datastore.NewQuery("Task").Ancestor(ancestor).Transaction(tx)
	var tasks []Task
	_, err = client.GetAll(ctx, query, &tasks)
	// [END transactional_single_entity_group_read_only]
	_ = err // Check error.
}

func Snippet_metadataNamespaces() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START namespace_run_query]
	const (
		startNamespace = "g"
		endNamespace   = "h"
	)
	query := datastore.NewQuery("__namespace__").
		Filter("__key__ >=", startNamespace).
		Filter("__key__ <", endNamespace).
		KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		log.Fatalf("client.GetAll: %v", err)
	}

	namespaces := make([]string, 0, len(keys))
	for _, k := range keys {
		namespaces = append(namespaces, k.Name)
	}
	// [END namespace_run_query]
}

func Snippet_metadataKinds() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START kind_run_query]
	query := datastore.NewQuery("__kind__").KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		log.Fatalf("client.GetAll: %v", err)
	}

	kinds := make([]string, 0, len(keys))
	for _, k := range keys {
		kinds = append(kinds, k.Name)
	}
	// [END kind_run_query]
}

func Snippet_metadataProperties() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START property_run_query]
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
	// [END property_run_query]
}

func Snippet_metadataPropertiesForKind() {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "my-proj")
	// [START property_by_kind_run_query]
	kindKey := datastore.NameKey("__kind__", "Task", nil)
	query := datastore.NewQuery("__property__").Ancestor(kindKey)

	type Prop struct {
		Repr []string `datastore:"property_representation"`
	}

	var props []Prop
	keys, err := client.GetAll(ctx, query, &props)
	// [END property_by_kind_run_query]
	_ = err  // Check error.
	_ = keys // Use keys to find property names, and props for their representations.
}
