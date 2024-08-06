// Copyright 2022 Google LLC
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

// wordcount is an example that counts words in Shakespeare and includes Beam
// best practices.
//
// The input file defaults to a public data set containing the text of King
// Lear, by William Shakespeare. You can override it and choose your own input
// with --input.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/testing/passert"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

var (
	// By default, this example reads from a public dataset containing the text of
	// King Lear. Set this option to choose a different input file or glob.
	input = flag.String("input", "gs://apache-beam-samples/shakespeare/kinglear.txt", "File(s) to read.")

	// Set this required option to specify where to write the output.
	output = flag.String("output", "", "Output file (required).")
)

func init() {
	register.DoFn3x0[context.Context, string, func(string)](&extractFn{})
	register.Emitter1[string]()
	register.Function2x1(formatFn)
}

var (
	wordRE          = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)
	empty           = beam.NewCounter("extract", "emptyLines")
	smallWordLength = flag.Int("small_word_length", 6, "length of small words (default: 9)")
	smallWords      = beam.NewCounter("extract", "smallWords")
	lineLen         = beam.NewDistribution("extract", "lineLenDistro")
)

// extractFn is a DoFn that emits the words in a given line and keeps a count for small words.
type extractFn struct {
	SmallWordLength int `json:"smallWordLength"`
}

// ProcessElement for extractFn processes a line at a time, emitting each word in that line
// while keeping count of how many words are below the SmallWordLength threshold.
func (f *extractFn) ProcessElement(ctx context.Context, line string, emit func(string)) {
	lineLen.Update(ctx, int64(len(line)))
	if len(strings.TrimSpace(line)) == 0 {
		empty.Inc(ctx, 1)
	}
	for _, word := range wordRE.FindAllString(line, -1) {
		// increment the counter for small words if length of words is
		// less than small_word_length
		if len(word) < f.SmallWordLength {
			smallWords.Inc(ctx, 1)
		}
		emit(word)
	}
}

// Format formats a KV of a word and its count as a string.
func Format(s beam.Scope, counted beam.PCollection) beam.PCollection {
	return beam.ParDo(s, formatFn, counted)
}

// formatFn is a DoFn that formats a word and its count as a string.
func formatFn(w string, c int) string {
	return fmt.Sprintf("%s: %v", w, c)
}

// CountWords is a composite transform that counts the words of a PCollection
// of lines. It expects a PCollection of type string and returns a PCollection
// of type KV<string,int>. The Beam type checker enforces these constraints
// during pipeline construction.
func CountWords(s beam.Scope, lines beam.PCollection) beam.PCollection {
	s = s.Scope("CountWords")

	// Convert lines of text into individual words.
	col := beam.ParDo(s, &extractFn{SmallWordLength: *smallWordLength}, lines)

	// Count the number of times each word occurs.
	return stats.Count(s, col)
}

// WordCountFromPCol counts the words from a PCollection and validates it.
func WordCountFromPCol(s beam.Scope, in beam.PCollection, hash string, size int) {
	out := Format(s, CountWords(s, in))
	passert.Hash(s, out, "out", hash, size)
}

func main() {
	flag.Parse()
	beam.Init()

	if *output == "" {
		log.Fatal("No output provided")
	}

	p := beam.NewPipeline()
	s := p.Root()

	lines := textio.Read(s, *input)
	counted := CountWords(s, lines)
	formatted := beam.ParDo(s, formatFn, counted)
	textio.Write(s, *output, formatted)

	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
