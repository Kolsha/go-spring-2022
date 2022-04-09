//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return nil
	}

	maxOut := guests[0].CheckOutDate
	for _, g := range guests {
		maxOut = max(maxOut, g.CheckOutDate)
	}

	table := make([]int, maxOut+1)
	for _, g := range guests {
		table[g.CheckInDate]++
		table[g.CheckOutDate]--
	}

	res := make([]Load, 0)

	count := 0
	for date, delta := range table {
		count += delta
		if delta == 0 {
			continue
		}
		res = append(res, Load{StartDate: date, GuestCount: count})
	}
	return res
}
