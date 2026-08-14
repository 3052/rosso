# Timeline of the Amazon Prime Video signup failure

**Day 1**
1. **Signed up for Amazon account** (email only, no phone) and **Prime Video + HBO Max channel add-on** — no error shown, movie played successfully.
2. On return, the site had **forgotten the Max subscription** and prompted to sign up again.
3. Investigation revealed **both Prime Video and HBO Max were canceled** on the account with no notice.
4. Received an email from Amazon: **order restrictions placed** on the account for "patterns differing from typical shopping behavior," with a 96-hour deadline to appeal.

**Appeal attempt (failed)**
5. Appeal page at `account-status.amazon.com` displayed only **"Switch
   countries to resolve issue"** — but the **"Switch countries" button was dead**
   (no redirect, no popup) in normal windows, private/incognito windows, and
   across devices.
6. Direct appeal form URL (`account-status.amazon.com/customer-appeal-form`) showed **"No Action Required"** — a dead end.
7. Country/Region setting in account preferences was found to be **"unset"**; setting it correctly and rejoining Prime resulted in **immediate auto-cancel of the subscription**.

**Second attempt (failed)**
8. Created a **new account with a new email and a new card**; Prime membership again **canceled immediately**.
9. Received a second email: account placed on hold for **"unusual payment activity"** pending payment-method validation.
10. Sign-in presented a verification page requiring **government ID or bank statement upload, with up to 24-hour review** before access is restored.

**Result:** After two accounts, two payment methods, and multiple
device/browser combinations, Amazon's automated fraud system blocked a
legitimate new customer from using a paid service at every step, with the only
appeal mechanism broken and the only remaining remedy requiring a 24-hour
manual document review.
