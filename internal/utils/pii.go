package utils

import "strings"

func TotallyRedactDetails(item string) string {
	n := len(item)
	return strings.Repeat("*", n)
}


/*
	This is for card details, BVN and NIN details.
	It would obscure the first 4 digits and last 3 digits of the original.
*/
func PartiallyRedactDetails(item string) string {
	n := len(item)
	if n < 7 {
		return TotallyRedactDetails(item)
	}
	
	// Mask all but positions 4,5,6
	return strings.Repeat("*", 4) + item[4:7] + strings.Repeat("*", n-7)
}