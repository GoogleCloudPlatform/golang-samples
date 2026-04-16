// Copyright 2026 Google LLC
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

package firestore

import (
	"bytes"
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestPipelineSnippets(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	databaseID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_ENTERPRISE_DATABASE")
	if databaseID == "" {
		t.Skip("Skipping firestore enterprise test. Set GOLANG_SAMPLES_FIRESTORE_ENTERPRISE_DATABASE.")
	}

	ctx := context.Background()

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	buf := new(bytes.Buffer)

	t.Run("pipelineConcepts", func(t *testing.T) {
		buf.Reset()
		if err := pipelineConcepts(buf, client); err != nil {
			t.Errorf("pipelineConcepts failed: %v", err)
		}
	})

	t.Run("basicRead", func(t *testing.T) {
		buf.Reset()
		if err := basicRead(buf, client); err != nil {
			t.Errorf("basicRead failed: %v", err)
		}
	})

	t.Run("fieldVsConstants", func(t *testing.T) {
		buf.Reset()
		if err := fieldVsConstants(buf, client); err != nil {
			t.Errorf("fieldVsConstants failed: %v", err)
		}
	})

	t.Run("inputStages", func(t *testing.T) {
		buf.Reset()
		if err := inputStages(buf, client); err != nil {
			t.Errorf("inputStages failed: %v", err)
		}
	})

	t.Run("wherePipeline", func(t *testing.T) {
		buf.Reset()
		if err := wherePipeline(buf, client); err != nil {
			t.Errorf("wherePipeline failed: %v", err)
		}
	})

	t.Run("aggregateGroups", func(t *testing.T) {
		buf.Reset()
		if err := aggregateGroups(buf, client); err != nil {
			t.Errorf("aggregateGroups failed: %v", err)
		}
	})

	t.Run("aggregateDistinct", func(t *testing.T) {
		buf.Reset()
		if err := aggregateDistinct(buf, client); err != nil {
			t.Errorf("aggregateDistinct failed: %v", err)
		}
	})

	t.Run("sort", func(t *testing.T) {
		buf.Reset()
		if err := sort(buf, client); err != nil {
			t.Errorf("sort failed: %v", err)
		}
	})

	t.Run("sortComparison", func(t *testing.T) {
		buf.Reset()
		if err := sortComparison(buf, client); err != nil {
			t.Errorf("sortComparison failed: %v", err)
		}
	})

	t.Run("functionsExample", func(t *testing.T) {
		buf.Reset()
		if err := functionsExample(buf, client); err != nil {
			t.Errorf("functionsExample failed: %v", err)
		}
	})

	t.Run("creatingIndexes", func(t *testing.T) {
		buf.Reset()
		if err := creatingIndexes(buf, client); err != nil {
			t.Errorf("creatingIndexes failed: %v", err)
		}
	})

	t.Run("sparseIndexes", func(t *testing.T) {
		buf.Reset()
		if err := sparseIndexes(buf, client); err != nil {
			t.Errorf("sparseIndexes failed: %v", err)
		}
	})

	t.Run("sparseIndexes2", func(t *testing.T) {
		buf.Reset()
		if err := sparseIndexes2(buf, client); err != nil {
			t.Errorf("sparseIndexes2 failed: %v", err)
		}
	})

	t.Run("coveredQuery", func(t *testing.T) {
		buf.Reset()
		if err := coveredQuery(buf, client); err != nil {
			t.Errorf("coveredQuery failed: %v", err)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		buf.Reset()
		if err := pagination(buf, client); err != nil {
			t.Errorf("pagination failed: %v", err)
		}
	})

	t.Run("collectionStage", func(t *testing.T) {
		buf.Reset()
		if err := collectionStage(buf, client); err != nil {
			t.Errorf("collectionStage failed: %v", err)
		}
	})

	t.Run("collectionGroupStage", func(t *testing.T) {
		buf.Reset()
		if err := collectionGroupStage(buf, client); err != nil {
			t.Errorf("collectionGroupStage failed: %v", err)
		}
	})

	t.Run("databaseStage", func(t *testing.T) {
		buf.Reset()
		if err := databaseStage(buf, client); err != nil {
			t.Errorf("databaseStage failed: %v", err)
		}
	})

	t.Run("documentsStage", func(t *testing.T) {
		buf.Reset()
		if err := documentsStage(buf, client); err != nil {
			t.Errorf("documentsStage failed: %v", err)
		}
	})

	t.Run("replaceWithStage", func(t *testing.T) {
		buf.Reset()
		if err := replaceWithStage(buf, client); err != nil {
			t.Errorf("replaceWithStage failed: %v", err)
		}
	})

	t.Run("sampleStage", func(t *testing.T) {
		buf.Reset()
		if err := sampleStage(buf, client); err != nil {
			t.Errorf("sampleStage failed: %v", err)
		}
	})

	t.Run("samplePercent", func(t *testing.T) {
		buf.Reset()
		if err := samplePercent(buf, client); err != nil {
			t.Errorf("samplePercent failed: %v", err)
		}
	})

	t.Run("unionStage", func(t *testing.T) {
		buf.Reset()
		if err := unionStage(buf, client); err != nil {
			t.Errorf("unionStage failed: %v", err)
		}
	})

	t.Run("unionStageStable", func(t *testing.T) {
		buf.Reset()
		if err := unionStageStable(buf, client); err != nil {
			t.Errorf("unionStageStable failed: %v", err)
		}
	})

	t.Run("unnestStage", func(t *testing.T) {
		buf.Reset()
		if err := unnestStage(buf, client); err != nil {
			t.Errorf("unnestStage failed: %v", err)
		}
	})

	t.Run("unnestStageEmptyOrNonArray", func(t *testing.T) {
		buf.Reset()
		if err := unnestStageEmptyOrNonArray(buf, client); err != nil {
			t.Errorf("unnestStageEmptyOrNonArray failed: %v", err)
		}
	})

	t.Run("countFunction", func(t *testing.T) {
		buf.Reset()
		if err := countFunction(buf, client); err != nil {
			t.Errorf("countFunction failed: %v", err)
		}
	})

	t.Run("countIfFunction", func(t *testing.T) {
		buf.Reset()
		if err := countIfFunction(buf, client); err != nil {
			t.Errorf("countIfFunction failed: %v", err)
		}
	})

	t.Run("countDistinctFunction", func(t *testing.T) {
		buf.Reset()
		if err := countDistinctFunction(buf, client); err != nil {
			t.Errorf("countDistinctFunction failed: %v", err)
		}
	})

	t.Run("sumFunction", func(t *testing.T) {
		buf.Reset()
		if err := sumFunction(buf, client); err != nil {
			t.Errorf("sumFunction failed: %v", err)
		}
	})

	t.Run("avgFunction", func(t *testing.T) {
		buf.Reset()
		if err := avgFunction(buf, client); err != nil {
			t.Errorf("avgFunction failed: %v", err)
		}
	})

	t.Run("minFunction", func(t *testing.T) {
		buf.Reset()
		if err := minFunction(buf, client); err != nil {
			t.Errorf("minFunction failed: %v", err)
		}
	})

	t.Run("maxFunction", func(t *testing.T) {
		buf.Reset()
		if err := maxFunction(buf, client); err != nil {
			t.Errorf("maxFunction failed: %v", err)
		}
	})

	t.Run("addFunction", func(t *testing.T) {
		buf.Reset()
		if err := addFunction(buf, client); err != nil {
			t.Errorf("addFunction failed: %v", err)
		}
	})

	t.Run("subtractFunction", func(t *testing.T) {
		buf.Reset()
		if err := subtractFunction(buf, client); err != nil {
			t.Errorf("subtractFunction failed: %v", err)
		}
	})

	t.Run("multiplyFunction", func(t *testing.T) {
		buf.Reset()
		if err := multiplyFunction(buf, client); err != nil {
			t.Errorf("multiplyFunction failed: %v", err)
		}
	})

	t.Run("divideFunction", func(t *testing.T) {
		buf.Reset()
		if err := divideFunction(buf, client); err != nil {
			t.Errorf("divideFunction failed: %v", err)
		}
	})

	t.Run("modFunction", func(t *testing.T) {
		buf.Reset()
		if err := modFunction(buf, client); err != nil {
			t.Errorf("modFunction failed: %v", err)
		}
	})

	t.Run("ceilFunction", func(t *testing.T) {
		buf.Reset()
		if err := ceilFunction(buf, client); err != nil {
			t.Errorf("ceilFunction failed: %v", err)
		}
	})

	t.Run("floorFunction", func(t *testing.T) {
		buf.Reset()
		if err := floorFunction(buf, client); err != nil {
			t.Errorf("floorFunction failed: %v", err)
		}
	})

	t.Run("roundFunction", func(t *testing.T) {
		buf.Reset()
		if err := roundFunction(buf, client); err != nil {
			t.Errorf("roundFunction failed: %v", err)
		}
	})

	t.Run("powFunction", func(t *testing.T) {
		buf.Reset()
		if err := powFunction(buf, client); err != nil {
			t.Errorf("powFunction failed: %v", err)
		}
	})

	t.Run("sqrtFunction", func(t *testing.T) {
		buf.Reset()
		if err := sqrtFunction(buf, client); err != nil {
			t.Errorf("sqrtFunction failed: %v", err)
		}
	})

	t.Run("expFunction", func(t *testing.T) {
		buf.Reset()
		if err := expFunction(buf, client); err != nil {
			t.Errorf("expFunction failed: %v", err)
		}
	})

	t.Run("lnFunction", func(t *testing.T) {
		buf.Reset()
		if err := lnFunction(buf, client); err != nil {
			t.Errorf("lnFunction failed: %v", err)
		}
	})

	t.Run("logFunction", func(t *testing.T) {
		buf.Reset()
		if err := logFunction(buf, client); err != nil {
			t.Errorf("logFunction failed: %v", err)
		}
	})

	t.Run("arrayConcatFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayConcatFunction(buf, client); err != nil {
			t.Errorf("arrayConcatFunction failed: %v", err)
		}
	})

	t.Run("arrayContainsFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayContainsFunction(buf, client); err != nil {
			t.Errorf("arrayContainsFunction failed: %v", err)
		}
	})

	t.Run("arrayContainsAllFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayContainsAllFunction(buf, client); err != nil {
			t.Errorf("arrayContainsAllFunction failed: %v", err)
		}
	})

	t.Run("arrayContainsAnyFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayContainsAnyFunction(buf, client); err != nil {
			t.Errorf("arrayContainsAnyFunction failed: %v", err)
		}
	})

	t.Run("arrayLengthFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayLengthFunction(buf, client); err != nil {
			t.Errorf("arrayLengthFunction failed: %v", err)
		}
	})

	t.Run("arrayReverseFunction", func(t *testing.T) {
		buf.Reset()
		if err := arrayReverseFunction(buf, client); err != nil {
			t.Errorf("arrayReverseFunction failed: %v", err)
		}
	})

	t.Run("equalFunction", func(t *testing.T) {
		buf.Reset()
		if err := equalFunction(buf, client); err != nil {
			t.Errorf("equalFunction failed: %v", err)
		}
	})

	t.Run("greaterThanFunction", func(t *testing.T) {
		buf.Reset()
		if err := greaterThanFunction(buf, client); err != nil {
			t.Errorf("greaterThanFunction failed: %v", err)
		}
	})

	t.Run("greaterThanOrEqualToFunction", func(t *testing.T) {
		buf.Reset()
		if err := greaterThanOrEqualToFunction(buf, client); err != nil {
			t.Errorf("greaterThanOrEqualToFunction failed: %v", err)
		}
	})

	t.Run("lessThanFunction", func(t *testing.T) {
		buf.Reset()
		if err := lessThanFunction(buf, client); err != nil {
			t.Errorf("lessThanFunction failed: %v", err)
		}
	})

	t.Run("lessThanOrEqualToFunction", func(t *testing.T) {
		buf.Reset()
		if err := lessThanOrEqualToFunction(buf, client); err != nil {
			t.Errorf("lessThanOrEqualToFunction failed: %v", err)
		}
	})

	t.Run("notEqualFunction", func(t *testing.T) {
		buf.Reset()
		if err := notEqualFunction(buf, client); err != nil {
			t.Errorf("notEqualFunction failed: %v", err)
		}
	})

	t.Run("existsFunction", func(t *testing.T) {
		buf.Reset()
		if err := existsFunction(buf, client); err != nil {
			t.Errorf("existsFunction failed: %v", err)
		}
	})

	t.Run("andFunction", func(t *testing.T) {
		buf.Reset()
		if err := andFunction(buf, client); err != nil {
			t.Errorf("andFunction failed: %v", err)
		}
	})

	t.Run("orFunction", func(t *testing.T) {
		buf.Reset()
		if err := orFunction(buf, client); err != nil {
			t.Errorf("orFunction failed: %v", err)
		}
	})

	t.Run("xorFunction", func(t *testing.T) {
		buf.Reset()
		if err := xorFunction(buf, client); err != nil {
			t.Errorf("xorFunction failed: %v", err)
		}
	})

	t.Run("notFunction", func(t *testing.T) {
		buf.Reset()
		if err := notFunction(buf, client); err != nil {
			t.Errorf("notFunction failed: %v", err)
		}
	})

	t.Run("condFunction", func(t *testing.T) {
		buf.Reset()
		if err := condFunction(buf, client); err != nil {
			t.Errorf("condFunction failed: %v", err)
		}
	})

	t.Run("equalAnyFunction", func(t *testing.T) {
		buf.Reset()
		if err := equalAnyFunction(buf, client); err != nil {
			t.Errorf("equalAnyFunction failed: %v", err)
		}
	})

	t.Run("notEqualAnyFunction", func(t *testing.T) {
		buf.Reset()
		if err := notEqualAnyFunction(buf, client); err != nil {
			t.Errorf("notEqualAnyFunction failed: %v", err)
		}
	})

	t.Run("maxLogicalFunction", func(t *testing.T) {
		buf.Reset()
		if err := maxLogicalFunction(buf, client); err != nil {
			t.Errorf("maxLogicalFunction failed: %v", err)
		}
	})

	t.Run("minLogicalFunction", func(t *testing.T) {
		buf.Reset()
		if err := minLogicalFunction(buf, client); err != nil {
			t.Errorf("minLogicalFunction failed: %v", err)
		}
	})

	t.Run("mapGetFunction", func(t *testing.T) {
		buf.Reset()
		if err := mapGetFunction(buf, client); err != nil {
			t.Errorf("mapGetFunction failed: %v", err)
		}
	})

	t.Run("mapSetFunction", func(t *testing.T) {
		buf.Reset()
		if err := mapSetFunction(buf, client); err != nil {
			t.Errorf("mapSetFunction failed: %v", err)
		}
	})

	t.Run("mapKeysFunction", func(t *testing.T) {
		buf.Reset()
		if err := mapKeysFunction(buf, client); err != nil {
			t.Errorf("mapKeysFunction failed: %v", err)
		}
	})

	t.Run("mapValuesFunction", func(t *testing.T) {
		buf.Reset()
		if err := mapValuesFunction(buf, client); err != nil {
			t.Errorf("mapValuesFunction failed: %v", err)
		}
	})

	t.Run("mapEntriesFunction", func(t *testing.T) {
		buf.Reset()
		if err := mapEntriesFunction(buf, client); err != nil {
			t.Errorf("mapEntriesFunction failed: %v", err)
		}
	})

	t.Run("byteLengthFunction", func(t *testing.T) {
		buf.Reset()
		if err := byteLengthFunction(buf, client); err != nil {
			t.Errorf("byteLengthFunction failed: %v", err)
		}
	})

	t.Run("charLengthFunction", func(t *testing.T) {
		buf.Reset()
		if err := charLengthFunction(buf, client); err != nil {
			t.Errorf("charLengthFunction failed: %v", err)
		}
	})

	t.Run("startsWithFunction", func(t *testing.T) {
		buf.Reset()
		if err := startsWithFunction(buf, client); err != nil {
			t.Errorf("startsWithFunction failed: %v", err)
		}
	})

	t.Run("likeFunction", func(t *testing.T) {
		buf.Reset()
		if err := likeFunction(buf, client); err != nil {
			t.Errorf("likeFunction failed: %v", err)
		}
	})

	t.Run("regexContainsFunction", func(t *testing.T) {
		buf.Reset()
		if err := regexContainsFunction(buf, client); err != nil {
			t.Errorf("regexContainsFunction failed: %v", err)
		}
	})

	t.Run("regexFindFunction", func(t *testing.T) {
		buf.Reset()
		if err := regexFindFunction(buf, client); err != nil {
			t.Errorf("regexFindFunction failed: %v", err)
		}
	})

	t.Run("regexFindAllFunction", func(t *testing.T) {
		buf.Reset()
		if err := regexFindAllFunction(buf, client); err != nil {
			t.Errorf("regexFindAllFunction failed: %v", err)
		}
	})

	t.Run("regexMatchFunction", func(t *testing.T) {
		buf.Reset()
		if err := regexMatchFunction(buf, client); err != nil {
			t.Errorf("regexMatchFunction failed: %v", err)
		}
	})

	t.Run("strConcatFunction", func(t *testing.T) {
		buf.Reset()
		if err := strConcatFunction(buf, client); err != nil {
			t.Errorf("strConcatFunction failed: %v", err)
		}
	})

	t.Run("strContainsFunction", func(t *testing.T) {
		buf.Reset()
		if err := strContainsFunction(buf, client); err != nil {
			t.Errorf("strContainsFunction failed: %v", err)
		}
	})

	t.Run("toUpperFunction", func(t *testing.T) {
		buf.Reset()
		if err := toUpperFunction(buf, client); err != nil {
			t.Errorf("toUpperFunction failed: %v", err)
		}
	})

	t.Run("toLowerFunction", func(t *testing.T) {
		buf.Reset()
		if err := toLowerFunction(buf, client); err != nil {
			t.Errorf("toLowerFunction failed: %v", err)
		}
	})

	t.Run("substrFunction", func(t *testing.T) {
		buf.Reset()
		if err := substrFunction(buf, client); err != nil {
			t.Errorf("substrFunction failed: %v", err)
		}
	})

	t.Run("strReverseFunction", func(t *testing.T) {
		buf.Reset()
		if err := strReverseFunction(buf, client); err != nil {
			t.Errorf("strReverseFunction failed: %v", err)
		}
	})

	t.Run("strTrimFunction", func(t *testing.T) {
		buf.Reset()
		if err := strTrimFunction(buf, client); err != nil {
			t.Errorf("strTrimFunction failed: %v", err)
		}
	})

	t.Run("strLTrimFunction", func(t *testing.T) {
		buf.Reset()
		if err := strLTrimFunction(buf, client); err != nil {
			t.Errorf("strLTrimFunction failed: %v", err)
		}
	})

	t.Run("strRTrimFunction", func(t *testing.T) {
		buf.Reset()
		if err := strRTrimFunction(buf, client); err != nil {
			t.Errorf("strRTrimFunction failed: %v", err)
		}
	})

	t.Run("strRepeatFunction", func(t *testing.T) {
		buf.Reset()
		if err := strRepeatFunction(buf, client); err != nil {
			t.Errorf("strRepeatFunction failed: %v", err)
		}
	})

	t.Run("strReplaceAllFunction", func(t *testing.T) {
		buf.Reset()
		if err := strReplaceAllFunction(buf, client); err != nil {
			t.Errorf("strReplaceAllFunction failed: %v", err)
		}
	})

	t.Run("strReplaceOneFunction", func(t *testing.T) {
		buf.Reset()
		if err := strReplaceOneFunction(buf, client); err != nil {
			t.Errorf("strReplaceOneFunction failed: %v", err)
		}
	})

	t.Run("strIndexOfFunction", func(t *testing.T) {
		buf.Reset()
		if err := strIndexOfFunction(buf, client); err != nil {
			t.Errorf("strIndexOfFunction failed: %v", err)
		}
	})

	t.Run("unixMicrosToTimestampFunction", func(t *testing.T) {
		buf.Reset()
		if err := unixMicrosToTimestampFunction(buf, client); err != nil {
			t.Errorf("unixMicrosToTimestampFunction failed: %v", err)
		}
	})

	t.Run("unixMillisToTimestampFunction", func(t *testing.T) {
		buf.Reset()
		if err := unixMillisToTimestampFunction(buf, client); err != nil {
			t.Errorf("unixMillisToTimestampFunction failed: %v", err)
		}
	})

	t.Run("unixSecondsToTimestampFunction", func(t *testing.T) {
		buf.Reset()
		if err := unixSecondsToTimestampFunction(buf, client); err != nil {
			t.Errorf("unixSecondsToTimestampFunction failed: %v", err)
		}
	})

	t.Run("timestampAddFunction", func(t *testing.T) {
		buf.Reset()
		if err := timestampAddFunction(buf, client); err != nil {
			t.Errorf("timestampAddFunction failed: %v", err)
		}
	})

	t.Run("timestampSubFunction", func(t *testing.T) {
		buf.Reset()
		if err := timestampSubFunction(buf, client); err != nil {
			t.Errorf("timestampSubFunction failed: %v", err)
		}
	})

	t.Run("timestampToUnixMicrosFunction", func(t *testing.T) {
		buf.Reset()
		if err := timestampToUnixMicrosFunction(buf, client); err != nil {
			t.Errorf("timestampToUnixMicrosFunction failed: %v", err)
		}
	})

	t.Run("timestampToUnixMillisFunction", func(t *testing.T) {
		buf.Reset()
		if err := timestampToUnixMillisFunction(buf, client); err != nil {
			t.Errorf("timestampToUnixMillisFunction failed: %v", err)
		}
	})

	t.Run("timestampToUnixSecondsFunction", func(t *testing.T) {
		buf.Reset()
		if err := timestampToUnixSecondsFunction(buf, client); err != nil {
			t.Errorf("timestampToUnixSecondsFunction failed: %v", err)
		}
	})

	t.Run("cosineDistanceFunction", func(t *testing.T) {
		buf.Reset()
		if err := cosineDistanceFunction(buf, client); err != nil {
			t.Errorf("cosineDistanceFunction failed: %v", err)
		}
	})

	t.Run("dotProductFunction", func(t *testing.T) {
		buf.Reset()
		if err := dotProductFunction(buf, client); err != nil {
			t.Errorf("dotProductFunction failed: %v", err)
		}
	})

	t.Run("euclideanDistanceFunction", func(t *testing.T) {
		buf.Reset()
		if err := euclideanDistanceFunction(buf, client); err != nil {
			t.Errorf("euclideanDistanceFunction failed: %v", err)
		}
	})

	t.Run("vectorLengthFunction", func(t *testing.T) {
		buf.Reset()
		if err := vectorLengthFunction(buf, client); err != nil {
			t.Errorf("vectorLengthFunction failed: %v", err)
		}
	})

	t.Run("stagesExpressionsExample", func(t *testing.T) {
		buf.Reset()
		if err := stagesExpressionsExample(buf, client); err != nil {
			t.Errorf("stagesExpressionsExample failed: %v", err)
		}
	})

	t.Run("createWhereData", func(t *testing.T) {
		buf.Reset()
		if err := createWhereData(buf, client); err != nil {
			t.Errorf("createWhereData failed: %v", err)
		}
	})

	t.Run("whereEqualityExample", func(t *testing.T) {
		buf.Reset()
		if err := whereEqualityExample(buf, client); err != nil {
			t.Errorf("whereEqualityExample failed: %v", err)
		}
	})

	t.Run("whereMultipleStagesExample", func(t *testing.T) {
		buf.Reset()
		if err := whereMultipleStagesExample(buf, client); err != nil {
			t.Errorf("whereMultipleStagesExample failed: %v", err)
		}
	})

	t.Run("whereComplexExample", func(t *testing.T) {
		buf.Reset()
		if err := whereComplexExample(buf, client); err != nil {
			t.Errorf("whereComplexExample failed: %v", err)
		}
	})

	t.Run("whereStageOrderExample", func(t *testing.T) {
		buf.Reset()
		if err := whereStageOrderExample(buf, client); err != nil {
			t.Errorf("whereStageOrderExample failed: %v", err)
		}
	})

	t.Run("unnestSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestSyntaxExample(buf, client); err != nil {
			t.Errorf("unnestSyntaxExample failed: %v", err)
		}
	})

	t.Run("unnestAliasIndexDataExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestAliasIndexDataExample(buf, client); err != nil {
			t.Errorf("unnestAliasIndexDataExample failed: %v", err)
		}
	})

	t.Run("unnestAliasIndexExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestAliasIndexExample(buf, client); err != nil {
			t.Errorf("unnestAliasIndexExample failed: %v", err)
		}
	})

	t.Run("unnestNonArrayDataExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestNonArrayDataExample(buf, client); err != nil {
			t.Errorf("unnestNonArrayDataExample failed: %v", err)
		}
	})

	t.Run("unnestNonArrayExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestNonArrayExample(buf, client); err != nil {
			t.Errorf("unnestNonArrayExample failed: %v", err)
		}
	})

	t.Run("unnestEmptyArrayDataExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestEmptyArrayDataExample(buf, client); err != nil {
			t.Errorf("unnestEmptyArrayDataExample failed: %v", err)
		}
	})

	t.Run("unnestEmptyArrayExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestEmptyArrayExample(buf, client); err != nil {
			t.Errorf("unnestEmptyArrayExample failed: %v", err)
		}
	})

	t.Run("unnestPreserveEmptyArrayExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestPreserveEmptyArrayExample(buf, client); err != nil {
			t.Errorf("unnestPreserveEmptyArrayExample failed: %v", err)
		}
	})

	t.Run("unnestNestedDataExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestNestedDataExample(buf, client); err != nil {
			t.Errorf("unnestNestedDataExample failed: %v", err)
		}
	})

	t.Run("unnestNestedExample", func(t *testing.T) {
		buf.Reset()
		if err := unnestNestedExample(buf, client); err != nil {
			t.Errorf("unnestNestedExample failed: %v", err)
		}
	})

	t.Run("sampleSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := sampleSyntaxExample(buf, client); err != nil {
			t.Errorf("sampleSyntaxExample failed: %v", err)
		}
	})

	t.Run("sampleDocumentsDataExample", func(t *testing.T) {
		buf.Reset()
		if err := sampleDocumentsDataExample(buf, client); err != nil {
			t.Errorf("sampleDocumentsDataExample failed: %v", err)
		}
	})

	t.Run("sampleDocumentsExample", func(t *testing.T) {
		buf.Reset()
		if err := sampleDocumentsExample(buf, client); err != nil {
			t.Errorf("sampleDocumentsExample failed: %v", err)
		}
	})

	t.Run("sampleAllDocumentsExample", func(t *testing.T) {
		buf.Reset()
		if err := sampleAllDocumentsExample(buf, client); err != nil {
			t.Errorf("sampleAllDocumentsExample failed: %v", err)
		}
	})

	t.Run("samplePercentageDataExample", func(t *testing.T) {
		buf.Reset()
		if err := samplePercentageDataExample(buf, client); err != nil {
			t.Errorf("samplePercentageDataExample failed: %v", err)
		}
	})

	t.Run("samplePercentageExample", func(t *testing.T) {
		buf.Reset()
		if err := samplePercentageExample(buf, client); err != nil {
			t.Errorf("samplePercentageExample failed: %v", err)
		}
	})

	t.Run("sortSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := sortSyntaxExample(buf, client); err != nil {
			t.Errorf("sortSyntaxExample failed: %v", err)
		}
	})

	t.Run("sortSyntaxExample2", func(t *testing.T) {
		buf.Reset()
		if err := sortSyntaxExample2(buf, client); err != nil {
			t.Errorf("sortSyntaxExample2 failed: %v", err)
		}
	})

	t.Run("sortDocumentIDExample", func(t *testing.T) {
		buf.Reset()
		if err := sortDocumentIDExample(buf, client); err != nil {
			t.Errorf("sortDocumentIDExample failed: %v", err)
		}
	})

	t.Run("selectSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := selectSyntaxExample(buf, client); err != nil {
			t.Errorf("selectSyntaxExample failed: %v", err)
		}
	})

	t.Run("selectPositionDataExample", func(t *testing.T) {
		buf.Reset()
		if err := selectPositionDataExample(buf, client); err != nil {
			t.Errorf("selectPositionDataExample failed: %v", err)
		}
	})

	t.Run("selectPositionExample", func(t *testing.T) {
		buf.Reset()
		if err := selectPositionExample(buf, client); err != nil {
			t.Errorf("selectPositionExample failed: %v", err)
		}
	})

	t.Run("selectBadPositionExample", func(t *testing.T) {
		buf.Reset()
		if err := selectBadPositionExample(buf, client); err != nil {
			t.Errorf("selectBadPositionExample failed: %v", err)
		}
	})

	t.Run("selectNestedDataExample", func(t *testing.T) {
		buf.Reset()
		if err := selectNestedDataExample(buf, client); err != nil {
			t.Errorf("selectNestedDataExample failed: %v", err)
		}
	})

	t.Run("selectNestedExample", func(t *testing.T) {
		buf.Reset()
		if err := selectNestedExample(buf, client); err != nil {
			t.Errorf("selectNestedExample failed: %v", err)
		}
	})

	t.Run("removeFieldsSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := removeFieldsSyntaxExample(buf, client); err != nil {
			t.Errorf("removeFieldsSyntaxExample failed: %v", err)
		}
	})

	t.Run("removeFieldsNestedDataExample", func(t *testing.T) {
		buf.Reset()
		if err := removeFieldsNestedDataExample(buf, client); err != nil {
			t.Errorf("removeFieldsNestedDataExample failed: %v", err)
		}
	})

	t.Run("removeFieldsNestedExample", func(t *testing.T) {
		buf.Reset()
		if err := removeFieldsNestedExample(buf, client); err != nil {
			t.Errorf("removeFieldsNestedExample failed: %v", err)
		}
	})

	t.Run("limitSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := limitSyntaxExample(buf, client); err != nil {
			t.Errorf("limitSyntaxExample failed: %v", err)
		}
	})

	t.Run("findNearestSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := findNearestSyntaxExample(buf, client); err != nil {
			t.Errorf("findNearestSyntaxExample failed: %v", err)
		}
	})

	t.Run("findNearestLimitExample", func(t *testing.T) {
		buf.Reset()
		if err := findNearestLimitExample(buf, client); err != nil {
			t.Errorf("findNearestLimitExample failed: %v", err)
		}
	})

	t.Run("findNearestDistanceDataExample", func(t *testing.T) {
		buf.Reset()
		if err := findNearestDistanceDataExample(buf, client); err != nil {
			t.Errorf("findNearestDistanceDataExample failed: %v", err)
		}
	})

	t.Run("findNearestDistanceExample", func(t *testing.T) {
		buf.Reset()
		if err := findNearestDistanceExample(buf, client); err != nil {
			t.Errorf("findNearestDistanceExample failed: %v", err)
		}
	})

	t.Run("offsetSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := offsetSyntaxExample(buf, client); err != nil {
			t.Errorf("offsetSyntaxExample failed: %v", err)
		}
	})

	t.Run("addFieldsSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := addFieldsSyntaxExample(buf, client); err != nil {
			t.Errorf("addFieldsSyntaxExample failed: %v", err)
		}
	})

	t.Run("addFieldsOverlapExample", func(t *testing.T) {
		buf.Reset()
		if err := addFieldsOverlapExample(buf, client); err != nil {
			t.Errorf("addFieldsOverlapExample failed: %v", err)
		}
	})

	t.Run("addFieldsNestingExample", func(t *testing.T) {
		buf.Reset()
		if err := addFieldsNestingExample(buf, client); err != nil {
			t.Errorf("addFieldsNestingExample failed: %v", err)
		}
	})

	t.Run("collectionInputSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := collectionInputSyntaxExample(buf, client); err != nil {
			t.Errorf("collectionInputSyntaxExample failed: %v", err)
		}
	})

	t.Run("collectionInputExampleData", func(t *testing.T) {
		buf.Reset()
		if err := collectionInputExampleData(buf, client); err != nil {
			t.Errorf("collectionInputExampleData failed: %v", err)
		}
	})

	t.Run("collectionInputExample", func(t *testing.T) {
		buf.Reset()
		if err := collectionInputExample(buf, client); err != nil {
			t.Errorf("collectionInputExample failed: %v", err)
		}
	})

	t.Run("subcollectionInputExampleData", func(t *testing.T) {
		buf.Reset()
		if err := subcollectionInputExampleData(buf, client); err != nil {
			t.Errorf("subcollectionInputExampleData failed: %v", err)
		}
	})

	t.Run("subcollectionInputExample", func(t *testing.T) {
		buf.Reset()
		if err := subcollectionInputExample(buf, client); err != nil {
			t.Errorf("subcollectionInputExample failed: %v", err)
		}
	})

	t.Run("collectionGroupInputSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := collectionGroupInputSyntaxExample(buf, client); err != nil {
			t.Errorf("collectionGroupInputSyntaxExample failed: %v", err)
		}
	})

	t.Run("collectionGroupInputExampleData", func(t *testing.T) {
		buf.Reset()
		if err := collectionGroupInputExampleData(buf, client); err != nil {
			t.Errorf("collectionGroupInputExampleData failed: %v", err)
		}
	})

	t.Run("collectionGroupInputExample", func(t *testing.T) {
		buf.Reset()
		if err := collectionGroupInputExample(buf, client); err != nil {
			t.Errorf("collectionGroupInputExample failed: %v", err)
		}
	})

	t.Run("databaseInputSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := databaseInputSyntaxExample(buf, client); err != nil {
			t.Errorf("databaseInputSyntaxExample failed: %v", err)
		}
	})

	t.Run("databaseInputSyntaxExampleData", func(t *testing.T) {
		buf.Reset()
		if err := databaseInputSyntaxExampleData(buf, client); err != nil {
			t.Errorf("databaseInputSyntaxExampleData failed: %v", err)
		}
	})

	t.Run("databaseInputExample", func(t *testing.T) {
		buf.Reset()
		if err := databaseInputExample(buf, client); err != nil {
			t.Errorf("databaseInputExample failed: %v", err)
		}
	})

	t.Run("documentInputSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := documentInputSyntaxExample(buf, client); err != nil {
			t.Errorf("documentInputSyntaxExample failed: %v", err)
		}
	})

	t.Run("documentInputExampleData", func(t *testing.T) {
		buf.Reset()
		if err := documentInputExampleData(buf, client); err != nil {
			t.Errorf("documentInputExampleData failed: %v", err)
		}
	})

	t.Run("documentInputExample", func(t *testing.T) {
		buf.Reset()
		if err := documentInputExample(buf, client); err != nil {
			t.Errorf("documentInputExample failed: %v", err)
		}
	})

	t.Run("unionSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := unionSyntaxExample(buf, client); err != nil {
			t.Errorf("unionSyntaxExample failed: %v", err)
		}
	})

	t.Run("aggregateSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := aggregateSyntaxExample(buf, client); err != nil {
			t.Errorf("aggregateSyntaxExample failed: %v", err)
		}
	})

	t.Run("aggregateGroupSyntax", func(t *testing.T) {
		buf.Reset()
		if err := aggregateGroupSyntax(buf, client); err != nil {
			t.Errorf("aggregateGroupSyntax failed: %v", err)
		}
	})

	t.Run("aggregateExampleData", func(t *testing.T) {
		buf.Reset()
		if err := aggregateExampleData(buf, client); err != nil {
			t.Errorf("aggregateExampleData failed: %v", err)
		}
	})

	t.Run("aggregateWithoutGroupExample", func(t *testing.T) {
		buf.Reset()
		if err := aggregateWithoutGroupExample(buf, client); err != nil {
			t.Errorf("aggregateWithoutGroupExample failed: %v", err)
		}
	})

	t.Run("aggregateGroupExample", func(t *testing.T) {
		buf.Reset()
		if err := aggregateGroupExample(buf, client); err != nil {
			t.Errorf("aggregateGroupExample failed: %v", err)
		}
	})

	t.Run("aggregateGroupComplexExample", func(t *testing.T) {
		buf.Reset()
		if err := aggregateGroupComplexExample(buf, client); err != nil {
			t.Errorf("aggregateGroupComplexExample failed: %v", err)
		}
	})

	t.Run("distinctSyntaxExample", func(t *testing.T) {
		buf.Reset()
		if err := distinctSyntaxExample(buf, client); err != nil {
			t.Errorf("distinctSyntaxExample failed: %v", err)
		}
	})

	t.Run("distinctExampleData", func(t *testing.T) {
		buf.Reset()
		if err := distinctExampleData(buf, client); err != nil {
			t.Errorf("distinctExampleData failed: %v", err)
		}
	})

	t.Run("distinctExample", func(t *testing.T) {
		buf.Reset()
		if err := distinctExample(buf, client); err != nil {
			t.Errorf("distinctExample failed: %v", err)
		}
	})

	t.Run("distinctExpressionsExample", func(t *testing.T) {
		buf.Reset()
		if err := distinctExpressionsExample(buf, client); err != nil {
			t.Errorf("distinctExpressionsExample failed: %v", err)
		}
	})

	t.Run("endsWithFunction", func(t *testing.T) {
		buf.Reset()
		if err := endsWithFunction(buf, client); err != nil {
			t.Errorf("endsWithFunction failed: %v", err)
		}
	})

	t.Run("whereHavingExample", func(t *testing.T) {
		buf.Reset()
		if err := whereHavingExample(buf, client); err != nil {
			t.Errorf("whereHavingExample failed: %v", err)
		}
	})

	t.Run("searchBasic", func(t *testing.T) {
		t.Skip("Requires a pre-created search index")
		buf.Reset()
		if err := searchBasic(buf, client); err != nil {
			t.Errorf("searchBasic failed: %v", err)
		}
	})

	t.Run("searchExact", func(t *testing.T) {
		t.Skip("Requires a pre-created search index")
		buf.Reset()
		if err := searchExact(buf, client); err != nil {
			t.Errorf("searchExact failed: %v", err)
		}
	})

	t.Run("searchTwoTerms", func(t *testing.T) {
		t.Skip("Requires a pre-created search index")
		buf.Reset()
		if err := searchTwoTerms(buf, client); err != nil {
			t.Errorf("searchTwoTerms failed: %v", err)
		}
	})

	t.Run("searchExcludeTerm", func(t *testing.T) {
		t.Skip("Requires a pre-created search index")
		buf.Reset()
		if err := searchExcludeTerm(buf, client); err != nil {
			t.Errorf("searchExcludeTerm failed: %v", err)
		}
	})

	t.Run("searchSpecialFields", func(t *testing.T) {
		t.Skip("Requires a pre-created search index")
		buf.Reset()
		if err := searchSpecialFields(buf, client); err != nil {
			t.Errorf("searchSpecialFields failed: %v", err)
		}
	})

	t.Run("defineExample", func(t *testing.T) {
		buf.Reset()
		if err := defineExample(buf, client); err != nil {
			t.Errorf("defineExample failed: %v", err)
		}
	})

	t.Run("toArrayExpression", func(t *testing.T) {
		buf.Reset()
		if err := toArrayExpression(buf, client); err != nil {
			t.Errorf("toArrayExpression failed: %v", err)
		}
	})

	t.Run("toScalarExpression", func(t *testing.T) {
		buf.Reset()
		if err := toScalarExpression(buf, client); err != nil {
			t.Errorf("toScalarExpression failed: %v", err)
		}
	})

	t.Run("forceIndexExamples", func(t *testing.T) {
		t.Skip("Requires a pre-created specific index")
		buf.Reset()
		if err := forceIndexExamples(buf, client); err != nil {
			t.Errorf("forceIndexExamples failed: %v", err)
		}
	})

	t.Run("pipelineUpdate", func(t *testing.T) {
		buf.Reset()
		if err := pipelineUpdate(buf, client); err != nil {
			t.Errorf("pipelineUpdate failed: %v", err)
		}
	})

	t.Run("pipelineDelete", func(t *testing.T) {
		buf.Reset()
		if err := pipelineDelete(buf, client); err != nil {
			t.Errorf("pipelineDelete failed: %v", err)
		}
	})

}
