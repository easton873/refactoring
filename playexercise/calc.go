package playexercise

import (
	"fmt"
	"math"
)

type Invoice struct {
	customer     string
	performances []PlayInfo
}

type PlayInfo struct {
	playID   string
	audience int
}

type Play struct {
	name string
	kind string
}

func statement(invoice Invoice, plays map[string]Play) (string, error) {
	totalAmount := 0.0
	volumeCredits := 0
	result := fmt.Sprintf("Statement for %s\n", invoice.customer)
	format := func(amount float64) string {
		return fmt.Sprintf("$%.2f", amount)
	}
	for _, perf := range invoice.performances {
		play := plays[perf.playID]
		thisAmount := 0.0

		switch play.kind {
		case "tragedy":
			thisAmount = 40000
			if perf.audience > 30 {
				thisAmount += 1000 * (float64(perf.audience) - 30)
			}
		case "comedy":
			thisAmount = 30000
			if perf.audience > 20 {
				thisAmount += 10000 + 500*(float64(perf.audience)-20)
			}
			thisAmount += 300 * float64(perf.audience)
		default:
			return "", fmt.Errorf("unknown type: %s", play.kind)
		}

		// add volume credits
		volumeCredits += int(math.Max(float64(perf.audience-30), 0))
		// add extra credit for every ten comedy attendees
		if "comedy" == play.kind {
			volumeCredits += int(math.Floor(float64(perf.audience / 5)))
		}

		// print line for this order
		result += fmt.Sprintf(" %s: %s (%d seats)\n", play.name, format(thisAmount/100), perf.audience)
		totalAmount += thisAmount
	}
	result += fmt.Sprintf("Amount owed is %s\n", format(totalAmount/100))
	result += fmt.Sprintf("You earned %d credits\n", volumeCredits)
	return result, nil
}
