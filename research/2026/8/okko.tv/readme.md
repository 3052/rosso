# Summary: Why Okko TV is not gettable for a US-based user

## The goal
Sign up for okko.tv from the United States.

## What was attempted

## Path 1: Direct signup (US location, email flow)

1. Set location to US
2. Click "Subscribe"
3. Choose email signup
4. Enter email and continue
5. Enter the verification code sent to email
6. Register the account
7. Select "New map" (new card) at checkout
8. Enter card number, expiry (MM/YY), and CVV
9. Confirm

**Result:** Payment fails with `DECLINED_BY_PROCESSING` (error code 214). This
happens with uBlock Origin disabled, and regardless of whether the location is
set to US or RU.

```json
{
    "status": 0,
    "callId": "8b2804d5d60441d1b871feb649237491",
    "authorized": true,
    "serviceInfo": {
        "serverTime": 1786916655932
    },
    "userInfo": {
        "id": "480522185",
        "type": "UNKNOWN",
        "email": "tony117@web-library.net",
        "hasSbUserId": false
    },
    "paymentStatus": {
        "transactionStatus": "FAILED",
        "errorReason": "DECLINED_BY_PROCESSING",
        "errorReasonCode": 214
    }
}
```

## Path 2: Sber ID signup (Russian location, Sberbank auth flow)

1. Set location to Russia
2. Okko redirects to `id.sber.ru` (Sberbank ID OAuth)
3. Sber ID requires a Russian phone number

**Attempts to get a Russian phone number:**

- **virtualsim.net** — "There are no SMS for listed phone(s)" — even when
  sending from own phone
- **5sim.net** — Requires VPN to non-US/non-Russia location; unconfirmed if
  Sber ID works
- **sms-man.com** — Sberbank support unclear; no "Other" service option available
- **onlinesim.ru** — Explicitly blocks `sberid` SMS in their blocked senders list

**Additional blockers:**

- Sber ID explicitly supports Russian numbers only
- New Russian anti-spam law (August 1, 2025) blocks all A2P SMS for opted-out users
- Russian SIM cards now require biometric verification (SNILS + Gosuslugi +
  photo + voice) since July 2025

## Path 3: YooMoney (Russian MIR virtual card)

Attempted to register at yoomoney.ru. API response: `{"errorCode":
"RegistrationForbidden"}` — US phone numbers blocked. Same result from both US
and RU VPN locations.

## Path 4: Pay in Russia (Russian MIR virtual card)

- **Unverified card:** Only top-up methods visible are crypto (greyed out for
  lower-tier card), UnionPay, WeChat, and Alipay — none usable from the US
- **Verified card (crypto-enabled):** US passport not accepted — only German,
  Italian, Latvian, Estonian, Mexican, or Colombian passports

## Path 5: Other MIR virtual card providers

- **Vizovcc** — Scam: bait-and-switch tactics, demands undisclosed "activation fees"
- **Cardn3** — Scam: same bait-and-switch pattern as Vizovcc

## Path 6: Shared accounts / gift cards

- **sharesub.com** — Failed
- **g2g.com** — Failed
- **z2u.com** — Failed

## Root cause

Okko TV is owned by JSC "New Opportunities," a UK-designated sanctioned entity
(sanctioned June 29, 2022). The international financial infrastructure is
configured to prevent payments to Okko from outside Russia:

- International payment processors block non-Russian cards at checkout
- Sber ID (the primary auth method) requires a Russian phone number, which now
  requires biometric verification
- Russian virtual card providers either block US users, require passports from
  non-US countries, or are scams
- YooMoney (Sberbank subsidiary) blocks US phone numbers entirely

There is no known working path for a US-based user with a US passport and US
payment methods to sign up for and pay Okko TV.
