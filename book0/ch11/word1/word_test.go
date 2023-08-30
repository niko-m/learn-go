package word

import "testing"

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") {
		t.Error(`IsPalindrome("detartrated") = false`)
	}
	if !IsPalindrome("kayak") {
		t.Error(`IsPalindrome("kayak") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
}

func TestFrenchPalindrome(t *testing.T) {
	if !IsPalindrome("ete") {
		t.Error(`IsPalindrome("ete") = false`)
	}
}

func TestCanalPalindrome(t *testing.T) {
	in := "A man, a plan, a canal: Panama"

	if !IsPalindrome(in) {
		t.Errorf("IsPalindrome(%q) = false", in)
	}
}
