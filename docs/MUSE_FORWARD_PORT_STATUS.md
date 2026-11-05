# Official TECHO5 v1.6.26 → Hey Muse forward-port

Updated: 2026-11-06 JST

## Objective
Build an independent `cronos` firmware from the **official stable TECHO5 code**, preserving Muse BLE pairing, wake-word, assistant turns and avatar, then add **Japanese spoken-input and Japanese audible-output**. Engine/format choices remain open; `2357gi/hey-muse` is reference-only.

## Verified source baselines
| Component | Pinned commit |
| --- | --- |
| Official TECHO5 stable `v1.6.26` | `82722e87fdd20f5e4e92ed6d3eadad55bcd56d6d` |
| Original Hey Muse `muse` | `1b574cca0a49cb4aa36a69600663434fd2d313fa` |
| Common ancestor | `8a4d96cdec7d76c8c80b1e203093b4efc214ec02` |

The integration branch was created *directly from the official stable commit*, not from the Muse fork. The original `muse` branch is untouched and protected. Initial diff inspection found 174 Muse-changed files, of which 34 were also changed upstream.

## Implemented on this integration branch
- Ported 122 isolated Muse-only source/assets files from original fork: Muse SDK, pairing, speech client, avatar, wakeword models and test suites.
- Added Muse configuration and protected SDK/TTS credentials without replacing official `News`, `Reader` and `ReplyPlayer` configuration.
- Added Muse to backend selection with realtime-slot override retained, turn-buffer isolation, wake readiness and optional avatar render paths.
- Integrated Muse pairing and audio endpoint controls into current official Setup.
- Ported BlueZ **Peripheral** alongside official BlueZ Central instead of overwriting its GATT implementation.
- Patched settings export to **exclude** Muse SDK and speech endpoint secrets and preserve them only for the same matching endpoint on restore.
- Kept optional old `Usage` screen behind `muse_usage_experimental` build tag; this is not part of Japanese speech baseline.

## Compile and test status
- First manually gated GitHub Actions run `42076036917`: **failed**, missing BlueZ Peripheral and two unredacted settings export secrets; both addressed.
- Second run `42076691591`: **failed** because of old Muse `utterance` dependency and an optional Usage screen package, fixed in subsequent changes.
- Third run `42077501584`: failed on two original Muse Setup assertions (English built-in voice warning and pairing status JSON), fixed.
- Fourth run `42078299872`: **success** — Go unit tests and ARMv7 cross-build.
- Fifth run `42079152900`: **success** — after adding swappable Wyoming STT → Muse text requests, fail-closed audio handling and offline unit tests. These are source/compile validations only; no device or end-to-end speech verification.
- One manual-gated, unsigned preview rootfs build `42079486348` triggered after green CI; do **not** install its artifact without independent hardware safety gate and approval.
- Workflow `.github/workflows/port-ci.yml` performs Go tests and ARMv7 compile only; no signing, auto-installation or artifact upload. Trigger manually or by editing `.ci/run-integration`.

## Outstanding integration
- [x] All `go test ./...` and ARMv7 daemon compile green on pinned v1.6.26 (run `42079152900`, Go 1.27 module).
- [ ] Preserve upstream display features and test Muse avatar during a turn.
- [ ] Verify local Wake word paths and Muse pairing client on actual `cronos`.
- [x] Optional Wyoming STT text path implemented and tested offline (including no audio upload after local STT failure).
- [ ] Japanese STT model deployment and real-device benchmark (Moonshine Tiny/Base, independent alternatives) pending.
- [ ] Japanese TTS: set a trusted Japanese-capable speech endpoint first. Built-in Pico is **English-only**, cannot serve as Japanese fallback. Investigate local Open JTalk/other engines separately.
- [ ] Document full rootfs packaging, signing, migration and rollback; treat trusted release keys as versioned trust roots and do not embed a private signing key or token.
- [ ] Before install: read-only `slotctl status`, free space, kernel/arch, current update source, pairing state, and backupability; require explicit owner approval for **every write**, reboot and pairing change.
- [ ] Hardware verification: voice input -> Japanese transcript -> Japanese Muse answer -> Japanese speaker output; test reconnect/OTA fallback; memory and latency.

## Firmware trust: important
The official base's `internal/update/trust.go` trusts **official** release signatures and `releases_cronos.go` points at official releases. The original Muse fork has **different** trust roots and update feed. A custom v1.6.26-derived image must not be distributed or auto-updated until an independent owner-controlled signing key and release channel are designed, recovery is validated, and the existing user's A/B slot/pairing files have been backed up. Signing and flashing are strictly separate approval gates.

## Resource and safety
Public repository standard GitHub-hosted Linux runners do not consume private-repository included-minute quota, but we still use one manual-gated integration test per milestone. Other commits carry `[skip ci]`. No release builds, no CI artifacts and no remote-server changes until explicitly authorized.

Current Echo firmware (user reported) is `v1.2.0-muse.3`; this repository work **does not change the device**. The persistent Muse pairing state is `layout.StateDir + "/muse.json"`; inspect/backup securely on-device only, never print or commit it.

## Next owner checkpoint
Once compile and code review are complete, present a non-destructive list of hardware inspection commands and explain exactly what the future firmware write will change, what data is preserved, and how to recover. Owner must explicitly approve install/slot switch/reboot. Maintain this document with real test evidence, not predictions.
