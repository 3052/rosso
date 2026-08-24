# TV+ (Turkey)

## Platform
- URL: `tvplus.com.tr/giris`
- Requires Turkish mobile number (+90) whose local part starts with 5
- Registration is SMS-verified
- No email, Apple, or Google signup available on the website

## Full Attempt Log

**The core constraint:** TV+ (`tvplus.com.tr/giris`) requires a Turkish mobile
number (+90) whose local part starts with 5 (that's the "must start with 5"
error). Registration is SMS-verified. There is no email, Apple, or Google
signup on the website. So the entire effort was about obtaining a Turkish +90
number that could actually *receive* an SMS while you're in the US.

### Phase 1 — Finding a platform that even has Turkey

**1. VirtualSIM (virtualsim.net)** — your preferred service, you'd used it
before for Russia. I checked their live ordering page and availability API.
They carry no Turkey on any of their six plans (EXTENDED24 $20/mo, AUTOMATIC
$30/mo, DEDICATED $36+/yr, ORDINAL $6+/yr, BUSINESS $3/mo, ONETIME $2). Their
country sets are Eastern Europe / CIS / Western Europe specialists (AT, CZ, DE,
DK, EE, FI, GE, IE, KG, LV, MA, MT, NL, PL, RU, UA, UK, USA). Their service
list had a few Turkish brands (papara.com, Paycell, BiP.com, Ozan.com) but no
Turkcell/TV+ entry, and no Turkish numbers to back them. **Verdict: no
Turkey.** No money spent.

**2. SMS-Activate (sms-activate.org)** — usually the top pick for this kind of
task. Checked before recommending: the platform shut down in December 2025; all
`sms-activate.*` domains went dark. **Verdict: defunct.** No money spent.

**3. SMSPool (smspool.net)** — my first recommendation. Verified: Turkey was in
   their country picker, they advertise real non-VoIP SIM numbers, accept
   credit/debit cards, and only charge when a code actually arrives (failed
   attempts free). You created an account and deposited **$5 by card**. On the
   Order page, both pool options failed:
- `pool mike` → "no numbers available at the moment please try again later"
- `pool charlie` → "your current selection is out of stock - please try again later"
**Verdict: Turkey out of stock.** $5 deposited, unused, refundable within 14 days.

**4. SMS-Man (sms-man.com)** — second recommendation, chosen because
independent testing flagged it as good for Turkey and it takes cards. You
signed up and tried to request a Turkey number. Their response: *"You are
requesting a very rare and popular number... According to statistics, the
number to such services is given with 5–6 attempts... we have added this
feature to telegram bot! Your API KEY: 9ypTwDBaQfh0e6LifS2xQItIUaEtogBC."* I
warned you to treat that API key like a password and not share it. I laid out a
manual-retry route (expect 5–6 attempts) vs. their Telegram bot. Your second
try returned: **"no numbers try again later."** **Verdict: stock depleted.**
Small balance unused, refundable.

**5. GrizzlySMS (grizzlysms.com)** — third recommendation. Verified they had a
dedicated Turkey section and a Turkcell service page, accept cards, auto-refund
if no code in 20 min. Their page showed "few" numbers for Turkcell+Turkey
(their FAQ says "few" = 0–2 in stock, fluctuates, may need to press GET 5–15
times). When you tried on the website, both **Turkcell** and **AnyOther** said
**"available only by API."** So we pivoted to their API, driven from your
browser address bar (sms-activate-compatible). This became the longest attempt
— see Phase 2.

### Phase 2 — GrizzlySMS API deep-dive

You got an API key from your account settings. We tested it with `getBalance` (free) — it returned `ACCESS_BALANCE` and worked.

