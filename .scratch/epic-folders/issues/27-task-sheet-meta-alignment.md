# 27: Task sheet metadata rows are misaligned

**What to build:** In the task sheet's metadata list (Stage, Owner, Labels, Added), every value starts at the same left edge and sits on the same baseline as its label.

**Blocked by:** none

**Status:** ready-for-agent

Seen on 2026-10-29 (screenshot of a Backlog task with no owner and the `epic engine-fixes` label):

- **Owner:** the add-owner icon button is indented to the right of the value column, while Backlog, the label chip and the date start further left. Likely the button's own padding.
- **Labels:** the chip is taller than a text row, so its text sits lower than the "Labels" label and the row spacing looks uneven.
- **Stage:** "Backlog" and the "Move to" button are fine horizontally, but check the button's padding doesn't push the row off the label's baseline.

Code: `client/src/components/organisms/TaskSheet.tsx`, `molecules/OwnerPicker.tsx`, `molecules/StagePicker.tsx`.

- [ ] All values in the metadata list share one left edge (icon buttons included)
- [ ] Each value is vertically centred on, or baseline-aligned with, its label
- [ ] Rows keep even spacing with and without an owner and labels
- [ ] Checked in the running app in light and dark mode (by hand)
