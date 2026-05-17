package main

type provider struct {
   actions []any
   result  string
}

var _ = map[string]provider{
   "privacy.com": {
      actions: []any{
         0: "mubi.com/memberships",
         1: "start your membership",
         2: "mubi",
         3: "next",
         4: map[string]any{
            "email address": map[string]any{
               "mail.tm":       "sorry this isn't a valid email",
               "mailsac.com":   "sorry this isn't a valid email",
               "tempmail.best": "valid for 1 hour",
            },
         },
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
   },
   "wise.com": {
      result: "works",
   },
   "bankofamerica.com": {
      actions: []any{
         0: "mubi.com",
         1: map[string]any{
            "email address": map[string]any{
               "tempmail.best": nil,
            },
         },
         2:  "get started",
         3:  "mubi",
         4:  "next",
         5:  "cardholder name",
         6:  "card number",
         7:  "expiry date",
         8:  "CVV",
         9:  "zip code",
         10: "start free trial",
      },
      result: "Something’s not right. Please check your details and try again",
   },
   "paypal.com": {
      actions: []any{
         0: map[string]any{
            "mubi.com/memberships": map[string]any{
               "7 day free trial": nil,
            },
         },
         1: "start your membership",
         2: "mubi",
         3: "next",
         4: map[string]any{
            "email address": map[string]any{
               "tempmail.best": nil,
               "mail.tm":       "sorry this isn't a valid email",
               "mailsac.com":   "sorry this isn't a valid email",
            },
         },
         5: "sign up",
         6: "agree membership terms",
         7: "pay with payPal",
         8: "agree and continue",
      },
      result: "This credit card is already associated with an existing account",
   },
}
