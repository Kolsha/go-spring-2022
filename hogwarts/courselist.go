//go:build !solution

package hogwarts

func learn(res []string, course string, prereqs map[string][]string, learned map[string]bool, pk map[string]bool) []string {
	if pk[course] {
		panic("a cycle")
	}
	if learned[course] {
		return res
	}

	pk[course] = true

	for _, prereq := range prereqs[course] {
		res = learn(res, prereq, prereqs, learned, pk)
	}
	learned[course] = true
	pk[course] = false
	res = append(res, course)
	return res
}
func GetCourseList(prereqs map[string][]string) []string {
	learned := make(map[string]bool)
	pk := make(map[string]bool)
	res := make([]string, 0)
	for course := range prereqs {
		res = learn(res, course, prereqs, learned, pk)
	}
	return res
}
