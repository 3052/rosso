# Movistar TV App (Colombia)

## Platform
- URL: `mimovil.movistar.co` (activation page)
- SMS PIN-based activation flow
- Form POSTs to `/MovistarPlay_Movil/Home/NumberRegis`
- Protected by Google reCAPTCHA
- Note: Tigo acquired Movistar Colombia (Nov 2025) — service may be transitioning

## Activation Flow
- Step 1: Enter 10-digit cell phone number → receive SMS PIN
- Step 2: Enter PIN to confirm

## Attempts
- US number (8175015193): "El abonado o número celular no existe."
- 5SIM virtual66 Colombian number: "El abonado o número celular no existe."

## Key Finding
- The form validates against Movistar Colombia's **subscriber database**
- Only numbers registered as active Movistar Colombia accounts pass
- Virtual numbers, VoIP, and SMS relay services all fail
- Email-only registration does not exist

## Movie: ID 4994505
- Original URL: `https://movistarplay.co/details/movie/4994505`
- Country: CO (Colombia), monetization: FLATRATE
- US access: connection timed out (geo-blocked)
- Colombia access: redirects to `https://tv.movistar.co`

## Alternative Access Attempts
- ShareSub: fail
- G2G: fail
- Z2U: fail
