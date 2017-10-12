package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

const IDLength = 20

type KademliaID [IDLength]byte

// NewKademliaID encodes data string to 160 bit hex ID.
// Returns ID as type KademliaID
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

// NewRandomKademliaID randomizes a 160 bit hex ID.
// Returns ID as type KademliaID
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// HashKademliaID encodes fileName to a 160 bit hex ID.
// Maximum 19 characters is allowed in fileName.
// Returns ID as KademliaID.
func HashKademliaID(fileName string) *KademliaID {
	fmt.Println("Fil Namn: " + fileName)
	f := hex.EncodeToString([]byte(fileName))
	if len(f) > 38 {
		fmt.Println(f)
		fmt.Println("Name of file can be maximum 19 characters, including file extension.")
	}
	f = f + "03"
	for len(f) < 40 {
		f = f + "01"
	}
	return NewKademliaID(f)
}

// DecodeHash will decode a hash made if HashKademliaID.
// Returns fileName as string.
func DecodeHash(hash string) string {
	byteArray := []byte(hash)

	for i := 0; i < 19; i++ {
		if (string(byteArray[i*2]) == "0") && (string(byteArray[(i*2)+1]) == "3") {
			fileName, _ := hex.DecodeString(string(byteArray[:(i)*2]))
			return string(fileName)
		}
	}
	return "Error when decoding dataID"
}

// Less compares kademliaID and otherKademliaID.
// Returns true if kademliaID is less than otherKAdemliaID else false.
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if both kademliaIDs is the same.
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance calculates the xor distance between two KademliaIDs.
// Returns distance as a KademliaID.
func (kademliaID *KademliaID) CalcDistance(target *KademliaID) *KademliaID {

	result := KademliaID{}

	for i := 0; i < IDLength; i++ {

		result[i] = kademliaID[i] ^ target[i]
	}

	return &result
}

// String returns KademliaID as a string.
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
