// Copyright 2021, Dávid Mikuš. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// most of the code copied from encoding/base64 package.
package base45

import (
	"bytes"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
	"testing"
)

type testpair struct {
	decoded, encoded string
}

var pairs = []testpair{
	{
		decoded: "AB",
		encoded: "BB8",
	},
	{
		decoded: "Hello!",
		encoded: "%69 VD92E",
	},
	{
		decoded: "Hello!!",
		encoded: "%69 VD92EX0",
	},
	{
		decoded: "base-45",
		encoded: "UJCLQE7W581",
	},
	{
		decoded: "base-45-",
		encoded: "UJCLQE7W5NW6",
	},
}

var bigTest = testpair{
	decoded: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Fusce aliquam vestibulum ipsum. Fusce tellus. Duis ante orci, molestie vitae vehicula venenatis, tincidunt ac pede. In laoreet, magna id viverra tincidunt, sem odio bibendum justo, vel imperdiet sapien wisi sed libero. Integer rutrum, orci vestibulum ullamcorper ultricies, lacus quam ultricies odio, vitae placerat pede sem sit amet enim. Et harum quidem rerum facilis est et expedita distinctio. Integer malesuada. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Etiam bibendum elit eget erat. Nulla pulvinar eleifend sem. Pellentesque arcu. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat. Integer tempor. Duis condimentum augue id magna semper rutrum. Fusce consectetuer risus a nunc. In sem justo, commodo ut, suscipit at, pharetra vitae, orci. Pellentesque sapien. Maecenas libero.",
	encoded: "$T9ZKE ZD$ED$QE ZDGVC*VDBJEPQESUEBEC7$C1Q5UPCF/DZ C7WENWE5$C944AVCM9EJQEZEDU1D: C-EDI$5$+8JQEDZCAEC%EDY$E ZDO/E QENED0%E1%EH44W9E1%EI$5$+8JQEDZC7WE VD7%EI$5KT8+ED944G/DDZC04EOPC1Q5P$DTVD QEQEDU44-ED3ECU44+ COED0%EOCCO/E1/D..DBWE9PES44ZEDOPCMVCG/D944-NCI9E6VCI$5XC9K44DECZKE7$C1Q5B$DI3DOCCPEDU44/ED5$CVKES44ZEDOPCMVCG/D1Q5LQE ZDV3E EDA44NED1$CMVC ZDSKD QEK2EU44: CH44Q$D5$CAVC7$CR44EECQEDM-DE4FPQER44  CK44NED5$CM2EU34G/D* C5$CQ448%E6LE3 DN44XKE2DDO/E QENED0%E1%ET44 VDBECUPC1LE5$CT44:VD*KEOPC6$C1Q5PVD PC.OEKFEBECT44:VD*KEOPC6$CN44AVCK2EU44-ED3ECO44PVDKPCVKESUEI9E6VCR440$CR44-ED944F$DSUE1$CYEDI$55Z8G44GEC1%EP44.$E6VC ZDZKE6LE ZDJ.COPCXVD.OE6$CSUE7$CD447AF  C-EDOCCAVC QEZEDZPC EDI$5XC97WE93DBJEB$DTVD$QE2EC%CCQ34*9F3$C7WE6%ER44ZEDSUEU3EGPCZ CIECB444%EPEDIECIECM44-3EO440LEPED1$C+UER442%ESUEZEDB440%EE9EP44.$EN44O.COEDMEDC446$C5$C2%ESUEP$D VD-ED9445/D ZDPEDD44 QEK440EC04E1%EI$55Z8MED ZD%JCYJC0/D1%ED44XVDSUE* C7$CD44VKE.UEZ340%EPVDO440%ES/E..DBJE: C- CN.C0/DR440$CI$5G7A VD1$C7WEXQE%$E944XKES EU343WEKFEDZCX C6LE ZDZKE6LE ZD:8D-NC7WE1/DNWEBJEOCCHQEM9E1$C7WEC44: CZ CNWE9PET44SUEJECSUEZKEOEDQED0/D+EDU44*3E4%E3WEBWE5KC.OEB$D EDZKE.OEAECMED.OEUPCF/D4$CY$ENWEBJEJECSUEI9E KE5$C1$CAVC.OEGVC*VD*KE5KC.OEHECI9E*KE04E6$CQ443$C: CPVD.UEU34G/D* C5$CS440$CS9EPJEP34.$E.OEUPC0/DYED1$CNWE ZDJECP3DDZCPEDL445EC..DR440$CI9EBJE6LEKWE1%EI$5$+8JQEDZCUPCF/DZ C7WENWE5$CQ44+ED7%E944M442%EAOCU34M-DLQE ZDSKD QEK2EB44+3EP$DGVCT44+UER447%EOPCM9ESUEIEC1Q5L9EGEC7$CVKEU44-ED3EC1Q504EOPCI$5G7A VD1$C7WEXQE%$ER44EECQED -DY343ECKPC..D.OEXVDYJC0LE11",
}

