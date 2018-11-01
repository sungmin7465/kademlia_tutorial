package kademlia

import (
	"container/list"
	"sort"

	"formalworks.com/kadem_tutorial/vector"
)

const BucketSize = 20

type RoutingTable struct {
	node    Contact
	buckets [IdLength * 8]*list.List
}

type ContactRecord struct {
	node    *Contact
	sortKey NodeID
}

func (rec *ContactRecord) Less(other interface{}) bool {
	return rec.sortKey.Less(other.(*ContactRecord).sortKey)
}

func NewRoutingTable(node *Contact) (ret *RoutingTable) {
	ret = new(RoutingTable)
	for i := 0; i < IdLength*8; i++ {
		ret.buckets[i] = list.New()
	}
	ret.node = *node
	return
}

func (table *RoutingTable) Update(contact *Contact) {
	prefix_length := contact.id.Xor(table.node.id).PrefixLen()
	bucket := table.buckets[prefix_length]

	var element *list.Element
	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(*Contact).id.Equals(table.node.id) {
			element = e
			return
		}
	}

	if element == nil {
		if bucket.Len() <= BucketSize {
			bucket.PushFront(contact)
		}
		// TODO: Handle insertion when the list is full by evicting old elements if
		// they don't respond to a ping.
	} else {
		bucket.MoveToFront(element.Value.(*list.Element))
	}
}

func copyToVector(start, end *list.Element, vec *vector.Vector, target NodeID) {
	for elt := start; elt != end; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		vec.Push(&ContactRecord{contact, contact.id.Xor(target)})
	}
}

func (table *RoutingTable) FindClosest(target NodeID, count int) (ret *vector.Vector) {
	ret = new(vector.Vector).Resize(0, count)

	bucket_num := target.Xor(table.node.id).PrefixLen()
	bucket := table.buckets[bucket_num]
	copyToVector(bucket.Front(), nil, ret, target)

	for i := 1; (bucket_num-i >= 0 || bucket_num+i < IdLength*8) && ret.Len() < count; i++ {
		if bucket_num-i >= 0 {
			bucket = table.buckets[bucket_num-i]
			copyToVector(bucket.Front(), nil, ret, target)
		}
		if bucket_num+i < IdLength*8 {
			bucket = table.buckets[bucket_num+i]
			copyToVector(bucket.Front(), nil, ret, target)
		}
	}

	sort.Sort(ret)
	if ret.Len() > count {
		ret.Cut(count, ret.Len())
	}
	return
}
