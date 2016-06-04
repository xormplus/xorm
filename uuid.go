package xorm

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"net"
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
