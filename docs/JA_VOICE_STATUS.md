# Japanese voice forward-port — engineering status

Updated: 2026-11-06 (JST)

## Current state
- GitHub fork verified: `dainbe/hey-muse`, original `muse` untouched.
- Work branch: `work/ja-voice-foundation`, created successfully after GitHub integration access was added.
- Source audit and integration bootstrap in progress. **No forward-port code, model benchmark, real-device operation or firmware update performed yet.**
- Repository settings: Public allows clone/read, while connected GitHub App authorization is necessary for writes. Write access was successfully verified by creating the work branch.

## Pinned reference points
| Role | Ref | Status |
| --- | --- | --- |
| Existing Hey Muse | `wobsoriano/hey-muse` `1b574cca0a49cb4aa36a69600663434fd2d313fa` | Confirmed |
| Official stable TECHO5 | `HuskerMinion/techo5` `v1.6.26` `82722e87fdd20f5e4e92ed6d3eadad55bcd56d6d` | Confirmed as of 2026-11-06 |
| Other reference only | `2357gi/hey-muse` | Code reviewed; not a mandatory base |
| Owner device | Echo Show 5 gen2 (`cronos`), `v1.2.0-muse.3` | User-reported; read-only recheck before writes |

## Technical assessment
- Official stable TECHO5 currently requires `go 1.27.0` in `echod/go.mod`, old Muse fork `go 1.26.0`. Match build toolchain to the checked-out module.
- Echo `cronos`: MediaTek MT8163, Cortex-A53 x4, about 996 MB RAM; arm64 Linux kernel with 32-bit userspace. On-device Japanese STT/TTS must be performance-tested.
- Muse supports spoken replies via an OpenAI-style `/audio/speech` endpoint, or built-in SVOX Pico; Pico is English-only and not valid Japanese TTS.
- STT backends and TTS backends should be selected independently. No ONNX-only constraint.
- Existing device pairing, wake-word and avatar must be preserved through forward-port. Trust roots for OTA signing differ across forks.

## Milestones
- [ ] **M0 Source audit:** Map upstream/Muse deltas, stable tag, pairing-state compatibility, release-signing and recovery.
- [ ] **M1 Forward port:** Build latest stable TECHO5 with Muse pairing, wake word and avatar; Unit Tests + ARMv7 build.
- [ ] **M2 Japanese voice:** Japanese transcription → Muse text reply → Japanese speech audio, initially with optional external TTS.
- [ ] **M3 Replaceable local STT:** Wyoming-compatible selector, Moonshine Tiny/Base and additional Japanese STT candidates; microbenchmarks and safe fallback.
- [ ] **M4 Local Japanese TTS:** On-device lightweight alternatives vs external Japanese TTS; report voice quality and performance.
- [ ] **M5 Device trial:** Owner-approved slot install, rollback, on-device memory/latency and regression tests.
- [ ] **M6 Later:** MagicMirror no-Home-Assistant web dashboard, Muse → Hermes read-only delegation.

## Acceptance gates
| Gate | Evidence | State |
| --- | --- | --- |
| Source pin | Official tag + Muse SHA recorded | Partial |
| Cross-build and Go tests | GitHub Actions logs for integration branch | Pending |
| Muse regression | Pairing, wake, response, avatar, rejoin | Pending |
| Japanese E2E | Actual Japanese speech in and out | Pending |
| Resource budget | P50/P95, peak RSS and OOM on `cronos` | Pending |
| Safe upgrade | A/B slot, backup, signing key verification | Pending |

## Non-destructive policy
Do not perform `slotctl install/commit/switch/mkstore`, `fastboot flash`, restart or modify running device, reset pairing, or change existing Ubuntu LLM infrastructure without explicit separate approval. Do not store tokens or recordings in this repository. Do not merge to `muse` automatically.
