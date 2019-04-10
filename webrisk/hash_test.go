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
package webrisk

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"reflect"
	"runtime"
	"sync"
	"testing"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	"github.com/golang/protobuf/proto"
)

var (
	testHashesCached [][]hashPrefix
	loadOnce         sync.Once
)

func getTestHashes() [][]hashPrefix {
	loadOnce.Do(func() {
		b, err := ioutil.ReadFile("testdata/hashes.gob")
		if err != nil {
			panic(err)
		}
		r := gob.NewDecoder(bytes.NewReader(b))
		if err := r.Decode(&testHashesCached); err != nil {
			panic(err)
		}
	})
	return testHashesCached
}

func TestHashValidate(t *testing.T) {
	vectors := []struct {
		hashes hashPrefixes
		valid  bool
	}{
		{hashPrefixes{"bbb"}, false},
		{hashPrefixes{"bbbb"}, true},
		{hashPrefixes{"bbbb", "aaaa"}, false},
		{hashPrefixes{"aaaa", "bbbb"}, true},
		{hashPrefixes{"aaaa", "bbbbb"}, true},
		{hashPrefixes{"aaaa", "bbbbb", "bbbbb"}, false},
		{hashPrefixes{"aaaa", "bbbbb", "bbbbc"}, true},
		{hashPrefixes{"aaaa", "bbbbb", "bbbbbb"}, false},
		{hashPrefixes{"aaaa", "bbbbb", "bbbbbc"}, false},
		{hashPrefixes{"aaaa", "bbbbbd", "bbbbbc"}, false},
		{hashPrefixes{"aaaa", "bbbbbc", "bbbbbd"}, true},
	}

	for i, v := range vectors {
		valid := v.hashes.Validate() == nil
		if valid != v.valid {
			t.Errorf("test %d, Validate() = %v, want %v", i, valid, v.valid)
		}
	}
}

func TestHashComputeSHA256(t *testing.T) {
	vectors := []struct {
		hashes hashPrefixes
		sha256 string
	}{{
		hashes: nil,
		sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}, {
		hashes: hashPrefixes{"xxxx", "yyyy", "zzzz"},
		sha256: "20ffb2c3e9532153b96b956845381adc06095f8342fa2db1aafba6b0e9594d68",
	}, {
		hashes: hashPrefixes{"aaaa", "bbbb", "cccc", "dddd"},
		sha256: "147eb9dcde0e090429c01dbf634fd9b69a7f141f005c387a9c00498908499dde",
	}}

	for i, v := range vectors {
		sha256 := hex.EncodeToString(v.hashes.SHA256())
		if sha256 != v.sha256 {
			t.Errorf("test %d, mismatching hash:\ngot  %s\nwant %s", i, sha256, v.sha256)
		}
	}
}

func TestHashFromPattern(t *testing.T) {
	got := hashFromPattern("foo.com/")
	want := hashPrefix(mustDecodeHex(t, "af30010e4c011fc77e1ab03ef898ab2f8dec17c0bd178a875ad0cd9b7a6e5973"))
	if got != want {
		t.Errorf("FromPattern() = %x, want %x", got, want)
	}
}

func TestHashSet(t *testing.T) {
	var testHashes = getTestHashes()

	type hashQuery struct {
		hash hashPrefix
		len  int
	}
	type testVector struct {
		hashes  hashPrefixes
		queries []hashQuery
		fail    bool
	}

	vectors := []testVector{{
		hashes:  hashPrefixes{"aaaa", "bbbbb", "bbbbc"},
		queries: []hashQuery{{"aaaa", 4}, {"bbbb", 0}, {"bbbbbbbbb", 5}, {"bbbbcbbbb", 5}},
	}, {
		hashes:  hashPrefixes{"aaaa", "bbbb"},
		queries: []hashQuery{{"aaaa", 4}, {"aaaaaaa", 4}, {"bbbb", 4}, {"cccc", 0}},
	}, {
		hashes:  hashPrefixes{"abcdefgh", "abcdefgi", "abcdefgj"},
		queries: []hashQuery{{"abcd", 0}, {"abcde", 0}, {"abcdef", 0}, {"abcdefg", 0}, {"abcdefgh", 8}, {"abcdefgz", 0}},
	}}

	// Add hashes based on actual test data.
	if !testing.Short() {
		for _, hs := range testHashes {
			var v testVector
			v.hashes = hs
			for _, h := range hs {
				v.queries = append(v.queries, hashQuery{h + "footer", len(h)})
			}
			for _, h := range hs {
				v.queries = append(v.queries, hashQuery{"header" + h, 0})
			}
			vectors = append(vectors, v)
		}
	}

	var hs hashSet
	for i, v := range vectors {
		var fail bool
		func() {
			defer func() { fail = recover() != nil }()
			hs.Import(v.hashes)

			for j, q := range v.queries {
				n := hs.Lookup(q.hash)
				if n != q.len {
					t.Errorf("test %d.%d, Lookup(%q) = %d, want %d", i, j, q.hash, n, q.len)
				}
			}

			hashes := hs.Export()
			hashPrefixes(hashes).Sort()
			if !reflect.DeepEqual(hashes, v.hashes) {
				t.Errorf("test %d, output hashes mismatch\ngot  %q\nwant %q", i, hashes, v.hashes)
			}
		}()

		if fail != v.fail {
			if fail {
				t.Errorf("test %d, unexpected test failure", i)
			} else {
				t.Errorf("test %d, unexpected test success", i)
			}
		}
	}
}