- **Wrong-country incident:** I initially guessed Turkey's country ID was 63. A
  purchase with `country=63` returned `ACCESS_NUMBER:578974278:420722690309` —
  a number starting **420**, which is the Czech Republic, not Turkey (90). I
  told you not to use it. Cancelling immediately returned
  **`EARLY_CANCEL_DENIED`** (their anti-abuse rule: you can't cancel within the
  first 2 minutes). After waiting 2 minutes, the same cancel call returned
  **`ACCESS_CANCEL`** and the $0.90 was refunded.
- **Country list lookup:** `getCountries` confirmed **Turkey = ID 62**, **Czech = ID 63**.
- **Stock check (free):** `getPrices&country=62` showed the `ot` ("Any other") service with **count 4,094 at $0.59**, and `gr_tk` with 4,865 at $0.86.
- **Purchase attempts:** every `getNumber&service=ot&country=62` returned
  **`NO_NUMBERS`**. Retried after waiting; still `NO_NUMBERS`. Tried
  `service=gr_tk`; still `NO_NUMBERS`. Tried the V2 endpoint with `maxPrice=3`;
  still `NO_NUMBERS`. Tried with a **Turkish VPN server** active; still
  `NO_NUMBERS`.

**Verdict:** the displayed stock counters were stale estimates, not live inventory — the live pool was empty on every attempt. **Balance $7.96, nothing spent, refundable.**

### Phase 3 — Eliminating the rest of the market

**6. Quackr** — their FAQ lists 24 supported countries; Turkey isn't among them. **Verdict: no Turkey.** No money spent.

**7. NexSMS (nexsms.net)** — has real Turkish *carrier* numbers (actual
Turkcell/Vodafone/Türk Telekom SIMs, not VoIP), with premium long-term numbers
valid 12–90 days. I initially dismissed it because payment is
crypto/USDT/Alipay/WeChat only — but you correctly pointed out you never said
cards were required. Later you signed up and saw Turkey listed at **$15.98 for
a regular number**, with "premium" and "boutique" tiers (premium had no Turkey
country option, only "projects" like WhatsApp/Telegram; boutique only had +1,
+44, +62 numbers). When we circled back to it later, the price had surged to
**$65.2281** — the same auction/surge pattern as SMSPin. **Verdict: unusable
price.** No money spent.

**8. TextVerification** — their 43 "online" Turkey numbers are free public
numbers anyone can read; unsafe for a paid account (someone could later take it
over). **Verdict: skipped as unsafe.** No money spent.

**9. SMSCodex (smscodex.com)** — a marketplace aggregating multiple number
sellers. Their public Turkey page listed **200 services, 544 numbers, from
$0.07/SMS**, including **Turkcell: 4 numbers at $0.12**. You created an
account. When you tried to pay, it said: **"foreign cards not accepted — use
cards issued by russian banks."** **Verdict: Russian cards only.** No money
spent.

**10. DialAnyone (dialanyone.com)** — markets private rented +90 numbers,
claims regular cards accepted. You signed up and looked for Turkey under
Numbers → Browse & Purchase: **not in the country list.** Clicking Messages →
Browse & Purchase looped back to the same page. I then found their marketing
page `dialanyone.com/receive-sms/turkey`, but its FAQ only said "a small
monthly fee... sign up is free, no credit card required to look at
availability" — no real inventory. **Verdict: the Turkey page is SEO marketing,
no real Turkish numbers.** No money spent.

**11. SmsPva** — checked live; Turkey not in their country list. **Verdict: no Turkey.** No money spent.

**12. OnlineSim** — their own page states their country selection is
"extensive, **except Turkey**." Free +90 numbers down; paid pool nonexistent.
**Verdict: no Turkey.** No money spent.

**13. SMSZ** — Turkey page exists but their catalog showed exactly one Turkey
combination (Telegram only), no Turkcell/TV+. SMSZ and 1001SMS appear to be the
same company (cross-linked). **Verdict: no usable Turkey service.** No money
spent.

