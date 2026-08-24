# Strim (Norway)

## Platform
- Registration URL: `https://www.strim.no/bli-kunde/registrer?productIds=9pvit`
- Strim is a Norwegian streaming service

## Registration Form
- Phone number field: Norwegian mobile numbers only, 8 digits
- Format check only — does not verify if the number is a real Strim subscriber
- Test number `41234567` (8 digits, starts with 4) was accepted

## Payment Methods Required
- Vipps (Norway's mobile payment app)
- PayPal (must be Norwegian PayPal account)
- Norwegian Visa/Mastercard (must be Norwegian-issued)

## Attempts
- Phone number `41234567`: accepted (format check only)
- Norwegian card: rejected — "To become a Strim subscriber, you must have a Norwegian card"
- PayPal: rejected — "To become a Strim subscriber, you must have a Norwegian PayPal account"
- Vipps: not attempted — requires Norwegian national ID (fødselsnummer/D-number) + Norwegian phone + Norwegian bank account

## Key Findings
- Phone field only checks format, not subscriber status — easily bypassed
- All three payment methods are locked to Norwegian identity/banking:
  - Card: requires Norwegian bank account
  - PayPal: requires Norwegian PayPal account (which itself requires Norwegian national ID, Norwegian bank/card, Norwegian proof of identity and address — PayPal shares data with Norwegian tax authorities via CRS/FATCA)
  - Vipps: requires Norwegian national ID + phone + bank
- BankID verification also required for trial/discount packages — BankID requires Norwegian national identity number (fødselsnummer) or D-number, impossible from US
- Norwegian PayPal account cannot be created from the US

## Alternative Access Attempts
- ShareSub: fail
- G2G: fail
- Z2U: fail

## Conclusion
- Strim blocks US users at the payment step, not registration
- All payment paths require Norwegian residency and identity
- No working registration path from US without Norwegian identity and banking
