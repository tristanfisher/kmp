package kmp

// preprocessTable pre-processes the search token and creates an array suitable for resuming search at the point of a failed match.
func preprocessTable(searchToken string) []int {

	// this is also a suffix for the prior index.  initialized as 0.
	// longestPrefix tracks the longest length of re-usable substrings in our searchToken
	longestPrefix := 0
	lenTarget := len(searchToken)

	// table tracks how much of a pattern can be (re-)used after a failed char comparison.
	// by definition, a string of length 1 has no prefix
	table := make([]int, lenTarget)

	// search for repeated substrings in a searchToken, bypassing the zeroth value of pattern as it has no usable prefix
	idx := 1
	for idx < lenTarget {

		// if the char at this index is the same as the longest match in the token, increase the longest prefix count
		if searchToken[idx] == searchToken[longestPrefix] {
			longestPrefix += 1
			table[idx] = longestPrefix
			idx++
			continue
		}

		// no repeated / substring found for pattern. set re-usable substring at this position to 0.
		if longestPrefix == 0 {
			table[idx] = 0
			idx++
			continue
		}

		// no matching prefix found with a non-zero longest prefix.
		// set the longest prefix to the previous largest value when non-zero.
		longestPrefix = table[idx-1]
		idx++
	}

	return table
}

// Search implements the Knuth-Morris-Pratt search algorithm.
// KMP uses the structure of a search input to avoid duplicate comparisons instead of restarting
// the search after each failed match of the complete search input.
//
// This is an improvement over searching a text for the first character in a target pattern, starting
// again when text[x] != targetPattern[y].
// This function takes two arguments:
// - body, which is the text to search
// - searchToken, is the target token for a search.
//
// The return is the position of matches and the count of matches.
func Search(body string, searchToken string) ([]int, int) {

	lenSearchSpace := len(body)
	lenTarget := len(searchToken)

	// zero search space or zero searchToken token.  no usable match logic.
	if lenSearchSpace == 0 || lenTarget == 0 {
		return []int{}, 0
	}

	table := preprocessTable(searchToken)

	idxBodySearch := 0
	idxTarget := 0

	// zero slice instead of declaration that could be nil
	startingPositionMatches := make([]int, 0)

	// walk each char in the search space
	for idxBodySearch < lenSearchSpace {
		// mismatch or no match
		if body[idxBodySearch] != searchToken[idxTarget] {
			// if we previously had a partial match of the token,
			// use value of previous index to skip unnecessary comparisons
			if idxTarget != 0 {
				idxTarget = table[idxTarget-1]
				continue
			}
			idxBodySearch++
			continue
		}

		// char match, advance search in both searchToken and body
		idxBodySearch++
		idxTarget++

		// complete token match in body
		if idxTarget == lenTarget {
			startingPositionMatches = append(startingPositionMatches, idxBodySearch-idxTarget)
			// use value of previous index to skip unnecessary comparisons
			idxTarget = table[idxTarget-1]
			// implicit continue
		}
	}

	return startingPositionMatches, len(startingPositionMatches)
}