func BenchmarkHashSet(b *testing.B) {
	var testHashes = getTestHashes()

	var queries []hashPrefix
	for _, h := range testHashes[1] {
		queries = append(queries, "header"+h)
		queries = append(queries, h+"footer")
	}

	var hs hashSet
	hs.Import(testHashes[1])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, h := range queries {
			hs.Lookup(h)
		}
	}
}

func BenchmarkHashSetMemory(b *testing.B) {
	var testHashes = getTestHashes()

	var ms1, ms2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&ms1)

	b.ReportAllocs()
	b.ResetTimer()
	var hs hashSet
	for i := 0; i < b.N; i++ {
		hs = hashSet{}
		hs.Import(testHashes[1])
	}

	runtime.GC()
	runtime.ReadMemStats(&ms2)
	b.Logf("mem_alloc: %dB, hashes: %dx", ms2.Alloc-ms1.Alloc, len(testHashes[1]))
}

func TestDecodeHashes(t *testing.T) {
	// These vectors were randomly generated using the server-side Rice
	// compression implementation.
	vectors := []struct {
		input  []string // Hex encoded serialized ThreatEntrySet proto
		output []string // Hex encoded expected hash prefixes
	}{{
		[]string{"122308a08fcb6d101c18062218dda588628aad88f883e2421a66384d10bce123dd22030202", "0a19081512151c9e466c435e51f99f059ff356185c730351d2f2b6"},
		[]string{"17f15426", "1c9e466c435e51f99f059ff356185c730351d2f2b6", "47ba02b7", "573373a2", "a0c7b20d", "a19edd3e", "d2c60aef", "f1fa25a2"},
	}, {
		[]string{"0a12080e120e8f991dc48f98c8647137d508974b", "0a1c081812180c698b1fc286b46c5ef5b96640b68a490e5135deebe15d02", "0a23081f121f40481597e49bc0768efbb174ca457f1b25eca550b611dd7385b49526b221cc"},
		[]string{"0c698b1fc286b46c5ef5b96640b68a490e5135deebe15d02", "40481597e49bc0768efbb174ca457f1b25eca550b611dd7385b49526b221cc", "8f991dc48f98c8647137d508974b"},
	}, {
		[]string{"0a0908051205f22897e85b", "0a11080d120da1e5504a06c508adac0441dcf5", "0a17081312139ccb416162bf1971b4017f2194026e6c309c91", "0a3408181230198dc5cba24feb2a0fba67e49bb747a8e242a62a194f4b1dc9ce9fb201c883313059b3438daeedb25160b0cfb64dbca3", "0a20081c121cb9181bc30742d0e5d1fb1bfa8f11603f6c39b2adfc83d0a4061ea490"},
		[]string{"198dc5cba24feb2a0fba67e49bb747a8e242a62a194f4b1d", "9ccb416162bf1971b4017f2194026e6c309c91", "a1e5504a06c508adac0441dcf5", "b9181bc30742d0e5d1fb1bfa8f11603f6c39b2adfc83d0a4061ea490", "c9ce9fb201c883313059b3438daeedb25160b0cfb64dbca3", "f22897e85b"},
	}, {
		[]string{"121808f8adb660101c1803220d5fe56d19f084ea8fffeacca708", "0a0b080712079a98ad5e44557e", "0a0d08091209ecd9224024b005f2ba", "0a12080e120ebc78a0caa431d7a503148002c4d7", "0a1c08181218c8783e05bad7fb601a6515e78b3212590af822dc1dca670a", "0a23081f121f9fa719c51211e2253306e021f181b75eb82accb20e9f495e000db597d4c350"},
		[]string{"2b9f705e", "8d4e735c", "9a98ad5e44557e", "9fa719c51211e2253306e021f181b75eb82accb20e9f495e000db597d4c350", "a085c4f2", "bc78a0caa431d7a503148002c4d7", "c8783e05bad7fb601a6515e78b3212590af822dc1dca670a", "ecd9224024b005f2ba", "f8960d0c"},
	}, {
		[]string{"121408dfea9d4e101c1802220991280dd01f2c9d5301"},
		[]string{"33341993", "5f75c709", "83bfca1d"},
	}, {
		[]string{"122008c2ee8cc801101c18052214b0e8007275b8ba8319e39f2d4ea7f1df8b1a1202"},
		[]string{"0f64be45", "42370319", "6650275f", "95ba6fd7", "9aab0322", "db7cbd52"},
	}, {
		[]string{"122b08ded89aba03101c1808221fce4f1cc9b81c12c9e0142610815a766d8771f0b282397d6779d2cbe98e2f01", "0a12080e120e6a32dc87b076fc4bf1b03bbc5528", "0a1e081a121a7ed093f0d8dd346a59bcef03c27941caea05f21ff6d593345f05"},
		[]string{"17ad48ca", "28471d40", "41e3df44", "45d4d43b", "51643a4b", "5eac4637", "6a32dc87b076fc4bf1b03bbc5528", "7e261bf6", "7ed093f0d8dd346a59bcef03c27941caea05f21ff6d593345f05", "beebab7b", "cc9d97ff"},
	}, {
		[]string{"121508e9f2f7ee01101c18022209ff93bf27d073d3bb37", "0a0c080812089d8cf590593d4033", "0a0e080a120a578aca0549f22647a724", "0a11080d120d74428c1c0cc815c0565bf9c717", "0a12080e120e3cea8c6ceba9149f70475cb3e622", "0a1d08191219b1409a6eab58a60f0ba3296daaa8fa11300f51cf4d3f25cec7", "0a23081f121f83d8521c22ef3dded0865c60efb48128098590d7cfc656a3518d8d5e60dda8"},
		[]string{"3cea8c6ceba9149f70475cb3e622", "578aca0549f22647a724", "5bf1e2c7", "69f9dd1d", "74428c1c0cc815c0565bf9c717", "83d8521c22ef3dded0865c60efb48128098590d7cfc656a3518d8d5e60dda8", "9d8cf590593d4033", "b1409a6eab58a60f0ba3296daaa8fa11300f51cf4d3f25cec7", "c96bdafe"},
	}, {
		[]string{"121108d793e4ff0a101c180122052f701aef01", "0a0f080b120baca3f0f63591bbeb500ef6"},
		[]string{"58dd71ff", "aca3f0f63591bbeb500ef6", "d709f9af"},
	}, {
		[]string{"12060895cd9af70d", "0a0e0805120a38c6373d5fa66ec3a2de", "0a0c0808120830ac22c848d99593", "0a1808141214a346b3224e920c47d054ab753256c2fbafe33a42", "0a1a08161216ce5f42ada2806e4c53dae6bde0aa7079682ced066dc3", "0a3a081b123620b4e512f3e760adc605751b9ba6a526b6c9a590132567f8a5ef53805511e1a0f84374fe2da9734fef68924b24abe734bc975a239874"},
		[]string{"20b4e512f3e760adc605751b9ba6a526b6c9a590132567f8a5ef53", "30ac22c848d99593", "38c6373d5f", "805511e1a0f84374fe2da9734fef68924b24abe734bc975a239874", "95a6e6de", "a346b3224e920c47d054ab753256c2fbafe33a42", "a66ec3a2de", "ce5f42ada2806e4c53dae6bde0aa7079682ced066dc3"},
	}}
loop:
	for i, v := range vectors {
		var got hashPrefixes
		for _, in := range v.input {
			set := &pb.ThreatEntryAdditions{}
			if err := proto.Unmarshal(mustDecodeHex(t, in), set); err != nil {
				t.Errorf("test %d, unexpected proto.Unmarshal error: %v", i, err)
				continue loop
			}

			hashes, err := decodeHashes(set)
			if err != nil {
				t.Errorf("test %d, unexpected decodeHashes error: %v", i, err)
				continue loop
			}
			got = append(got, hashes...)
		}

		got.Sort()
		want := make(hashPrefixes, 0, len(v.output))
		for _, h := range v.output {
			want = append(want, hashPrefix(mustDecodeHex(t, h)))
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("test %d, DecodeHashes() = %x; want %x", i, got, want)
		}
	}
}

