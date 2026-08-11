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

*Restored — the earlier "removed" assessment was based on RootMyTV v1/v2 (patched ~mid-2022), but the webosbrew scene has moved on.*

Newer exploits, notably **faultmanager**, cover webOS 4.0–9. webOS 4.0/4.5
(2018–2019 models) still has no patched firmware, and DejaVuln remains
unpatched on webOS 3.5 (2017 models) — only the newest firmware (webOS 5–10) is
mostly patched.

Before buying, ask the seller for the **"webOS TV Version"** from the settings
menu (not "Software Version" — they're different) and check it against
[cani.rootmy.tv](https://cani.rootmy.tv/). Target 2017–2019 models:

| Year | Model | webOS Version | Price |
|---|---|---|---|
| 2018 | LG OLED B8 | 4.0 | $400 |
| 2019 | LG UM7300 | 4.5 | $598.33 |
| 2019 | LG SM8600 | 4.5 | $799.99 |
| 2018 | LG UK6300 | 4.0 | $850 |
| 2018 | LG OLED C8 | 4.0 | $1449 |
| 2019 | LG OLED B9 | 4.5 | Out of stock online |
| 2019 | LG OLED C9 | 4.5 | Out of stock online |
| 2018 | LG UK6500 | 4.0 | Out of stock online |
| 2017 | LG OLED B7 | 3.5 | Out of stock online |
| 2017 | LG OLED C7 | 3.5 | Out of stock online |
| 2017 | LG UJ6300 | 3.5 | Out of stock online |
| 2017 | LG UJ7700 | 3.5 | Out of stock online |

On first boot, skip network setup (or block LG update servers at the router)
until rooted and the built-in update blocker is enabled. Do **not** install
LG's Developer Mode app — it interferes with rooting; the Homebrew Channel
replaces it.

Post-root, the Homebrew Channel auto-installs, SSH can be toggled on, and
`ares-cli` gives SDK access from a PC. The key advantage: webOS apps (including
Peacock) are HTML/JS packages — once rooted, the app's source is readable and
modifiable on-device, so EME/PlayReady calls can be observed directly, or a
minimal EME page can be written to probe the license server. This is the Edge
experiment, but on a platform where PlayReady is native and
hardware-provisioned.

**Step zero (before rooting):** DNS-log the TV (Pi-hole or router) while
playing Peacock content to confirm the app actually negotiates a **PlayReady**
license endpoint rather than Widevine.

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
