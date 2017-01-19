package main

/* Behavior for numbers in Shenzhen IO.
   They are capped from -999 to 999
	 with no overflow.
*/

type Number int

func (n Number) GetValue() int {
	return int(n)
}

func (n *Number) Add(i int) {
	new_n := int(*n) + i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}

func (n *Number) Sub(i int) {
	new_n := int(*n) - i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}
