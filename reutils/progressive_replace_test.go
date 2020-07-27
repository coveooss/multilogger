package reutils

import (
	"fmt"
	"regexp"
)

func ExampleProgressiveReplaceAll() {
	re := regexp.MustCompile(`(?s){[^x]*}`)
	code := `{a
	{b
		{c
			{d

			d}
		c}
	b}
a}

adfsfdasx

{

}

cx
{{{{

}}}}
   `

	var i int
	new := ProgressiveReplaceAll(re, code, func(s string, pos *int) string {
		i++
		*pos = 5
		return fmt.Sprint("[", i)
	})
	fmt.Println(new)

	// Output:
}
