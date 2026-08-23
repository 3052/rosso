# Movistar Play / Movistar TV (Argentina)

## Platform
- Registration URL: `https://idp.movistar.com.ar/self-registration/?returnTo=https%3A%2F%2Ftv.movistar.com.ar%2F`
- Streaming URL: `https://tv.movistar.com.ar/`
- JSP-based multi-step form (Step 1 POSTs to `step-1.jsp`)

## Step 1 Required Fields
- Document type: DNI or Pasaporte (Passport accepted)
- Document number: 8 digits max (DNI), numbers only
- Número de trámite: 11 digits max (transaction number from Argentine DNI)
- Gender: Femenino, Masculino, or No Binario
- Date of birth: hidden field (default `1960-01-01`)

## Key Findings
- Passport is accepted as document type (unlike Chile's RUT-only)
- Step 1 doesn't ask for phone number (unlike Colombia's SMS flow)
- Número de trámite required even for passport holders — this is an Argentine-specific ID field
- Hidden date of birth field suggests validation against government records

## Attempts
- Passport option reviewed but not attempted — número de trámite requirement and likely Argentine payment method in later steps make it impractical from US

## Alternative Access Attempts
- ShareSub: fail
- G2G: fail
- Z2U: fail
