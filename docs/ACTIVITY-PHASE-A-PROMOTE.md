# Activity Phase A — promote with 2.0

Anonymous visitor Activity (public `POST /activity`, `surf_aid` / `surf_sid`) is implemented in surfing-service and the home FE tracker.

When tagging **dasmlab-home 2.0.0** for production:

1. Confirm surfing-service image in prod already includes Phase A handlers (verified earlier on shared API).
2. Roll home FE so production serves the always-on tracker + history-mode routes + hubs.
3. Owner Activity panel remains dual-gated (`admin` + `ACTIVITY_VIEWERS`).

No separate Feature flag required — demote is not needed for Phase A.
