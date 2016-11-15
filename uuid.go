package xorm

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"math"
	"math/big"
	"net"
	"sort"
	"strings"
	"time"
)

// The UUID represents Universally Unique IDentifier (which is 128 bit long).
type UUID [16]byte

var (
	// NIL is defined in RFC 4122 section 4.1.7.
	// The nil UUID is special form of UUID that is specified to have all 128 bits set to zero.
	NIL = &UUID{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	// NameSpaceDNS assume name to be a fully-qualified domain name.
	// Declared in RFC 4122 Appendix C.
	NameSpaceDNS = &UUID{
		0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceURL assume name to be a URL.
	// Declared in RFC 4122 Appendix C.
	NameSpaceURL = &UUID{
		0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceOID assume name to be an ISO OID.
	// Declared in RFC 4122 Appendix C.
	NameSpaceOID = &UUID{
		0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
	// NameSpaceX500 assume name to be a X.500 DN (in DER or a text output format).
	// Declared in RFC 4122 Appendix C.
	NameSpaceX500 = &UUID{
		0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8,
	}
)

// Version of the UUID represents a kind of subtype specifier.
func (u *UUID) Version() int {
	return int(binary.BigEndian.Uint16(u[6:8]) >> 12)
}

// String returns the human readable form of the UUID.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u *UUID) WithoutDashString() string {
	return fmt.Sprintf("%x", u[:])
}

func (u *UUID) variantRFC4122() {
	u[8] = (u[8] & 0x3f) | 0x80
}

// NewV3 creates a new UUID with variant 3 as described in RFC 4122.
// Variant 3 based namespace-uuid and name and MD-5 hash calculation.
func NewV3(namespace *UUID, name []byte) *UUID {
	uuid := newByHash(md5.New(), namespace, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x30
	return uuid
}

func newByHash(hash hash.Hash, namespace *UUID, name []byte) *UUID {
	hash.Write(namespace[:])
	hash.Write(name[:])

	var uuid UUID
	copy(uuid[:], hash.Sum(nil)[:16])
	uuid.variantRFC4122()
	return &uuid
}

type stamp [10]byte

var (
	mac      []byte
	requests chan bool
	answers  chan stamp
)

const gregorianUnix = 122192928000000000 // nanoseconds between gregorion zero and unix zero

func init() {
	mac = make([]byte, 6)
	rand.Read(mac)
	requests = make(chan bool)
	answers = make(chan stamp)
	go unique()
	i, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, d := range i {
		if len(d.HardwareAddr) == 6 {
			mac = d.HardwareAddr[:6]
			return
		}
	}
}

// NewV1 creates a new UUID with variant 1 as described in RFC 4122.
// Variant 1 is based on hosts MAC address and actual timestamp (as count of 100-nanosecond intervals since
// 00:00:00.00, 15 October 1582 (the date of Gregorian reform to the Christian calendar).
func NewV1() *UUID {
	var uuid UUID
	requests <- true
	s := <-answers
	copy(uuid[:4], s[4:])
	copy(uuid[4:6], s[2:4])
	copy(uuid[6:8], s[:2])
	uuid[6] = (uuid[6] & 0x0f) | 0x10
	copy(uuid[8:10], s[8:])
	copy(uuid[10:], mac)
	uuid.variantRFC4122()
	return &uuid
}

func unique() {
	var (
		lastNanoTicks uint64
		clockSequence [2]byte
	)
	rand.Read(clockSequence[:])

	for range requests {
		var s stamp
		nanoTicks := uint64((time.Now().UTC().UnixNano() / 100) + gregorianUnix)
		if nanoTicks < lastNanoTicks {
			lastNanoTicks = nanoTicks
			rand.Read(clockSequence[:])
		} else if nanoTicks == lastNanoTicks {
			lastNanoTicks = nanoTicks + 1
		} else {
			lastNanoTicks = nanoTicks
		}
		binary.BigEndian.PutUint64(s[:], lastNanoTicks)
		copy(s[8:], clockSequence[:])
		answers <- s
	}
}

// NewV4 creates a new UUID with variant 4 as described in RFC 4122. Variant 4 based on pure random bytes.
func NewV4() *UUID {
	buf := make([]byte, 16)
	rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	var uuid UUID
	copy(uuid[:], buf[:])
	uuid.variantRFC4122()
	return &uuid
}

// NewV5 creates a new UUID with variant 5 as described in RFC 4122.
// Variant 5 based namespace-uuid and name and SHA-1 hash calculation.
func NewV5(namespaceUUID *UUID, name []byte) *UUID {
	uuid := newByHash(sha1.New(), namespaceUUID, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x50
	return uuid
}

// NewNamespaceUUID creates a namespace UUID by using the namespace name in the NIL name space.
// This is a different approach as the 4 "standard" namespace UUIDs which are timebased UUIDs (V1).
func NewNamespaceUUID(namespace string) *UUID {
	return NewV5(NIL, []byte(namespace))
}

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

func (u *UUID) UnmarshalText(text []byte) (err error) {
	if len(text) < 32 {
		err = fmt.Errorf("uuid: UUID string too short: %s", text)
		return
	}

	t := text[:]
	braced := false

	if bytes.Equal(t[:9], urnPrefix) {
		t = t[9:]
	} else if t[0] == '{' {
		braced = true
		t = t[1:]
	}

	b := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			if t[0] != '-' {
				err = fmt.Errorf("uuid: invalid string format")
				return
			}
			t = t[1:]
		}

		if len(t) < byteGroup {
			err = fmt.Errorf("uuid: UUID string too short: %s", text)
			return
		}

		if i == 4 && len(t) > byteGroup &&
			((braced && t[byteGroup] != '}') || len(t[byteGroup:]) > 1 || !braced) {
			err = fmt.Errorf("uuid: UUID string too long: %s", text)
			return
		}

		_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])
		if err != nil {
			return
		}

		t = t[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

type StringSet struct {
	set    map[string]bool
	list   []string
	sorted bool
}

func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool), make([]string, 0), false}
}

func (set *StringSet) Add(i string) bool {
	_, found := set.set[i]
	set.set[i] = true
	if !found {
		set.sorted = false
	}
	return !found //False if it existed already
}

func (set *StringSet) Contains(i string) bool {
	_, found := set.set[i]
	return found //true if it existed already
}

func (set *StringSet) Remove(i string) {
	set.sorted = false
	delete(set.set, i)
}

func (set *StringSet) Len() int {
	return len(set.set)
}

func (set *StringSet) ItemByIndex(idx int) string {
	set.Sort()
	return set.list[idx]
}

func (set *StringSet) Index(c string) int {
	for i, s := range set.list {
		if c == s {
			return i
		}
	}
	return 0
}

func (set *StringSet) Sort() {
	if set.sorted {
		return
	}
	set.list = make([]string, 0)
	for s, _ := range set.set {
		set.list = append(set.list, s)
	}
	sort.Strings(set.list)
	set.sorted = true
}

func (set *StringSet) String() string {
	set.Sort()
	return strings.Join(set.list, "")
}

const (
	DEFAULT_ALPHABET = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type ShortUUID struct {
	alphabet *StringSet
}

func NewShortUUID() *ShortUUID {
	suid := &ShortUUID{}
	suid.SetAlphabet(DEFAULT_ALPHABET)
	return suid
}

func NewShortUUIDWithAlphabet(alphabet string) *ShortUUID {

	suuid := &ShortUUID{}
	if alphabet == "" {
		alphabet = DEFAULT_ALPHABET
	}
	suuid.SetAlphabet(alphabet)
	return suuid
}

func (s *ShortUUID) SetAlphabet(alphabet string) {
	set := NewStringSet()
	for _, a := range alphabet {
		set.Add(string(a))
	}
	set.Sort()
	s.alphabet = set
}

func (s ShortUUID) String() string {
	return s.UUID("")
}

var (
	NamespaceDNS, _  = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

func (s *ShortUUID) UUID(name string) string {
	var _uuid *UUID
	if name == "" {
		_uuid = NewV4()
	} else if strings.HasPrefix(name, "http") {
		_uuid = NewV5(&NamespaceDNS, []byte(name))
	} else {
		_uuid = NewV5(&NamespaceURL, []byte(name))
	}

	return s.Encode(_uuid)
}

// Encodes a UUID into a string (LSB first) according to the alphabet
// If leftmost (MSB) bits 0, string might be shorter
func (s *ShortUUID) Encode(uuid *UUID) string {
	padLen := s.encodeLen(len(uuid))
	number := uuidToInt(uuid)
	return s.numToString(number, padLen)
}

func (s *ShortUUID) Decode(input string) (UUID, error) {
	_uuid, err := FromString(s.stringToNum(input))
	return _uuid, err
}

func (s *ShortUUID) encodeLen(numBytes int) int {
	factor := math.Log(float64(25)) / math.Log(float64(s.alphabet.Len()))
	length := math.Ceil(factor * float64(numBytes))
	return int(length)
}

//Covert a number to a string, using the given alphabet.
func (s *ShortUUID) numToString(number *big.Int, padToLen int) string {
	output := ""
	var digit *big.Int
	for number.Uint64() > 0 {
		number, digit = new(big.Int).DivMod(number, big.NewInt(int64(s.alphabet.Len())), new(big.Int))
		output += s.alphabet.ItemByIndex(int(digit.Int64()))
	}
	if padToLen > 0 {
		remainer := math.Max(float64(padToLen)-float64(len(output)), 0)
		output = output + strings.Repeat(s.alphabet.ItemByIndex(0), int(remainer))
	}

	return output
}

// Convert a string to a number(based uuid string),using the given alphabet.
func (s *ShortUUID) stringToNum(input string) string {
	n := big.NewInt(0)
	for i := len(input) - 1; i >= 0; i-- {
		n.Mul(n, big.NewInt(int64(s.alphabet.Len())))
		n.Add(n, big.NewInt(int64(s.alphabet.Index(string(input[i])))))
	}

	x := fmt.Sprintf("%x", n)
	x = x[0:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:32]
	return x
}

func uuidToInt(_uuid *UUID) *big.Int {
	var i big.Int
	i.SetString(strings.Replace(_uuid.String(), "-", "", 4), 16)
	return &i
}
