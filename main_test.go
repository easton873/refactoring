package refactoring

import "testing"

func Test_statement(t *testing.T) {
	type args struct {
		invoice Invoice
		plays   map[string]Play
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Empty", args{}, `Statement for 
Amount owed is $0.00
You earned 0 credits
`, false,
		},
		{
			"Book example", args{
				invoice: Invoice{
					customer: "BigCo",
					performances: []PlayInfo{
						{playID: "hamlet", audience: 55},
						{playID: "as-like", audience: 35},
						{playID: "othello", audience: 40},
					},
				},
				plays: map[string]Play{
					"hamlet":  {name: "Hamlet", kind: "tragedy"},
					"as-like": {name: "As You Like It", kind: "comedy"},
					"othello": {name: "Othello", kind: "tragedy"},
				},
			},
			`Statement for BigCo
 Hamlet: $650.00 (55 seats)
 As You Like It: $580.00 (35 seats)
 Othello: $500.00 (40 seats)
Amount owed is $1730.00
You earned 47 credits
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := statement(tt.args.invoice, tt.args.plays)
			if (err != nil) != tt.wantErr {
				t.Errorf("statement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("statement() got = %v, want %v", got, tt.want)
			}
		})
	}
}
