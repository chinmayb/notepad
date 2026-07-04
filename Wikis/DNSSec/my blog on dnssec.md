
Previously, we talked about the evolution of DNS.  We went through the how the resolvers get the answer for a domain and cache the responses achieving scale.
If you notice there's a possibility of bad actors come in between and act as authoritative server for a particular domain and inject the fake server IPs into the resolver's cache and resolver wouldn't come to know about it. This is called DNS cache poising [1]. More on that later.

So the important thing to fix is to know for resolver "how do I know if the answer for a particular domain really came from the valid name server that owns the zone?"  This is not like HTTPS, where the data is encypted, The zone is dnssec signed , but its record data is visible through out the wire.



[1] DNS cache poisoning 