func TestDecodeIndices(t *testing.T) {
	// These vectors were randomly generated using the server-side Rice
	// compression implementation.
	vectors := []struct {
		input  []string // Hex encoded serialized ThreatEntrySet proto
		output []int32  // Expected output indices
	}{{
		[]string{"121c08ac01101c18052213720000c0210000100400001a01006017000000"},
		[]int32{172, 229, 364, 494, 776, 963},
	}, {
		[]string{"1222084b101c1807221a34010000110000300500000a0000e0100000a80100007a000000"},
		[]int32{75, 229, 297, 463, 473, 608, 714, 958},
	}, {
		[]string{"121e0823101c18062216f80100800f000050050000c50000c0020000b4020000"},
		[]int32{35, 287, 349, 519, 716, 738, 911},
	}, {
		[]string{"120308e607"},
		[]int32{998},
	}, {
		[]string{"121408c101101c1803220b360300c02b0000b8040000"},
		[]int32{193, 604, 779, 930},
	}, {
		[]string{"1223088001101c1807221a140000c000000088050000520000c03200004402000024000000"},
		[]int32{128, 138, 141, 318, 400, 806, 951, 1023},
	}, {
		[]string{"121c088f02101c1805221348000000560000c0010000cd00000002000000"},
		[]int32{271, 307, 651, 707, 912, 928},
	}, {
		[]string{"121808f103101c1804220fde0000c00000006805000085000000"},
		[]int32{497, 608, 611, 784, 917},
	}}

loop:
	for i, v := range vectors {
		var got []int32
		for _, in := range v.input {
			set := &pb.ThreatEntryRemovals{}
			if err := proto.Unmarshal(mustDecodeHex(t, in), set); err != nil {
				t.Errorf("test %d, unexpected proto.Unmarshal error: %v", i, err)
				continue loop
			}

			indices, err := decodeIndices(set)
			if err != nil {
				t.Errorf("test %d, unexpected decodeIndices error: %v", i, err)
				continue loop
			}
			got = append(got, indices...)
		}
		if !reflect.DeepEqual(got, v.output) {
			t.Errorf("failure in test case %d. DecodeIndices() = %v; want %v", i, got, v.output)
		}
	}
}

