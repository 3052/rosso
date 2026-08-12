# peacock

https://peacocktv.com/help/article/what-devices-and-platforms-are-supported-by-peacock

I've been trying to get Peacock H.265 working. It seems that H.265 only works
with UHD content, and I couldn't find any public Peacock scripts that support
UHD.

## Microsoft Edge

MS Edge 125+ (Windows only — PlayReady is not available in Edge on Mac).

*Note: Edited the Edge client to trigger a PlayReady request, but the server rejected it.*

## Xbox One

*Note: Retail mode can't run Fiddler; dev mode App Store returns no search
results, so retail apps (including Peacock) won't work. Sideloading the Peacock
app directly fails — it's an encrypted `.eappxbundle` that can't be installed
in dev mode.*

## LG webOS — Primary Target

| Year | Model | webOS Version | Price |
|---|---|---|---|
| 2018 | LG OLED B8 | 4.0 | $400 |
| 2019 | LG UM7300 | 4.5 | $598.33 |
| 2019 | LG SM8600 | 4.5 | $799.99 |
| 2018 | LG UK6300 | 4.0 | $850 |
| 2018 | LG OLED C8 | 4.0 | $1449 |
| 2017 | LG OLED B7 | 3.5 | Out of stock online |
| 2019 | LG OLED B9 | 4.5 | Out of stock online |
| 2017 | LG OLED C7 | 3.5 | Out of stock online |
| 2019 | LG OLED C9 | 4.5 | Out of stock online |
| 2017 | LG SJ8000 * | 3.5 | Out of stock online |
| 2017 | LG UJ6300 | 3.5 | Out of stock online |
| 2017 | LG UJ6540 * | 3.5 | Out of stock online |
| 2017 | LG UJ7700 | 3.5 | Out of stock online |
| 2018 | LG UK6500 | 4.0 | Out of stock online |
| 2018 | LG UK7700 * | 4.0 | Out of stock online |
| 2019 | LG UM7400 * | 4.5 | Out of stock online |

`*` = rows I added. All models are rootable on any firmware (LG never patched
webOS 3.5–4.5).

## Samsung Smart TV — Second Choice

Models from 2017 or later. Rootable via SamyGO, but it's a per-model/per-firmware slog with forum-era tooling and no pre-purchase compatibility checker. Requires blocking updates.

## VIZIO SmartCast TV — Third Choice

SmartCast TV (2016 and newer). Rootable via ViziOwn (clean pre-auth remote
root) on unpatched firmware, but the app is a cloud-driven thin client — little
to inspect on-device. Target older SmartCast-branded models (~2016–2021:
D/V/M/P-Series) bought used and kept offline during setup.

## Hisense VIDAA

*Removed: no known root path.*

## Roku

*Removed: locked bootloader, no public root.*

## Xbox Series X/S

*Removed: no exploit exists.*

## Xumo / Spectrum Xumo Stream Box

*Removed: no known root; Comcast-leased hardware.*

## PlayStation

*Removed: rootable only on old firmware (PS4 ≤ 12.02, PS5 ≤ 5.50), but Peacock requires PSN access on current firmware, so a jailbroken console can't actually run the app.*
