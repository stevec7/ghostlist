package ghostlist

import (
    "sort"
)

/*
func (s *sortListT) Sort() {

}
*/

func sortHostlist(hostlist *[]string) error {
	sort.Strings(*hostlist)

	return nil
}