func TestRiceDecoder(t *testing.T) {
	// These vectors were randomly generated using the server-side Rice
	// compression implementation.
	vectors := []struct {
		k      uint32   // Golomb-Rice value
		input  string   // The input encoded string in hexadecimal format
		output []uint32 // The expected output values
	}{{
		2,
		"f702",
		[]uint32{15, 9},
	}, {
		5,
		"00",
		[]uint32{0},
	}, {
		10,
		"",
		[]uint32{},
	}, {
		28,
		"54607be70a5fc1dcee69defe583ca3d6a5f2108c4a595600",
		[]uint32{62763050, 1046523781, 192522171, 1800511020, 4442775, 582142548},
	}, {
		28,
		"06861b2314cb46f2af0708c988541f4104d51a03ebe63a8013917bbf83f3b785f12918b36109",
		[]uint32{26067715, 344823336, 8420095, 399843890, 95029378, 731622412, 35811335, 1047558127, 1117722715, 78698892},
	}, {
		27,
		"8998d875bc4491eb390c3e309a78f36ad4d9b19ffb703e443ea3086742c22b46698e3cebd9105a439a32a52d4e770f877820b6ab7198480c9e9ed7230c13432ca901",
		[]uint32{225846818, 328287420, 166748623, 29117720, 552397365, 350353215, 558267528, 4738273, 567093445, 28563065, 55077698, 73091685, 339246010, 98242620, 38060941, 63917830, 206319759, 137700744},
	}, {
		28,
		"21c50291f982d757b8e93cf0c84fe8648d776204d6853f1c9700041b17c6",
		[]uint32{339784008, 263128563, 63871877, 69723256, 826001074, 797300228, 671166008, 207712688},
	}, {
		28,
		"959c7db08fe8d9bdfe8c7f81530d75dc4e40180c9a453da8dcfa2659409e16084377c34e0401a4e65d00",
		[]uint32{471820069, 196333855, 855579133, 122737976, 203433838, 85354544, 1307949392, 165938578, 195134475, 553930435, 49231136},
	}, {
		27,
		"1a4f692a639af6c62eaf73d06fd731eb771d43e32b93ce678b59f998d4da4f3c6fb0e8a5788d623618fe081e78d814322484611cf33763c4a0887b74cb64c85cba05",
		[]uint32{87336845, 129291033, 30906211, 433549264, 30899891, 53207875, 11959529, 354827862, 82919275, 489637251, 53561020, 336722992, 408117728, 204506246, 188216092, 9047110, 479817359, 230317256},
	}, {
		28,
		"f1940a876c5f9690e3abf7c0cb2de976dbf85963c16f7c99e3875fc704deb9468e54c0ac4a030d6c8f00",
		[]uint32{297968956, 19709657, 259702329, 76998112, 1023176123, 29296013, 1602741145, 393745181, 177326295, 55225536, 75194472},
	}, {
		28,
		"412ce4fe06dc0dbd31a504d56edd9b43b73f11245210804f964bd48067b2dd52c94e02c6d760de0692521edd356471262cfecf8146b27901",
		[]uint32{532220688, 780594691, 436816483, 163436269, 573044456, 1069604, 39629436, 211410997, 227714491, 381562898, 75610008, 196754597, 40310339, 15204118, 99010842},
	}, {
		28,
		"b22c263acd669cdb5f072e6fe6f9211052d594f4822248f99d24f6ff2ffc6d3f21651b363456eac42100",
		[]uint32{219354713, 389598618, 750263679, 554684211, 87381124, 4523497, 287633354, 801308671, 424169435, 372520475, 277287849},
	}}

loop:
	for i, v := range vectors {
		br := newBitReader(mustDecodeHex(t, v.input))
		rd := newRiceDecoder(br, v.k)

		vals := []uint32{}
		for i := 0; i < len(v.output); i++ {
			val, err := rd.ReadValue()
			if err != nil {
				t.Errorf("test %d, unexpected error: %v", i, err)
				continue loop
			}
			vals = append(vals, val)
		}
		if !reflect.DeepEqual(vals, v.output) {
			t.Errorf("test %d, output mismatch:\ngot  %v\nwant %v", i, vals, v.output)
		}
	}
}

