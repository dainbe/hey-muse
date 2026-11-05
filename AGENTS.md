# Agent development contract — Japanese voice on Echo Show 5

## Goal
Support natural Japanese spoken input and audible Japanese replies on Echo Show 5 second generation (`cronos`), retaining Hey Muse pairing, wake word and avatar. STT and TTS engines are implementation details, not requirements. `2357gi/hey-muse` is reference-only.

## Source baselines
- Original Hey Muse: `wobsoriano/hey-muse`, `muse` at `1b574cca0a49cb4aa36a69600663434fd2d313fa`.
- TECHO5 stable baseline verified 2026-11-06: `HuskerMinion/techo5` `v1.6.26` (`82722e87fdd20f5e4e92ed6d3eadad55bcd56d6d`). Recheck before actual port and pin the selected SHA.
- Existing Echo runs `v1.2.0-muse.3` per owner. Do not infer that its SSH, pairing backup, rootfs slot, or signing key is available.

## Implementation sequence
1. Read-only source audit: diff official stable TECHO5 vs Muse, including config migration, pairing state, update signing and A/B rollback.
2. Forward-port Muse BLE pairing, wake word, dialog state, display/avatar hooks and setup UI onto an official stable baseline. Keep the original `muse` branch unchanged.
3. Prove an end-to-end Japanese voice turn. Initially use Muse transcription or an external STT and a trusted Japanese TTS endpoint if needed.
4. Add a swappable on-device STT backend; Wyoming is a suitable boundary but ONNX is not mandatory. Evaluate Moonshine Tiny/Base JA and independent alternatives on the actual `cronos`.
5. Evaluate local Japanese TTS and retain an external Japanese TTS option. Built-in SVOX Pico is English-only and is not a Japanese TTS fallback.
6. After voice is stable, separately implement MagicMirror without Home Assistant and optional Muse → Hermes read-only delegation.

## Validation
- Run Go unit tests, ARMv7 cross-build, and static checks in CI; separate cross-compilation from hardware success.
- First acceptance: Japanese user utterance → correct Japanese transcript → Muse response → audible Japanese speech from Echo.
- Report latency (P50/P95, including first audio), model load time, peak memory/RSS, OOM/swap behavior, and fallback on physical hardware.
- A model is not considered supported on the ~1 GB Echo solely because the binary compiles.

## Safety
- Never flash, reboot, change slot, commit/rollback firmware, unpair, reset pairing, change network/firewall, rotate credentials, delete data, or alter the existing Ubuntu LLM server without an explicit user approval after impact and rollback are explained.
- Keep Muse tokens, signing keys, device identity, local secrets, Wi-Fi credentials, private voice recordings and calendar data out of Git, CI, PR comments and logs.
- Use feature branches and draft PRs. Do not merge to `muse`, publish releases or sign installable images without explicit approval.
- CI artifacts should be marked development/test; do not suggest automated installation.
- Maintain `docs/JA_VOICE_STATUS.md` with verified evidence, tests, risks and current blockers.