**14. Turkcell official eSIM** — I checked whether you could buy a legitimate
Turkish number directly from Turkcell. They now sell tourist eSIMs online to
foreigners (passport scan + live video verification via their app), but the
packages **only activate when you arrive in Turkey/Europe** — won't work
roaming in the US. **Verdict: not viable from the US.** No money spent.

### Phase 4 — The two platforms that got closest

**15. SMSPin (smspin.io)** — a dedicated Turkey page, crypto payment
(Cryptomus: USDT/BTC/ETH), auto-refund if no code in ~20 min. Their live
pricing page showed **Turkey +90: ~1.6 million numbers in stock, from
$1.18/SMS**, non-VoIP (real Android devices with real SIMs in Turkey). You
created an account and funded **$5**. On the Get Number page (Turkey + "Any
other"), the price kept escalating as you clicked buy: **$0.20 → "No numbers
available for this operator right now. Try another." → $0.56 → "This deal just
sold out. Refresh to see updated pricing." → $0.72 → $1.18 → $4 →
"unavailable."** The Operator menu listed operator 6, 7, 10, 4, 8 — all sold
out. The Rent option required choosing a service (no "Any other"). **Verdict:
pool empty, auction bidding up a phantom number.** ~$5 deposited, unused,
refundable.

**16. 1001SMS (1001sms.com)** — verified Turkey + "Other" service with 120+
   numbers from $0.95, card + crypto payment, refund if no SMS. You signed up.
   Card top-up failed: **"Cards issued in the United States, Russia, Kazakhstan,
   and Turkey are not supported by this payment method."** I confirmed they accept
   crypto (BTC, ETH, USDT TRC20/ERC20, TRX, BNB; $3 min; crypto deposits
   non-refundable but failed numbers refund to balance). You funded **$10.15 in
   crypto**. On the Turkey page, "Other" showed **$36.28** (gouged catch-all
   price, 2,136 numbers). I told you not to buy it and to pick the cheapest named
   service instead. You pasted the service dropdown; cheapest were **Discord $0.28
   (7,418 available)** and **Baidu $0.31 (44,464 available)**.
- Clicking "Receive SMS" on Discord revealed three operator pools: `mike` success 39% $1.21, `charlie` success 39% $1.21, `default` no success rate $0.28.
- I told you to pick **mike** ($1.21); it returned **"choose another operator"** (mike empty).
- You switched to **charlie** ($1.21) — **it issued a number.**
- You entered it on TV+ (Turkey +90, 10 digits starting with 5) and requested
  the SMS. After ~1 minute no code; the 1001SMS page showed "Still no OTP?
  Cancel this number and get a new one. You'll receive a full refund
  automatically if no SMS was received." I said wait 5 more minutes then
  Resend, don't cancel early.
- After 5 min + a Resend + another minute, **no code arrived.** You suspected
  1001SMS was blocking it because you'd bought under "Discord," not TV+. I
  suggested looking for an "All messages" filter; you found none.

**Verdict:** charlie's pool issued numbers but no TV+ code arrived — likely
service filtering (dashboard only shows codes for the purchased service) or a
dead number at the 39% success rate. **$10.15 crypto unused, refundable via
support.** (Note: "mike"/"charlie" are the same pool names SMSPool showed on
day one — these platforms share backend providers, which is why Turkey was
empty everywhere at once.)

**17. HotTelecom (hottelecom.biz)** — verified a dedicated +90 Turkey SMS
number: **$50 setup + $60/month**, crypto accepted, connected within 24 hours,
assigned to you (no lockout risk), SMS-only with all texts visible in the
dashboard. Fixed price, no surge. You rejected it as too expensive (~$110
total). **Verdict: real but too costly.** No money spent.

**18. PVAPins (pvapins.com)** — a dedicated Turkey/crypto page, one-time
numbers and rentals, different operator pool. You confirmed **we'd already
tried this one** earlier. **Verdict: already attempted.**

## Alternative Access Attempts
- ShareSub: fail
- G2G: fail
- Z2U: fail