func TestBitReader(t *testing.T) {
	vectors := []struct {
		cnt int    // Number of bits to read
		val uint32 // Expected output value to read
		rem int    // Number of bits remaining in the bitReader
	}{
		{cnt: 0, val: 0, rem: 56},
		{cnt: 1, val: 1, rem: 55},
		{cnt: 1, val: 0, rem: 54},
		{cnt: 1, val: 1, rem: 53},
		{cnt: 1, val: 1, rem: 52},
		{cnt: 8, val: 0x20, rem: 44},
		{cnt: 32, val: 0x40000000, rem: 12},
		{cnt: 9, val: 0x00000170, rem: 3},
		{cnt: 3, val: 0x00000001, rem: 0},
	}

	// Test bitReader with data.
	br := newBitReader(mustDecodeHex(t, "0d020000000437"))
	for i, v := range vectors {
		val, err := br.ReadBits(v.cnt)
		if err != nil {
			t.Errorf("test %d, unexpected error: %v", i, err)
		}
		if val != v.val {
			t.Errorf("test %d, ReadBits() = 0x%08x, want 0x%08x", i, val, v.val)
		}
		if rem := br.BitsRemaining(); rem != v.rem {
			t.Errorf("test %d, BitsRemaining() = %d, want %d", i, rem, v.rem)
		}
	}

	// Test empty bitReader.
	br = newBitReader(mustDecodeHex(t, ""))
	if rem := br.BitsRemaining(); rem != 0 {
		t.Errorf("BitsRemaining() = %d, want 0", rem)
	}
	if _, err := br.ReadBits(1); err == nil {
		t.Errorf("unexpected ReadBits success")
	}
}
