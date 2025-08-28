package main

func plural[T ~uint | ~int](value T, singular, plural string) string {
	if value == 1 {
		return singular
	}
	return plural
}
