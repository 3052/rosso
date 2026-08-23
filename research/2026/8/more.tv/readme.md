# Wink (Russia) — Registration Attempts

## Platform
- Original URL: `https://more.tv/devushka-s-tatuirovkoy-drakona`
- Redirected to: `https://wink.ru/`
- More.tv merged into Wink (Rostelecom's streaming service)

## Registration Requirements
- Russian +7 phone number required (website, Android app, and Android TV all phone-only)
- SMS 4-digit code verification
- Email registration: **removed** (was available in older app versions, no longer supported per current Dec 2025 FAQ)

## Attempts
- Website login form: phone numbers only (cannot type letters)
- Android app (`ru.rt.video.app.mobile`): also phone-only
- Kazakhstan +7 number (Grizzly SMS): "You can not send an SMS to the specified number" — Wink checks country code, not just +7 prefix

## SMS Provider Attempts
- **5SIM:** Russia blocked (country not found error)
- **SMS-Activate:** shut down December 29, 2025
- **Onlinesim:** no Wink service listed, no "Other" option available
- **1001SMS:** no Russia numbers available
- **SMS-MAN:** no Russia numbers available
- **GetSMSCode:** explicitly lists Wink for Russia (8,953 numbers in stock at
  $0.20) — $15 deposited via USDT ERC20 but payment page timed out, balance
  still $0, support emailed with no response
- **Grizzly SMS:** Kazakhstan +7 number rejected by Wink (country code check)
- **TurboN.Rent:** no Russia numbers available for any service
- **Luchibb:** requires TRC20 (Tron) USDT, Coinbase doesn't support TRC20 withdrawals, bridging fees would eat most of the $8.53 USDT balance

## Conclusion
- Wink requires a Russian +7 phone number for all registration paths
- No email registration path exists on any platform (web, Android, Android TV)
- Western SMS providers have largely pulled Russian numbers due to sanctions
- Russian-based providers (Onlinesim) don't list Wink as a service
- GetSMSCode is the only confirmed provider with Wink for Russia, but payment credit issue remains unresolved
- No working registration path identified from US without resolving the GetSMSCode balance issue
