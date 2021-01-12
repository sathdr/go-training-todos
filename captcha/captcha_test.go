package captcha_test

import (
	"fmt"
	"testing"

	"github.com/pallat/todos/captcha"

	"github.com/stretchr/testify/assert"
)

func TestCaptchaPattern1(t *testing.T) {
	operands := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for givingOperand, want := range operands {
		t.Run(fmt.Sprintf("right operand %d", givingOperand), func(t *testing.T) {
			p := 1
			lo := 1
			op := 1
			ro := givingOperand

			want := fmt.Sprintf("1 + %s", want)

			cc := captcha.New(p, lo, op, ro)
			get := cc.String()
			assert.Equal(t, get, want, fmt.Sprintf("given %d %d %d %d want %q but get %q", p, lo, op, ro, want, get))
		})
	}

	for givingOperand := 1; givingOperand <= 9; givingOperand++ {
		t.Run(fmt.Sprintf("left operand %d", givingOperand), func(t *testing.T) {
			p := 1
			lo := givingOperand
			op := 1
			ro := 1

			want := fmt.Sprintf("%d + one", givingOperand)

			cc := captcha.New(p, lo, op, ro)
			get := cc.String()
			assert.Equal(t, get, want, fmt.Sprintf("given %d %d %d %d want %q but get %q", p, lo, op, ro, want, get))
		})
	}

}

func TestCaptchaPattern2(t *testing.T) {
	operands := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for givingOperand, want := range operands {
		t.Run(fmt.Sprintf("left operand %d", givingOperand), func(t *testing.T) {
			p := 2
			lo := givingOperand
			op := 1
			ro := 1

			want := fmt.Sprintf("%s + 1", want)

			cc := captcha.New(p, lo, op, ro)
			get := cc.String()
			assert.Equal(t, get, want, fmt.Sprintf("given %d %d %d %d want %q but get %q", p, lo, op, ro, want, get))
		})
	}

	for givingOperand := 1; givingOperand <= 9; givingOperand++ {
		t.Run(fmt.Sprintf("right operand %d", givingOperand), func(t *testing.T) {
			p := 2
			lo := 1
			op := 1
			ro := givingOperand

			want := fmt.Sprintf("one + %d", givingOperand)

			cc := captcha.New(p, lo, op, ro)
			get := cc.String()
			assert.Equal(t, get, want, fmt.Sprintf("given %d %d %d %d want %q but get %q", p, lo, op, ro, want, get))
		})
	}

}
