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

The integration branch was created *directly from the official stable commit*, not from the Muse fork. `baseline/techo5-v1.6.26` pins the same official commit, and Draft PR #2 targets that baseline for a reviewable official-to-Muse diff. The original `muse` branch is untouched and protected. Initial diff inspection found 174 Muse-changed files, of which 34 were also changed upstream.

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
- Unsigned preview run `42079486348` failed on a Muse pairing test race, fixed by waiting for fixture `Run` startup. Next `42079935858` passed Go tests but Ubuntu 26.04 runner no longer offers a concrete `qemu-user-static` package; changed to `qemu-user-binfmt`, confirmed by official Ubuntu package index.
- Rootfs preview `42080449780`: verified Go tests, ARM cross-compile, QEMU package setup and public inputs, but failed because mandatory AirPlay 2 helper binaries were not yet staged.
- Rootfs preview `42081277834`: verified official `rootfs-v1.6.26.tar.gz` with its published SHA256SUMS; extracted official AirPlay 2, FFmpeg, Tailscale helpers (confirmed ARMv7/arm64). Hit GitHub runner UID/GID ownership restrictions during Alpine rootfs extraction.
- Rootfs preview `42081730021`: failed a CI-only namespace `chown 0:42` preflight. Left normal WSL rootless mode unchanged; explicit temporary GitHub Runner sudo build mode now available.
- Rootfs preview `42081277834`: **failed** Alpine service-user ownership in a rootless GitHub Runner namespace. Separate developer WSL default remains rootless.
- Rootfs preview `42081730021`: CI-only `chown 0:42` preflight confirmed hosted-runner UID/GID mapping limitation, no device touched.
- **Rootfs preview `42082277352`: success** — all Go tests, QEMU setup, public input fetch, pinned official v1.6.26 SHA256-verified helpers, complete custom unsigned rootfs tar, SHA256SUMS/CONTENTS and 1-day-retention artifact upload. Artifact ID `13441421638`, 49,303,810 bytes; expires at `2026-11-06T18:16:23Z`. [Run](https://github.com/dainbe/hey-muse/actions/runs/42082277352).
- Rootfs generated in a **temporary GitHub-hosted runner only**; isolated opt-in `TECHO5_CI_ROOTFS_SUDO=1` used for packages that require file ownership IDs. Default developer WSL path stays rootless. Heavy workflow restored to manual-only.
- **Never install** unsigned preview output before hardware safety gate and explicit approval; no physical Echo STT/TTS operation is verified.
- Workflow `.github/workflows/port-ci.yml` performs Go tests and ARMv7 compile only; no signing, auto-installation or artifact upload. Trigger manually or by editing `.ci/run-integration`.

## Latest CI follow-up (2026-11-06)
- Upstream PR-check run `42082288309` failed on panel-specific `go vet`: `internal/feature/muse/gadget.go:33` declared `const StatePath = layout.StateDir + "/muse.json"`, but `layout.StateDir` is a variable on `panel` builds.
- Corrected `StatePath` to `var`, matching the upstream `layout` path definitions. This preserves the existing pairing filename and directory; it does **not** change the data or device state.
- The follow-up PR check `42125847545` confirmed the pairing `StatePath` error was resolved, then found the same `const`/runtime-directory mismatch in `feature/display/render_avatar.go`. That `avatarDir` has also been changed to `var` without changing its value or persistent data. The next full PR check will validate all Build Tags.

## Formatting and device-tag CI follow-up
- Full PR check `42126780427` completed all `go build` and `go vet` variants (show, dot, dot+donut, spot, panel), including the two corrected variable paths. Its only failure was `gofmt -l` reporting nine integration-edited Go files.
- A one-time, exact-path formatter workflow `42127654760` applied `gofmt` to only those nine files and committed `523c2d79483fcc2241d89fa211000e1b318f968a` to the integration branch. It did not rebuild/flash firmware or touch baseline/Muse branches.
- **Final PR check `42127736906`: SUCCESS**. All Go Build/Vet tag variants (show, dot, dot+donut, spot, panel), gofmt, translation status, installer tests, and Docker-based boot layer tests passed. No firmware install, signed release, or physical-device tests were run.

## Outstanding integration
- [x] All `go test ./...` and ARMv7 daemon compile green on pinned v1.6.26 (run `42079152900`, Go 1.27 module).
- [ ] Preserve upstream display features and test Muse avatar during a turn.
- [ ] Verify local Wake word paths and Muse pairing client on actual `cronos`.
- [x] Optional Wyoming STT text path implemented and tested offline (including no audio upload after local STT failure).
- [ ] Japanese STT model deployment and real-device benchmark (Moonshine Tiny/Base, independent alternatives) pending.
- [ ] Japanese TTS: set a trusted Japanese-capable speech endpoint first. Built-in Pico is **English-only**, cannot serve as Japanese fallback. Investigate local Open JTalk/other engines separately.
- [x] Compile and stage a reproducible **unsigned preview** rootfs using official signed stable helpers and public inputs; [run #42082277352](https://github.com/dainbe/hey-muse/actions/runs/42082277352).
- [ ] Owner-controlled signing, pairing migration and A/B rollback procedure; no private signing key or token may be committed.
- [ ] Before install: read-only `slotctl status`, free space, kernel/arch, current update source, pairing state, and backupability; require explicit owner approval for **every write**, reboot and pairing change.
- [ ] Hardware verification: voice input -> Japanese transcript -> Japanese Muse answer -> Japanese speaker output; test reconnect/OTA fallback; memory and latency.

## Firmware trust: important
The official base's `internal/update/trust.go` still trusts **official** release signatures, but the integration branch changes `releases_cronos.go` to this fork's currently empty release feed so that it cannot automatically replace Muse with the official non-Muse firmware. The original Muse fork has **different** trust roots and update feed. A custom v1.6.26-derived image must not be distributed or auto-updated until an independent owner-controlled signing key and release channel are designed, recovery is validated, and the existing user's A/B slot/pairing files have been backed up. Signing and flashing are strictly separate approval gates.

## Resource and safety
Public repository standard GitHub-hosted Linux runners do not consume private-repository included-minute quota, but we minimize runs. `port-ci.yml` is explicitly sentinel-gated, and `rootfs-preview.yml` has been restored to manual-only after bounded runs, with a 1-day artifact retention and unsigned DO-NOT-INSTALL label. There are no public releases, signed firmware, or remote-server changes.

Current Echo firmware (user reported) is `v1.2.0-muse.3`; this repository work **does not change the device**. The persistent Muse pairing state is `layout.StateDir + "/muse.json"`; inspect/backup securely on-device only, never print or commit it.

## Next owner checkpoint
Offline Build + CI compilation and unsigned rootfs preview are complete. The next gate is **owner-assisted read-only device preflight, followed by an explicit approval decision**. Once compile and code review are complete, present a non-destructive list of hardware inspection commands and explain exactly what the future firmware write will change, what data is preserved, and how to recover. Owner must explicitly approve install/slot switch/reboot. Maintain this document with real test evidence, not predictions.

## OTA isolation
This branch points `releases_cronos.go` at this fork's empty Release feed; owner-controlled signing and migration from the old Muse trust key remain pending.
