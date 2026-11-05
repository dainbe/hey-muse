# Echo Show 5 Gen 2 (cronos): read-only preflight

**Purpose:** collect the minimum device facts needed before planning any update from `v1.2.0-muse.3` to an official TECHO5 v1.6.26-derived Muse build. **No update is authorized.**

## Read-only commands, executed on the Echo only when the owner is present

These commands do not modify firmware, pairing or services:

```sh
uname -a
cat /proc/meminfo | head -n 6
df -h / /data /store 2>/dev/null
command -v slotctl
slotctl status
/usr/local/bin/techo5 --version
# Do not display contents: existence and size only
test -f /data/misc/techo5/muse.json && stat -c 'pairing state bytes: %s' /data/misc/techo5/muse.json
test -d /data/misc/techo5/muse-avatar && du -sh /data/misc/techo5/muse-avatar
```

Some paths or tools may differ by installed firmware. Report missing paths/commands rather than attempting automatic recovery. Do not run `slotctl mkstore`, `install`, `switch`, `commit`, `rollback` or a reboot.

**Do not paste:** `/data/misc/techo5/muse.json`, saved config, Muse token, BLE keys, private key, authorized SSH keys, Wi-Fi credentials, full `env`, token-bearing logs or private recordings. Report only whether the pairing state file exists and its size, without showing its contents, if necessary.

## Hardware gates that require explicit separate owner approval
1. Back up the current slot and persistence; verify the backups and their secure storage.
2. Confirm that the updater's release feed and public trust key match the intended custom release; never reuse the official private signing key.
3. Build a reversible installation plan, including which slot will be written, what happens to existing pairing, and a proven rescue path if the new kernel/rootfs cannot boot.
4. Only then request approval for rootfs install, slot switch, reboot, OTA trust changes, or repair.
5. Post-boot: confirm network/SSH, Muse BLE pairing, wake word, avatar, Japanese voice input/output and rollback ability.

The permanent Muse pairing file in the old fork is `layout.StateDir + "/muse.json"`. Its contents must not be copied to a public repository or included in CI artifacts.

## Hardware performance criteria
Report the on-device STT's model load time, recognition delay for a fixed 5-second Japanese sample, peak RSS and remaining memory. Separately measure first-audio latency and quality of Japanese TTS. Only one neural inference model should normally stay loaded on the ~1 GB RAM device until benchmarks justify more.

No claim of hardware compatibility is established by successful Go tests or ARMv7 cross-compile.
