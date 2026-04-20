---
title: "dig and delv"
source: "https://kb.isc.org/docs/aa-01152"
author:
published:
created: 2026-04-13
description: "BIND 9.10 introduced a new debugging tool as a successor to dig. So, of course, we had to name it delv. It works like dig but understands DNSSEC better."
tags:
  - "clippings"
---
### Delv

Surely you remember the old nursery rhyme that begins with "One, two, buckle my shoe." You might have forgotten that its sixth verse is "Eleven, twelve, dig and delve." How many nursery rhymes have their own [Wikipedia page](http://en.wikipedia.org/wiki/One,_Two,_Buckle_My_Shoe)?

BIND 9.10 contains a new debugging tool that is a successor to `dig`. So, of course, we had to name it `delv`. It works very much like `dig`, but it understands DNSSEC better.

`delv` checks the DNSSEC validation chain using the same code that is used by the BIND 9 DNS server itself. Compared to `dig +sigchase`, `delv` is much closer to what really happens inside a DNS server.

Like all BIND tools, `delv` is fully documented in Appendix B of the BIND ARM. In general, you should use `delv` the same way you have been using `dig`. But `delv` is not an enhanced version of `dig`; it is an entirely new program. So here are a few tips.

#### delv FQDN

If you just say `delv name`, and name is in a signed zone, `delv` will report "fully validated" and give you the RRSIG as part of the answer.

#### delv +multi

DNSSEC-related records are often very long. To make its output more readable, `delv` has a `+multi` option that formats large records into multiline reports that are readable in a standard 80-column text window. Men and Mice has provided this example of the use of `+multi`:

```
% delv dnsworkshop.org soa +multi
; fully validated
dnsworkshop.org.        3600 IN SOA ns1.myinfrastructure.org. hostmaster.strotmann.de. (
               86         ; serial
               86400      ; refresh (1 day)
               7200       ; retry (2 hours)
               3542400    ; expire (5 weeks 6 days)
               3600       ; minimum (1 hour)
               )
dnsworkshop.org.        3600 IN RRSIG SOA 8 2 3600 (
               20140321030247 20140219020247 63654 dnsworkshop.org.
               O8mmiuNdXIWG6huaLiQrvKabDY3qivQ3R5qRUZ1IG3wp
               bd0UBnvpazpG01ntk8uZ7wEStScmiY7oYtvRGIHG37mG
               8GFI60CUx3pdXJIpmodfoUBk8cfGsJXFQODIZCTUQiyk
               Pv9I6+wjyseDJJTYlrsBCvAEabPExFKZc7v+L+k= )
```

Like `dig`, `delv` is fully compatible with IPv6:

```
% delv dnsworkshop.org AAAA +multi
; fully validated
dnsworkshop.org.        7200 IN AAAA 2001:470:1f08:f1d::2
dnsworkshop.org.        7200 IN RRSIG AAAA 8 2 7200 (
               20140321025727 20140219020247 63654 dnsworkshop.org.
               gqkc1Xq/UveKrhcXpqOwDsN5HFSqMsPkxXOyCqu9bMyx
               dtnkh0J0Iqukv+uHL/dDQLnPcxjdFqs3N5Jf3BFHdgkG
               tf0UPhNKsuhlsRdo2H5O+TqmLvA1zCsYhH/72vVvxslR
               MiiuZ1ILGpLA2EOyiZu70/ZIU3Ypc3nb8+ydgx4= )
```

#### delv +multi +rtrace

The command `delv +rtrace` prints the extra DNS lookups that `delv` needs to make while validating the reply to a query. This example is from Men and Mice:

```
% delv dnsworkshop.org mx +multi +rtrace
;; fetch: dnsworkshop.org/MX
;; fetch: dnsworkshop.org/DNSKEY
;; fetch: dnsworkshop.org/DS
;; fetch: org/DNSKEY
;; fetch: org/DS
;; fetch: ./DNSKEY
; fully validated
dnsworkshop.org.        3600 IN MX 100 mail.strotmann.de.
dnsworkshop.org.        3600 IN RRSIG MX 8 2 3600 (
                                20140308193355 20140206183355 63654 dnsworkshop.org.
                                hCOcPJrDCXpcVS82FgGEdUhaUmW3XkxXEuEa4AFvzkzi
                                mDcokYNjrW/Hay4NclSWV0jrBwrXABXik5dh7w7KsPkD
                                WKhw/qVvkuiFCm+T5lb9OVkGQAuPhBOplbVgdbZce9L7
                                N2IVTQTLMECKfzCTfKeOtwupJAMPXCt/Xskd5o4= )\`\`
```

#### delv +multi +vtrace

The `+vtrace` option shows the entire DNSSEC chain of validation. This example again from Men and Mice:

```
shell> delv _443._tcp.dnsworkshop.org TLSA  +multi +vtrace
;; fetch: _443._tcp.dnsworkshop.org/TLSA
;; validating _443._tcp.dnsworkshop.org/TLSA: starting
;; validating _443._tcp.dnsworkshop.org/TLSA: attempting positive response validation
;; fetch: dnsworkshop.org/DNSKEY
;; validating dnsworkshop.org/DNSKEY: starting
;; validating dnsworkshop.org/DNSKEY: attempting positive response validation
;; fetch: dnsworkshop.org/DS
;; validating dnsworkshop.org/DS: starting
;; validating dnsworkshop.org/DS: attempting positive response validation
;; fetch: org/DNSKEY
;; validating org/DNSKEY: starting
;; validating org/DNSKEY: attempting positive response validation
;; fetch: org/DS
;; validating org/DS: starting
;; validating org/DS: attempting positive response validation
;; fetch: ./DNSKEY
;; validating ./DNSKEY: starting
;; validating ./DNSKEY: attempting positive response validation
;; validating ./DNSKEY: verify rdataset (keyid=19036): success
;; validating ./DNSKEY: signed by trusted key; marking as secure
;; validating org/DS: in fetch_callback_validator
;; validating org/DS: keyset with trust secure
;; validating org/DS: resuming validate
;; validating org/DS: verify rdataset (keyid=33655): success
;; validating org/DS: marking as secure, noqname proof not needed
;; validating org/DNSKEY: in dsfetched
;; validating org/DNSKEY: dsset with trust secure
;; validating org/DNSKEY: verify rdataset (keyid=21366): success
;; validating org/DNSKEY: marking as secure (DS)
;; validating dnsworkshop.org/DS: in fetch_callback_validator
;; validating dnsworkshop.org/DS: keyset with trust secure
;; validating dnsworkshop.org/DS: resuming validate
;; validating dnsworkshop.org/DS: verify rdataset (keyid=24209): success
;; validating dnsworkshop.org/DS: marking as secure, noqname proof not needed
;; validating dnsworkshop.org/DNSKEY: in dsfetched
;; validating dnsworkshop.org/DNSKEY: dsset with trust secure
;; validating dnsworkshop.org/DNSKEY: verify rdataset (keyid=2611): success
;; validating dnsworkshop.org/DNSKEY: marking as secure (DS)
;; validating _443._tcp.dnsworkshop.org/TLSA: in fetch_callback_validator
;; validating _443._tcp.dnsworkshop.org/TLSA: keyset with trust secure
;; validating _443._tcp.dnsworkshop.org/TLSA: resuming validate
;; validating _443._tcp.dnsworkshop.org/TLSA: verify rdataset (keyid=63654): success
;; validating _443._tcp.dnsworkshop.org/TLSA: marking as secure, noqname proof not needed
; fully validated
_443._tcp.dnsworkshop.org. 3544 IN TLSA 3 0 1 (
                                3E5E70BBA957CA0DAFCB799F15F6236133C0F6C73FA7
                                3762BFFBCA4AF92389CA )
_443._tcp.dnsworkshop.org. 3544 IN RRSIG TLSA 8 4 3600 (
                                20140309145739 20140207135739 63654 dnsworkshop.org.
                                JYkLiFqvrjqiIlm/bA4CaffJ3Iikos31bfEVb2njjIR+
                                /7dudq9pAj898OVZrtqIjmfD7knyCT2nt6Gp/yFYif4k
                                Tt7W2XMhnWecwRnFexhVYp1zg2dkZSw4XcBRMz/F2NkM
                                0xziG9dNFg/6AAs/0ehMurLvRj1ula/UIO/wU5w= )\`\`
```