---
title: "DNSSEC Key and Signing Policy"
source: "https://kb.isc.org/docs/dnssec-key-and-signing-policy"
author:
published:
created: 2026-04-13
description: "Since BIND 9.16 there is a new way to maintain DNSSEC."
tags:
  - "clippings"
---
BIND 9.16 introduced a new method to maintain DNSSEC on your zones. In addition to the `inline-signing` and `auto-dnssec` configuration options, there is now `dnssec-policy` (also see the [configuration reference in the ARM](https://bind9.readthedocs.io/en/stable/reference.html#namedconf-statement-dnssec-policy)).

`dnssec-policy` replaces `auto-dnssec`

`auto-dnssec` has been removed since the 9.19.16 development release in favour of the newer, more flexible `dnssec-policy`.

With `dnssec-policy`, you can specify a **K** ey **A** nd **S** igning **P** olicy (KASP) and group all KASP-related configurations together, making your `named` configuration more intuitive when it comes to DNSSEC and making it easier to enable DNSSEC for your zones. Zones can be secured with DNSSEC as simply as this:

```
zone "example.com" {
    type primary;
    file "db.example.com";
    dnssec-policy default;
};
```

`default` is a built-in policy that introduces a single ECDSAP256SHA key and accompanying signatures for the zone `example.com` in this case. Additionally it uses NSEC (not NSEC3 so far) for denial of existence.

Dynamic Zones or Inline Signing Currently Required

The above is not enough to enable DNSSEC. The zones must either be dynamically-updatable (e.g. be configured for dynamic updates, such as from a Kea server) or have `inline-signing` enabled before the above would work. In BIND 9.18, explicit configuration is required for either. As of BIND 9.20, `inline-signing` defaults to `yes`, and associated application of `dnssec-policy` is sufficient.

As this is a built-in policy you don't have to configure it explicitly. For reference, it look like this:

```
dnssec-policy "default" {
    keys {
        csk key-directory lifetime unlimited algorithm ecdsa256;
    };
    ...
    # Key timings
    dnskey-ttl PT1H;
    publish-safety PT1H;
    retire-safety PT1H;
    purge-keys P90D;
    ...
    # Signature timings
    signatures-refresh P5D;
    signatures-validity P14D;
    signatures-validity-dnskey P14D;
    ...
    # Zone parameters
    max-zone-ttl P1D;
    zone-propagation-delay PT5M;
    parent-ds-ttl P1D;
    parent-propagation-delay PT1H;
};
```

You can specify your own policy if you prefer a separate KSK and ZSK, a different cryptographic algorithm, or NSEC3. You can also tweak the timings, such as for signature validity and TTLs.

Rollover correctness is guaranteed by [DNSSEC records state machines](https://nlnetlabs.nl/downloads/publications/satin2012-Schaeffer.pdf). BIND 9 tries to follow the timings from the policy, but does not apply them if they would result in the zone becoming bogus.

Below we discuss some migration strategies, as well as how to change your DNSSEC policy.

## Turn on DNSSEC

If your zone is currently unsigned, it is now very easy to turn signing on: simply add `dnssec-policy default;` to your zone configuration. Automatic DNSSEC management requires either the zone to be dynamic (allowing dynamic updates), or `inline-signing` to be enabled. Then reconfigure `named`.

It will take some time before it is safe to submit the DS to your parent, to avoid a situation where a resolver fetches the DS but still has a corresponding DNSKEY RRset in the negative cache. When the CDS and CDNSKEY records for the key are published it is safe to submit the DS, see the section on [KSK rollover](https://kb.isc.org/v1/docs/dnssec-key-and-signing-policy#ksk-rollover) below.

If you have many zones and you want to turn on DNSSEC for all of them with the same policy, you can also set the `dnssec-policy` option at `option` or `view` level.

Multiple zones with the same policy

BIND 9 does not share keys between zones. In other words, if you use the same policy for different zones, a separate set of keys is created for each zone.

BIND 9 as a recursive name server

When running BIND 9 as both a recursive and authoritative name server it is strongly recommended to set `dnssec-policy` at the `zone` level. When setting a `dnssec-policy` at `option` or `view` level, this will also apply to the default built-in empty "arpa" zones.

## Migrate to dnssec-policy

If your zone is already signed, hats-off to you! But how do you migrate to the new `dnssec-policy`? There are a couple of items to take into account.

### Wait for rollovers to complete

When you are currently rolling keys, it is not the best time to migrate to `dnssec-policy`. BIND 9 is able to detect when keys are active or inactive, but it is good practice to not mix the two operational procedures. You can check if a rollover is currently in progress using `dig DNSKEY <domain>`. If you see two ZSKs or two KSKs for the same algorithm you are in the process of rolling. ZSK have flags of 256. KSK 257. As shown below where there are two keys of algorithm 257 (the 257 just after each DNSKEY for example.com in the answer section).

```
dig +multi DNSKEY example.com
...
;; ANSWER SECTION:
example.com.    2350 IN    DNSKEY 257 3 13 (
                /ffIEAFFXxGkHmy2a+YnDCBXZbG03Lhao3btWacaPP0h
                5MS/TEEAsBsKvJrR4dhx8RukCEeQcP96ixk3AeI5sw==
                ) ; KSK; alg = ECDSAP256SHA256 ; key id = 2295
example.com.    2350 IN    DNSKEY 257 3 13 (
                9J9CRaIfbbuwdORoV9dJ3Bs0kEDkvrXjwUx23JwUMu/J
                U0lCpb6T/nmbB0YeY4quXGGayS0bJLR34TwPgHZxpQ==
                ) ; KSK; alg = ECDSAP256SHA256 ; key id = 13297
...
```

### Match your existing keys

First, inspect your existing keys. The built-in `default` policy uses a single ECDSAP256SHA key (combined signing key (CSK), performing both the duties of a KSK and a ZSK, also known as Single Type Signing Scheme); if your zone uses another algorithm or has a separate KSK and ZSK, setting the `default` policy to your zone immediately starts a rollover because the existing keys do not match the given policy. This is fine if you want to roll over the keys anyway, but if you want to use the same DNSSEC strategy you should first create a `dnssec-policy` that matches your current situation.

BIND 9.16 with Single Type Signing Scheme zones

For BIND 9.16.20 and earlier, the legacy key metadata has no information about the role of keys. It determines the role from the key flags: 256 means it is a ZSK, and 257 is converted to a KSK. In other words, migrating a CSK won't work. Versions 9.16.21 and higher are able to migrate keys in Single Type Signing Scheme zones (those signed with just one key) to CSK, so if you are using at least 9.16.21, the keys will be migrated properly.

For example, if your zone is now signed with a RSASHA256 2048-bit KSK and a RSASHA256 1024-bit ZSK, you should create the following policy:

```
dnssec-policy "myway" {
    keys {
        ksk lifetime unlimited algorithm rsasha256 2048;
        zsk lifetime P60D algorithm rsasha256 1024;
    };
};
...
zone "example.com" {
    dnssec-policy myway;
};
```

In this example the KSK `lifetime` is set to `unlimited`, so no KSK rollovers are scheduled. The ZSK `lifetime` uses an ISO 8601 format to indicate that the ZSK needs to be replaced every 60 days. You can use different values for the key lifetimes; the timing metadata in the existing key files ("; Inactive:...", "Removed:...") is adjusted accordingly. You do not need to specify the same time period as you were using previously.

`algorithm` can take an algorithm number or mnemonic. Specifying the length is optional, but when migrating to `dnssec-policy` it is better to be explicit, ensuring the policy matches your existing keys.

Note: Private key format must be v1.3

In order to use existing keys for smart signing, the `Private-key-format` must be set to `v1.3`. Keys with older versions do not allow metadata and thus are not compatible with `dnssec-policy`. You can use `dnssec-settime -f` to migrate your existing keys to the newer format. You can check the `Private-key-format` by examining the contents of your zone's private key file, which should contain a line similar to this:  
`Private-key-format: v1.3`  
The default location of the file is `/var/cache/bind`. The filename will contain the name of the zone and end in `.private`. For example: `Kexample.com.+008+12345.private`.

#### Key states

When you reconfigure `named` to introduce a `dnssec-policy` that matches your pre-existing keys, their key files are retained so that they continue to be used to maintain your zone's DNSSEC signatures. In addition, `named` creates the accompanying `dnssec-policy` state files needed for each key (e.g. `Kexample.com.+008+12345.state`).

When your existing active keys have been used for a long time, the key states should have been initialized to `omnipresent` when you restarted `named` to migrate to `dnssec-policy`. That means that they have existed in the zone long enough that all resolvers either have these records currently in cache, or they have to refresh them from the authoritative servers. The significant factor is that no resolver would be expected to have a previous instance of this zone's DNSKEY RRset in cache.

If a key has only recently been introduced, the key state will most likely be set to `rumoured`, meaning that some resolvers may be using it already, whereas others probably haven't seen it yet. When determining key states during a migration, `named` looks at the key timing metadata and takes into account TTLs and propagation delays to determine if the key has existed in the zone for "long enough" to be given state `omnipresent` versus `rumoured`.

Also, inactive keys are migrated; their key states are set to either `unretentive` (a predecessor key that is being rolled) or `hidden` (the key has been removed from the zone long enough that no resolver knows about it).

Let's examine a key state file:

```
; This is the state of key 12345, for example.com.
Algorithm: 8
Length: 2048
Lifetime: 0
KSK: yes
ZSK: no
Generated: 20210818090144 (Wed Aug 18 11:01:44 2021)
Published: 20210818090144 (Wed Aug 18 11:01:44 2021)
Active: 20210818090144 (Wed Aug 18 11:01:44 2021)
PublishCDS: 20210818090144 (Wed Aug 18 11:01:44 2021)
DNSKEYChange: 20210818090149 (Wed Aug 18 11:01:49 2021)
KRRSIGChange: 20210818090149 (Wed Aug 18 11:01:49 2021)
DSChange: 20210818090149 (Wed Aug 18 11:01:49 2021)
DNSKEYState: rumoured
KRRSIGState: rumoured
DSState: rumoured
GoalState: omnipresent
```

This file contains all the information required to determine the next action with respect to key maintenance. The first five lines determine the role of the key: This key is a 2048 bit RSASHA256 KSK for the zone `example.com` that has no lifetime set. Next, there are is a bunch of timing information about when this key is expected to be published, when it should start signing, and when the CDS (and CDNSKEY) should be published. We also keep track of when the state machine for various resource records have changed states, and what state they are in (`hidden`, `rumoured`, `omnipresent`, or `unretentive`).

You can read more about key states [here](https://kb.isc.org/docs/dnssec-key-states).

DSState stuck in rumoured?

If you see the `DSState` stuck in `rumoured` after the migration, you need to run `rndc dnssec -checkds published example.com` to tell BIND that the DS is already published in the parent zone. Be sure and confirm that the DS has actually been published before performing the command (see [KSK rollover](https://kb.isc.org/v1/docs/dnssec-key-and-signing-policy#ksk-rollover) for details about checking the DS state).

#### dnssec-settime

You can modify the key state by hand or with `dnssec-settime`, although this is not recommended unless you know what you are doing. When adding `-s` to `dnssec-settime`, the state file also will be modified. In addition to the key timing metadata, you can update the key goal and the key states:

```
-g state
    This option sets the goal state for this key. 

-d state date/offset
    This option sets the DS state for this key.

-k state date/offset
    This option sets the DNSKEY state for this key.

-r state date/offset
    This option sets the RRSIG (KSK) state for this key.
              
-z state date/offset
    This option sets the RRSIG (ZSK) state for this key.
```

### NSEC3

If you want to use NSEC3, set the `nsec3param` option:

```
dnssec-policy "myway" {
    nsec3param;
};
```

The `nsec3param` option can take sub-options for additional tuning, but current best practice is to use the defaults. Do not use extra iterations, salt, or opt-out unless their implications are fully understood. See [RFC 9276](https://datatracker.ietf.org/doc/html/rfc9276.html).

## Change your policy

If you decide to change your policy - for example, you want to do an algorithm rollover - it is as simple as editing your policy and reconfiguring `named`.

Here is an example where we have changed the `myway` policy to use a single ECDSAP256SHA256 key instead of the RSASHA256 KSK/ZSK split approach. After reconfiguring `named`, the corresponding zones will introduce a new key and phase out the old keys:

```
dnssec-policy "myway" {
    keys {
        # ksk lifetime unlimited algorithm rsasha256 2048;
        # zsk lifetime P60D algorithm rsasha256 1024;
        csk lifetime unlimited algorithm ecdsa256;
    };
};

zone "example.com" {
    ...
    dnssec-policy myway;
};
```

This policy now is the same as the *default* policy, so we could have also changed the configuration like this:

```
zone "example.com" {
    ...
    # dnssec-policy myway;
    dnssec-policy default;
};
```

Rather than changing the global `myway` policy, we are attaching a different policy to a single zone. The result is the same for this zone, but allows you to keep some zones using the global (or per-view) `myway` policy while at the same time migrating one or more other zones to the `default`.

## Key rollover

Key rollovers are done automatically and on a schedule. How often rollovers occur depends on the `lifetime` of the key. If set to `unlimited`, no key rollovers will be scheduled; you would need to use `rndc dnssec -rollover` to start or schedule a manual key rollover.

A note on signing strategies

ZSK rollovers use a Pre-Publish signing strategy; that is, the DNSKEY is prepublished alongside the existing DNSKEY and once propagated, the RRSIGs in the zone are gradually replaced. This happens fully automatically; no additional action from the operator is required.

KSK rollovers use a Double-KSK signing strategy. Just as with ZSK rollovers, the DNSKEY is prepublished next to the existing DNSKEY and after enough time has passed, the DS in the parent may be submitted.

CSK rollovers use a combination of the two signing strategies. The DNSKEY is prepublished and then the zone RRSIGs can be gradually replaced, while at the same time the DS can be submitted.

The ZSK Pre-Publish rollover and KSK Double-KSK are described in more detail in [RFC 7583](https://tools.ietf.org/html/rfc7583).

### Manual key rollover

To perform a rollover independently of any schedule, you can do that using `rndc`:

```
rndc dnssec -rollover -key 12345 example.com
```

This will start a key rollover for key with ID 12345 in zone *example.com*. You can also schedule a rollover at a specific time by adding `-when <datetime>` where the `datetime` is in `YYYYMMDDHHMMSS` format (UTC). The key lifetime is reduced, and when it is time to introduce the new successor key, BIND automatically does so.

You can check the progress of the key rollover as shown below:

```
$ rndc dnssec -status example.com
dnssec-policy: default
current time:  Thu Feb  2 12:00:27 2023

key: 2295 (ECDSAP256SHA256), CSK
  published:      yes - since Sat Dec  3 18:45:40 2022
  key signing:    yes - since Sat Dec  3 18:45:40 2022
  zone signing:   yes - since Sat Dec  3 18:45:40 2022

  Key will retire on Fri Feb  3 09:05:00 2023
  - goal:           hidden
  - dnskey:         omnipresent
  - ds:             omnipresent
  - zone rrsig:     omnipresent
  - key rrsig:      omnipresent

key: 13297 (ECDSAP256SHA256), CSK
  published:      yes - since Thu Feb  2 11:58:46 2023
  key signing:    yes - since Thu Feb  2 11:58:46 2023
  zone signing:   yes - since Thu Feb  2 14:03:46 2023

  No rollover scheduled
  - goal:           omnipresent
  - dnskey:         rumoured
  - ds:             hidden
  - zone rrsig:     rumoured
  - key rrsig:      rumoured
```

Note that two keys are shown. One specifies that it will soon be retired and the second specifies how long it has existed.

#### Importing DNSSEC keys

A DNSSEC Practice Statement may require that the key generation happens on a different system than the signing system. For example, the key may be sanity-checked before it is installed.

If you don't want BIND 9 to create keys, you can import an external key into the `key-directory` before you start a manual rollover. If it matches the `dnssec-policy`, BIND 9 will select the key, rather than generate a new one. A key will only be selected if it has not been used before, i.e. the key has no timing metadata set. To pregenerate a key that has no timing metadata you can use the `dnssec-keygen -G` option. For example:

```
dnssec-keygen -G -a ECDSAP256SHA256 example.com
```

#### Scheduling a rollover

If you want to schedule a key rollover, rather than start one immediately, you can use the `-when` parameter when running the `rndc dnssec -rollover` command. BIND 9 will create a key when it is required to pre-publish the new key.

Again, if you don't want BIND 9 to create the key for you, you can import the key into the `key-directory`.

For example, if you want to replace the key with id 12345 for the zone `example.net` in 60 days time, and you have a key created with id 54321. Then you would first put the pregenerated key into the key directory and then start a rollover at the desired date and time. This may look like:

```
cp Kexample.com.+013+54321.* keys/
when=$(date -d "now + 60 days" +"%Y%m%d%H%M%S")
rndc dnssec -rollover -key 12345 -when $when example.com
```

Reminder: times are assumed to be in UTC

As mentioned above, times default to UTC and date defaults to the system timezone setting. If you wish to be more explicit, you can use the TZ variable; for example:

```
when=$(TZ=UTC date -d "now + 60 days" +"%Y%m%d%H%M%S")
```

### KSK rollover

ZSK rollovers can be fully automated. However, this is not possible for KSK or CSK, because you still need to submit the DS to the parent zone manually before it can be used by resolvers that use DNSSEC validation.

Therefore, it is less usual to give a time-limited key lifetime to a KSK or CSK and more common to instigate a KSK or CSK rollover manually.

The DS for a given key may be submitted if the corresponding CDS and CDNSKEY records have been published.

Lets take the `example.com` where we replace a key with ID 54321 with its successor - the key with ID 12345. The DS record to submit can be retrieved with `dnssec-dsfromkey` utility.

```
dnssec-dsfromkey -2 Kexample.com+008+12345
```

`named` started the key rollover for you, but does not continue through some steps until you signal to BIND that certain important actions have taken place in the parent zone:

- When you have seen that the new DS RR for your new KSK or CSK has been added to the parent zone, you must signal this fact to BIND in order to continue the rollover. You can check the current DS RR by querying a parent zone server `dig +dnssec +norecurse +multi @172.0.2.1 example.com` and inspect the returned DS RRSIG.

```
example.com.        3600 IN    RRSIG DS 8 2 3600 (
                20230222152214 20230201142214 54321 com.
                Aj5nORNdG03mZPo22XMcWvbSLN75EBGDYbJPis5m/6Zt
                zkisSFZb1h9fPaETGdzfu/tmZ1zrGU5QzFarP5GdFh9O
                iZlqJgPWLgEosQIC30BXPfY3uOtqpxhjOz5SimSTSILv
                Jy00fWeJw9dWmeJ6R4lXqGTu0/YOCuqQ2in8UlQ= )
```

Note the 54321 in the above output. Once 54321 is replaced with the new ID of 12345 (and the new key output following), then you know the parent zone has published the new DS RR.

- Similarly, you also need to signal to BIND when a KSK or CSK that is being replaced has had its DS removed from the parent zone.

Once the two keys have had their DS records swapped in the parent zone, you can run the following commands:

```
rndc dnssec -checkds -key 12345 published example.com
rndc dnssec -checkds -key 54321 withdrawn example.com
```

Automation through Parental Agents

The above can be automated by adding `parental agents` into the zone configuration, e.g. `parental-agents { 192.0.2.1; };` where `192.0.2.1` is a trusted IP address of an upstream DNS server that will be able to answer the DS keys that will be published and withdrawn. BIND will notice that the new key is available and will perform the above steps automatically. See the section on [Key Rollover](https://bind9.readthedocs.io/en/latest/chapter5.html#key-rollover) in the ARM for futher details.

## Dynamic zones, inline-signing and files.

This section is a summary of the effect of having either or both an `update-policy` or `inline-signing` configured.

As has been said already, to have `named` sign and maintain your zones automatically using `dnssec-policy`, at least one of the following conditions must be true:

- Zones must be configured to allow dynamic updates, using the `update-policy` statement. Or...
- `inline-signing` must be enabled. This is the default as of BIND 9.20. Older versions will require explicit configuration. This can be done either within the zone definition block, or globally, or a combination of both (if, for instance, you have a global policy that applies to all zones except ones that have a different, local policy).

### What happens in the file system when either of these is enabled?

If neither of these is enabled and a zone is not signed, there will be one zone file containing an SOA, one or more NS and any other records you have provisioned. This is the static data you have in a zone, with no additional data such as DNSSEC signature records (RRSIG), NSEC or NSEC3 records.

Both dynamic update and `dnssec-policy` need a mechanism to be able to store changes to the zone temporarily as `named` is working through processiing the changes. This is done in one or two ways, through the use of journal files.

#### What is a journal file?

Journal files have the extension `.jnl` (journal) or `.jbk` (journal backup) and are created by `named` when it needs to store information temporarily during processing of some task.

Depending on the task being performed and on its configuration, `named` will create these files with names that are relevant to the task. For example, a journal file created because of dynamic updates to a zone will have `.jnl` appended to the name of the zone file.

Journal files are binary

Since journal files are a scratchpad, used by `named` to keep track of data that will eventually be synchronised, the information contained in them is in a compact binary format, specific to BIND. It is not intended to be human-readable directly. However, the utility `named-journalprint` is available to read them and to print the contents in a human-readable format.

#### Dynamic updates

When a zone is enabled for dynamic updates an additional file is created by `named`; the *journal* file. Its name is the same as the zone file, plus the extension `.jnl`. The journal file acts as a temporary cache; a place for `named` to store dynamic updates to the zone as they happen. These updates will be merged into the zone file in the background.

The same is true of DNSSEC signatures and other records when **only** dynamic updates are enabled for a zone. Additional DNSSEC records are added to the original zone file, which now contains regular (e.g. **A**, **AAAA**, **SRV** etc.) RRSETs **and** their corresponding signatures, in the form of RRSIG records.

#### inline-signing

If a zone is not enabled for dynamic updates, there needs to be another mechanism for `named` to create and store signatures as it works its way through signing the data in the zone and `inline-signing` achieves that and more. It works in a similar way to dynamic updates, in that it also creates a journal file in which to store its work temporarily. Once that is finished, changes are merged into **a version of** the zone file in the background. However, unlike dynamic updates, where all work and changes are made in the existing zone file and its corresponding journal file, `inline-signing` creates new files for the DNSSEC material and for some other purposes. These files are:

| filename | purpose |
| --- | --- |
| <raw\_zone\_filename>.signed | Contains all records from the raw zone file plus their signatures and other DNSSEC material |
| <raw\_zone\_filename>.signed.jnl | Contains working data for the signing process |
| <key\_filename>.state | A human-readable text file containing live information about the current state of a key or keys for a given zone |

With `inline-signing`, its working process and the results are kept separate from the original data on which they are based, so there is no overlap with the files used by `update-policy`.

To summarise:

- If only dynamic update is enabled, `dnssec-policy` will work in its journal and store signatures alongside the records that have been signed in the original zone file.
- With `inline-signing` enabled, work is done in dedicated journals and signatures are stored alongside the records that have been signed in **new** `...signed` files.

It may be important to be aware of this difference in behaviour if `inline-signing` is enabled in a configuration where zones (with dynamic updates enabled) were already being signed by `dnssec-policy`.  
Before the change, the original zone files will contain regular data (managed dynamically through `nsupdate` **plus** DNSSEC data.  
If `inline-signing` is now enabled, which is the default for `dnssec-policy` as of BIND version 9.20, `named` will now try to sign everything in what it considers to be the raw zone file, which already contains signatures. The result is a `...signed` zone file that contains new signatures for regular records plus new signatures of the previous signatures.

If this would be the case, there are three ways to avoid it happening:

1. Convert zones from secure to insecure, thus removing signatures and reverting to a pre-DNSSEC state. This can be done using the built-in `dnssec-policy "insecure";`. Following this, an upgrade can proceed, using the defaults for `inline-signing` and new `...signed` zone files will be produced, as described above.
2. Before upgrading to 9.20 from a version in which `inline-signing` was not an option within the `dnssec-policy` block (therefore the default for `inline-signing` per zone was "no"), add `inline-signing no;` explicitly to every zone that would be affected. This more specific configuration overrides the default of "yes" that applies globally and means that zones will carry on working as they did before.
3. In the `dnssec-policy` configuration block, add `inline-signing no;` before upgrading, thus overriding the global default and keeping zones as they are

---

We have discussed a couple of `dnssec-policy` migration and update strategies and potential pitfalls to be aware of when using it. More information can be found in our [DNSSEC Guide](https://bind9.readthedocs.io/en/latest/dnssec-guide.html).