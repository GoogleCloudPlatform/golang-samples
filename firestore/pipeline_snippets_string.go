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
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func byteLengthFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_byte_length]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.ByteLength(firestore.FieldOf("title")).As("titleByteLength"),
		)).
		Execute(ctx)
	// [END firestore_byte_length]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func charLengthFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_char_length]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.CharLength(firestore.FieldOf("title")).As("titleCharLength"),
		)).
		Execute(ctx)
	// [END firestore_char_length]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func startsWithFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_starts_with]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StartsWith(firestore.FieldOf("title"), "The").As("needsSpecialAlphabeticalSort"),
		)).
		Execute(ctx)
	// [END firestore_starts_with]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func likeFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_like]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Like(firestore.FieldOf("genre"), "%Fiction").As("anyFiction"),
		)).
		Execute(ctx)
	// [END firestore_like]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func regexContainsFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_regex_contains]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.RegexContains(firestore.FieldOf("title"), "Firestore (Enterprise|Standard)").As("isFirestoreRelated"),
		)).
		Execute(ctx)
	// [END firestore_regex_contains]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func regexFindFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_regex_find]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.RegexFind(firestore.FieldOf("email"), "@[A-Za-z0-9.-]+").As("domain"),
		)).
		Execute(ctx)
	// [END firestore_regex_find]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func regexFindAllFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_regex_find_all]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.RegexFindAll(firestore.FieldOf("comment"), "@[A-Za-z0-9_]+").As("mentions"),
		)).
		Execute(ctx)
	// [END firestore_regex_find_all]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func regexMatchFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_regex_match]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.RegexMatch(firestore.FieldOf("title"), "Firestore (Enterprise|Standard)").As("isFirestoreExactly"),
		)).
		Execute(ctx)
	// [END firestore_regex_match]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strConcatFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_str_concat]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StringConcat(firestore.FieldOf("title"), " by ", firestore.FieldOf("author")).As("fullyQualifiedTitle"),
		)).
		Execute(ctx)
	// [END firestore_str_concat]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strContainsFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_string_contains]
	snapshot := client.Pipeline().
		Collection("articles").
		Select(firestore.Fields(
			firestore.StringContains(firestore.FieldOf("body"), "Firestore").As("isFirestoreRelated"),
		)).
		Execute(ctx)
	// [END firestore_string_contains]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func toUpperFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_to_upper]
	snapshot := client.Pipeline().
		Collection("authors").
		Select(firestore.Fields(
			firestore.ToUpper(firestore.FieldOf("name")).As("uppercaseName"),
		)).
		Execute(ctx)
	// [END firestore_to_upper]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func toLowerFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_to_lower]
	snapshot := client.Pipeline().
		Collection("authors").
		Select(firestore.Fields(
			firestore.Equal(firestore.ToLower(firestore.FieldOf("genre")), "fantasy").As("isFantasy"),
		)).
		Execute(ctx)
	// [END firestore_to_lower]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func substrFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_substr_function]
	snapshot := client.Pipeline().
		Collection("books").
		Where(firestore.StartsWith(firestore.FieldOf("title"), "The ")).
		Select(firestore.Fields(
			firestore.Substring(firestore.FieldOf("title"), firestore.ConstantOf(4), firestore.FieldOf("title").CharLength()).As("titleWithoutLeadingThe"),
		)).
		Execute(ctx)
	// [END firestore_substr_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strReverseFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_str_reverse]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Reverse(firestore.FieldOf("name")).As("reversedName"),
		)).
		Execute(ctx)
	// [END firestore_str_reverse]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strTrimFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_trim_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Trim(firestore.FieldOf("name")).As("whitespaceTrimmedName"),
		)).
		Execute(ctx)
	// [END firestore_trim_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strLTrimFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_ltrim_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.LTrim(firestore.FieldOf("name")).As("ltrimmedName"),
		)).
		Execute(ctx)
	// [END firestore_ltrim_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strRTrimFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_rtrim_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.RTrim(firestore.FieldOf("name")).As("rtrimmedName"),
		)).
		Execute(ctx)
	// [END firestore_rtrim_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strRepeatFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_string_repeat_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StringRepeat(firestore.FieldOf("title"), 2).As("repeatedTitle"),
		)).
		Execute(ctx)
	// [END firestore_string_repeat_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strReplaceAllFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_string_replace_all_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StringReplaceAll(firestore.FieldOf("title"), "The", "A").As("replacedTitle"),
		)).
		Execute(ctx)
	// [END firestore_string_replace_all_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strReplaceOneFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_string_replace_one_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StringReplaceOne(firestore.FieldOf("title"), "The", "A").As("replacedTitle"),
		)).
		Execute(ctx)
	// [END firestore_string_replace_one_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func strIndexOfFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_string_index_of_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.StringIndexOf(firestore.FieldOf("title"), "The").As("indexOfThe"),
		)).
		Execute(ctx)
	// [END firestore_string_index_of_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
