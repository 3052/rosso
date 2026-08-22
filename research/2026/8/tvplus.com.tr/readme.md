# https://tvplus.com.tr/film-izle/asiricilar--204922332

## TV+ Signup from US — Attempt Log

Goal: create a TV+ account from the US. Requires a Turkish mobile number (+90
5xx) with SMS verification. No email signup exists.

1. **VirtualSIM** — no Turkey on any plan (checked their live API).
2. **SMS-Activate** — platform shut down December 2025.
3. **SMSPool** — Turkey out of stock. $5 deposited, unused.
4. **SMS-Man** — "no numbers, try later." Small balance unused.
5. **GrizzlySMS (via API)** — stock counters showed thousands of Turkey
   numbers, but every buy returned NO_NUMBERS. Tried `ot`, `gr_tk`, the V2
   endpoint, and a Turkish VPN. $7.96 unused.
6. **Quackr** — doesn't support Turkey.
7. **NexSMS** — crypto/Alipay only, no cards.
8. **TextVerification** — shared public numbers only, unsafe.
9. **SMSCodex** — live Turkcell stock found, but Russian bank cards only.
10. **DialAnyone** — Turkey page is SEO marketing. No Turkey in the actual country list.
11. **Turkcell official eSIM** — only activates inside Turkey, not the US.

Wrong-number incident: accidentally bought a Czech number on GrizzlySMS
(guessed the country ID). Hit a 2-minute cancel lock (EARLY_CANCEL_DENIED),
then got the refund (ACCESS_CANCEL, $0.90 back).

Key facts learned:
- GrizzlySMS API: Turkey = country ID 62, Czech = 63, cancel = setStatus&status=8
- Stock counters are stale estimates, not live inventory
- Turkish numbers are currently the hardest to obtain across the whole market

Unused balances recoverable: ~$7.96 on GrizzlySMS, ~$5 on SMSPool, small amount
on SMS-Man. Contact their support for refunds.

Conclusion: no working automated route found.