func testEqual(t *testing.T, msg string, args ...interface{}) bool {
	t.Helper()
	if args[len(args)-2] != args[len(args)-1] {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func TestEncode(t *testing.T) {
	for _, p := range pairs {
		got := EncodeToString([]byte(p.decoded))
		testEqual(t, "Encode(%q) = %q, want %q", p.decoded,
			got, p.encoded)
	}
}

func TestEncoder(t *testing.T) {
	for _, p := range pairs {
		bb := &bytes.Buffer{}
		encoder := NewEncoder(bb)
		encoder.Write([]byte(p.decoded))
		encoder.Close()
		testEqual(t, "Encode(%q) = %q, want %q", p.decoded, bb.String(), p.encoded)
	}
}

func TestDecode(t *testing.T) {
	for _, p := range pairs {
		dbuf := make([]byte, DecodedLen(len(p.encoded)))
		count, err := Decode(dbuf, []byte(p.encoded))
		testEqual(t, "Decode(%q) = error %v, want %v", p.encoded, err, error(nil))
		testEqual(t, "Decode(%q) = length %v, want %v", p.encoded, count, len(p.decoded))
		testEqual(t, "Decode(%q) = %q, want %q", p.encoded, string(dbuf[0:count]), p.decoded)

		dbuf, err = DecodeString(p.encoded)
		testEqual(t, "DecodeString(%q) = error %v, want %v", p.encoded, err, error(nil))
		testEqual(t, "DecodeString(%q) = %q, want %q", p.encoded, string(dbuf), p.decoded)
	}
}

func TestDecoder(t *testing.T) {
	for _, p := range pairs {
		// We can't read everything in one read
		if len(p.encoded)%3 != 0 {
			continue
		}
		decoder := NewDecoder(strings.NewReader(p.encoded))
		dbuf := make([]byte, DecodedLen(len(p.encoded)))
		count, err := decoder.Read(dbuf)
		if err != nil && err != io.EOF {
			t.Fatal("Read failed", err)
		}
		testEqual(t, "Read from %q = length %v, want %v", p.encoded, count, len(p.decoded))
		testEqual(t, "Decoding of %q = %q, want %q", p.encoded, string(dbuf[0:count]), p.decoded)
		if err != io.EOF {
			_, err = decoder.Read(dbuf)
		}
		testEqual(t, "Read from %q = %v, want %v", p.encoded, err, io.EOF)
	}
}

func TestDecoderBuffering(t *testing.T) {
	for bs := 1; bs <= 12; bs++ {
		decoder := NewDecoder(strings.NewReader(bigTest.encoded))
		buf := make([]byte, len(bigTest.decoded)+12)
		var total int
		var n int
		var err error
		for total = 0; total < len(bigTest.decoded) && err == nil; {
			n, err = decoder.Read(buf[total : total+bs])
			total += n
		}
		if err != nil && err != io.EOF {
			t.Errorf("Read from %q at pos %d = %d, unexpected error %v", bigTest.encoded, total, n, err)
		}
		testEqual(t, "Decoding/%d of %q = %q, want %q", bs, bigTest.encoded, string(buf[0:total]), bigTest.decoded)
	}
}

func TestDecodeCorrupt(t *testing.T) {
	testCases := []struct {
		input  string
		offset int // -1 means no corruption.
	}{
		{"", -1},
		{"FGW", -1},
		{"GGW", 3},
		{"FaW", 1},
	}
	for _, tc := range testCases {
		dbuf := make([]byte, DecodedLen(len(tc.input)))
		_, err := Decode(dbuf, []byte(tc.input))
		if tc.offset == -1 {
			if err != nil {
				t.Error("Decoder wrongly detected corruption in", tc.input)
			}
			continue
		}
		switch err := err.(type) {
		case CorruptInputError:
			testEqual(t, "Corruption in %q at offset %v, want %v", tc.input, int(err), tc.offset)
		default:
			t.Error("Decoder failed to detect corruption in", tc)
		}
	}
}

func TestDecodeBounds(t *testing.T) {
	var buf [32]byte
	s := EncodeToString(buf[:])
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Decode panicked unexpectedly: %v\n%s", err, debug.Stack())
		}
	}()
	n, err := Decode(buf[:], []byte(s))
	if n != len(buf) || err != nil {
		t.Fatalf("StdEncoding.Decode = %d, %v, want %d, nil", n, err, len(buf))
	}
}

func TestBig(t *testing.T) {
	n := 3*1000 + 1
	raw := make([]byte, n)
	const alpha = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		raw[i] = alpha[i%len(alpha)]
	}
	encoded := new(bytes.Buffer)
	w := NewEncoder(encoded)
	nn, err := w.Write(raw)
	if nn != n || err != nil {
		t.Fatalf("Encoder.Write(raw) = %d, %v want %d, nil", nn, err, n)
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("Encoder.Close() = %v want nil", err)
	}
	decoded, err := io.ReadAll(NewDecoder(encoded))
	if err != nil {
		t.Fatalf("io.ReadAll(NewDecoder(...)): %v", err)
	}

	if !bytes.Equal(raw, decoded) {
		var i int
		for i = 0; i < len(decoded) && i < len(raw); i++ {
			if decoded[i] != raw[i] {
				break
			}
		}
		t.Errorf("Decode(Encode(%d-byte string)) failed at offset %d", n, i)
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		EncodeToString(data)
	}
}

func BenchmarkDecodeString(b *testing.B) {
	sizes := []int{2, 4, 8, 64, 8192}
	benchFunc := func(b *testing.B, benchSize int) {
		data := EncodeToString(make([]byte, benchSize))
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = DecodeString(data)
		}
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			benchFunc(b, size)
		})
	}
}
