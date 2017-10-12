package main

const bucketSize = 20

type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
}

// NewRoutingTable initializes RoutingTable object.
// Populates RoutingTable bucket list with new Bucket objects.
// Returns RoutingTable as *RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

// AddContact adds a contact to routing table.
// If contact is RoutingTable owner he is not added.
// Finds appropriate Bucket and add contact to it.
func (routingTable *RoutingTable) AddContact(contact Contact) {
	if contact.ID != routingTable.me.ID {
		bucketIndex := routingTable.GetBucketIndex(contact.ID)
		bucket := routingTable.buckets[bucketIndex]
		bucket.AddContact(contact)
	}
}

// AddContactNetwork adds contact to RoutingTable.
// Finds appropriate Bucket and add Contact.
func (routingTable *RoutingTable) AddContactNetwork(contact Contact, network *Network) {
	bucketIndex := routingTable.GetBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContactNetwork(contact, network)
}

// FindClosestContacts finds 20 closest contacts to target in RoutingTable
// Returns a sorted list with Contacts. Closest to target first.
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.GetBucketIndex(target)

	bucket := routingTable.buckets[bucketIndex]
	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}
	return candidates.GetContacts(count)
}

// GetBucketIndex finds which bucket id is supposed to be added to.
func (routingTable *RoutingTable) GetBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}
