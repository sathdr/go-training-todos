package captcha

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Captcha represents captcha
type Captcha struct {
	pattern      int
	leftOperand  int
	operator     int
	rightOperand int
}

// New create new captcha
func New(p, lo, op, ro int) Captcha {
	return Captcha{p, lo, op, ro}
}

var numbers = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
var operators = []string{"", "+", "-", "*"}

func (cc Captcha) String() string {

	if cc.pattern == 1 {
		return fmt.Sprintf("%d %s %s", cc.leftOperand, operators[cc.operator], numbers[cc.rightOperand])
	} else if cc.pattern == 2 {
		return fmt.Sprintf("%s %s %d", numbers[cc.leftOperand], operators[cc.operator], cc.rightOperand)
	}

	return ""
}

var src = rand.NewSource(time.Now().UnixNano())
var rnd = rand.New(src)
var store = map[string]int{}

// KeyQuestion return a new captcha string
func KeyQuestion() (string, string) {
	pattern, leftOperand, operator, rightOperand := rnd.Intn(1)+1, rnd.Intn(9)+1, rnd.Intn(3)+1, rnd.Intn(9)+1
	answer := 0
	switch operator {
	case 1:
		answer = leftOperand + rightOperand
	case 2:
		answer = leftOperand - rightOperand
	case 3:
		answer = leftOperand * rightOperand
	}

	cc := New(pattern, leftOperand, operator, rightOperand)
	key := uuid.New().String()
	store[key] = answer
	return key, cc.String()
}

var mux sync.Mutex

// Answer check answer
func Answer(key string, ans int) bool {
	// prevent race condition
	mux.Lock()
	defer mux.Unlock()

	if v, ok := store[key]; ok {
		delete(store, key)
		return v == ans
	}
	log.Printf("not found %s key in store\n", key)
	return false
}
