package mubi

var _ = map[string]provider{
   "privacy.com": {
      actions: []any{
         0:  "mubi.com/memberships",
         1:  "start your membership",
         2:  "mubi",
         3:  "next",
         4:  "email address",
         5:  "sign up",
         6:  "enter code",
         7:  "cardholder name",
         8:  "card number (MUST USE NEW CARD EACH TIME)",
         9:  "expiry date",
         10: "CVV",
         11: "zip code",
         12: "I agree that after my first 7-days free, my membership will automatically renew",
         13: "start free trial",
      },
      result: "Please check your card has sufficient funds to complete the purchase",
      date:   "2026-05-17",
   },
   "paypal.com": {
      actions: []any{
         0: "mubi.com/memberships",
         1: "start your membership",
         2: "mubi",
         3: "next",
         4: "email address",
         5: "sign up",
         6: "agree membership terms",
         7: "pay with payPal",
         8: "agree and continue",
      },
      result: "This credit card is already associated with an existing account",
      date:   "2026-05-17",
   },
   "bankofamerica.com": {
      actions: []any{
         0:  "mubi.com",
         1:  "ENABLE JAVASCRIPT",
         2:  "email address",
         3:  "get started",
         4:  "mubi",
         5:  "next",
         6:  "cardholder name",
         7:  "card number",
         8:  "expiry date",
         9:  "CVV",
         10: "zip code",
         11: "start free trial",
         12: "mubi.com/subscription",
         13: "cancel mubi subscription",
         14: "skip to cancel",
         15: "cancel my subscription",
      },
      result: "If you cancel now, your trial will end immediately",
      date:   "2026-05-17",
   },
   "wise.com": {
      actions: []any{
         0: "mubi.com/memberships",
         1: "start your membership",
         2: "mubi",
         3: "next",
         4: map[string]any{
            "email address": map[string]string{
               "tempmail.best": "valid for 1 hour",
               "mail.tm":       "sorry this isn't a valid email",
               "mailsac.com":   "sorry this isn't a valid email",
            },
         },
         5:  "sign up",
         6:  "enter code",
         7:  "cardholder name",
         8:  "card number",
         9:  "expiry date",
         10: "CVV",
         11: "zip code",
         12: "I agree that after my first 7-days free, my membership will automatically renew",
         13: "start free trial",
      },
      result: "pass",
      date:   "2026-05-17",
   },
}

type provider struct {
   actions []any
   result  string
   date    string
}
