package scanner

import "19/matrix"

type uniqueArr []matrix.Vector

func (u *uniqueArr) add(v matrix.Vector) {
	for _, value := range *u {
		if value == v {
			return
		}
	}
	*u = append(*u, v)
}