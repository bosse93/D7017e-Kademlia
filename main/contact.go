package main

import (
	"fmt"
	"sort"
)

type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

// NewContact returns a type Contact with fields ID, address and distance.
// Distance is always nil.
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

//CalcDistance sets distance field of contact to distance to target.
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less compares two Contact distances and returns the lowest one.
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns Contact formated as string.
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

type ContactCandidates struct {
	contacts []Contact
}

// NewContactCandidates return new ContactCandidates object.
func NewContactCandidates() ContactCandidates {
	return ContactCandidates{}
}

// Append two ContactCandidates lists together to one.
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns count first Contacts in ContactCandidates list.
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort ContactCandidates list on distance.
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns length of ContactCandidates list.
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap two Contacts in ContatCandidates list.
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less compares two Contacts distances and returns the lowest one.
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}
