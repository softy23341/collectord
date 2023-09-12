package util

import "sort"

// Int64Slice TBD
type Int64Slice []int64

// UniqInt64 TBD
func UniqInt64(ids Int64Slice) Int64Slice {
	idsMap := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		idsMap[id] = struct{}{}
	}

	var idsList []int64
	for id := range idsMap {
		idsList = append(idsList, id)
	}
	return idsList
}

func (sl Int64Slice) Len() int           { return len(sl) }
func (sl Int64Slice) Swap(i, j int)      { sl[i], sl[j] = sl[j], sl[i] }
func (sl Int64Slice) Less(i, j int) bool { return sl[i] < sl[j] }

// Sort TBD
func Sort(ids Int64Slice) Int64Slice {
	sort.Sort(ids)
	return ids
}

// Int64Set TBD
type Int64Set map[int64]struct{}

// Contains TBD
func (s Int64Set) Contains(x int64) (contains bool) {
	_, contains = s[x]
	return
}

// Join TBD
func (s Int64Set) Join(other Int64Set) {
	for x := range other {
		s[x] = struct{}{}
	}
}

// ToSlice TBD
func (s Int64Set) ToSlice() (slice []int64) {
	slice = make([]int64, 0, len(s))
	for x := range s {
		slice = append(slice, x)
	}
	return
}

// ToSet TBD
func (sl Int64Slice) ToSet() Int64Set {
	return Int64SliceToSet([]int64(sl))
}

// All TBD
func (sl Int64Slice) All(x int64) bool {
	for _, el := range sl {
		if el != x {
			return false
		}
	}
	return true
}

// Uniq TBD
func (sl Int64Slice) Uniq() Int64Slice {
	return sl.ToSet().ToSlice()
}

// Delete TBD
func (sl Int64Slice) Delete(x int64) Int64Slice {
	res := make(Int64Slice, 0, len(sl))
	for _, v := range sl {
		if v != x {
			res = append(res, v)
		}
	}
	return res
}

// DeleteAll TBD
func (sl Int64Slice) DeleteAll(rms []int64) Int64Slice {
	res := make(Int64Slice, 0, len(sl))
	for _, v := range sl {
		skip := false
		for _, rm := range rms {
			if v == rm {
				skip = true
			}
		}
		if skip {
			continue
		}
		res = append(res, v)
	}
	return res
}

// PushBack TBD
func (sl Int64Slice) PushBack(x int64) Int64Slice {
	res := make(Int64Slice, 0, len(sl)+1)
	res = append(res, []int64{}...)
	res = append(res, x)
	return res
}

// Int64SliceToSet TBD
func Int64SliceToSet(slice []int64) (set Int64Set) {
	set = make(Int64Set, len(slice))
	for _, x := range slice {
		set[x] = struct{}{}
	}
	return
}

// Equal TBD
func (sl Int64Slice) Equal(another Int64Slice) bool {
	found := false

	// naive version
	for _, r := range sl {
		found = false
		for _, a := range another {
			if a == r {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
