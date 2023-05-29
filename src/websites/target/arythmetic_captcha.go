package target

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func SolveArythmeticCaptcha(doc *goquery.Document) (answer string) {
	captcha := doc.Find("#fn_fast_order").Find("div.secret_number").Text()
	if captcha == "" {
		log.Info("No arithmetic captcha found, skipping captcha answer")
		return ""
	}

	log.Info("Arithmetic captcha found, solving: " + captcha)
	return SolveExpression(captcha)
}

func SolveExpression(captcha string) (answer string) {
	//assuming captcha consists of numbers, +, -, =, and ? on either side of the expression
	sides := strings.Split(captcha, "=")

	splitFn := func(c rune) bool {
		return c == ' '
	}
	leftSide, rightSide := strings.FieldsFunc(sides[0], splitFn), strings.FieldsFunc(sides[1], splitFn)
	diff := solveExpression(rightSide) - solveExpression(leftSide)
	diff *= getDiffSign(leftSide, rightSide)

	log.Info("Captcha solution: ", diff)
	return strconv.Itoa(diff)
}

func solveExpression(e []string) int {
	result := 0

	for i := 0; i < len(e); i += 2 {
		if val, err := strconv.Atoi(e[i]); err == nil {
			if i == 0 || e[i-1] == "+" {
				result += val
			} else if e[i-1] == "-" {
				result -= val
			} else {
				panic("Unknown operation: " + e[i-1])
			}
		}
	}
	return result
}

func isVariable(e string) bool {
	_, err := strconv.Atoi(e)
	return err != nil && e != "+" && e != "-" && e != "="
}

func getDiffSign(leftSide, rightSide []string) int {
	for i, e := range leftSide {
		if isVariable(e) {
			if i == 0 || leftSide[i-1] == "+" {
				return 1
			} else {
				return -1
			}
		}
	}

	for i, e := range rightSide {
		if isVariable(e) {
			if i == 0 || rightSide[i-1] == "+" {
				return -1
			} else {
				return 1
			}
		}
	}
	panic("Variable not found in the expression")
}
