
Read data from the system:
- Cache aside
- Read through

Write data to the system:
- Write around
- Write back
- Write through


1. Always have a function to reload the -entire- cache.
Remember it's a cache, durability not guaranteed.
2. Have a function to prefill. Then you can start faster.
3. Have a function to validate it. Even better have it asynchronous.
4. Know when a little more memory for the DB will invalidate you needing to use the cache, DBs are amazing


https://www.linkedin.com/feed/update/urn:li:activity:7088906907037360128/



---

The most common cache pattern used in applications is the "Cache-Aside" or "Lazy-Loading" pattern. This pattern involves manually managing the cache, where the application code is responsible for checking the cache before accessing the data source (e.g., a database) and populating the cache if the data is not already present. Here's how the Cache-Aside pattern works:

1. **Cache Check**:
   When the application needs to access data, it first checks the cache to see if the data is already stored there. If the data is found in the cache, it is returned to the application, saving the need to fetch it from the data source.

2. **Cache Miss**:
   If the data is not found in the cache (a cache miss), the application fetches the data from the data source (e.g., a database, an external API).

3. **Cache Population**:
   After fetching the data from the data source, the application populates the cache with the newly fetched data, associating it with a key for future access. The application typically sets an expiration time for the cached data to ensure it doesn't become stale.

4. **Subsequent Accesses**:
   On subsequent accesses to the same data, the application first checks the cache. If the data is still within its expiration period, it is returned from the cache. If it has expired, the application follows the same process of fetching, populating, and updating the cache.

The Cache-Aside pattern is simple and widely used because it provides control over what data gets cached and when it should be updated. However, it also requires careful management to avoid cache staleness and ensure that the cache doesn't consume too much memory.

It's worth noting that there are other cache patterns as well, including:

- **Write-Through**: In this pattern, data is written or updated both in the cache and the data source simultaneously, ensuring that the cache is always up-to-date.

- **Write-Behind (Write-Behind Cache)**: Write-Behind caching involves writing data to the cache immediately and asynchronously updating the data source. This pattern can improve write performance at the expense of potential data source inconsistencies.

- **Read-Through**: In the Read-Through pattern, the cache is responsible for reading data from the data source when a cache miss occurs. The application code doesn't interact with the data source directly.

- **Write-Through with Write-Behind**: This pattern combines elements of Write-Through and Write-Behind caching, providing a balance between write performance and data consistency.

The choice of cache pattern depends on your application's specific requirements, data access patterns, and performance considerations. Cache-Aside is commonly used because it's straightforward and provides a good balance between control and performance. However, for more complex scenarios, other cache patterns may be more suitable.






---


##### Redis ***Streams***

TL;DR persistence, ack, at-least-once, consumer, immutable log

##### Redis ***Pub\/Sub*** is an:

TL;DR no persistence, no ack, real-time, low latency, high throughput, at-most-once

- **at-most-once** messaging system that allows publishers to broadcast messages to one or more channels. More precisely, Redis Pub/Sub is designed for real-time communication between instances where low latency is of the utmost importance, and as such doesn’t feature any form of persistence or acknowledgment. The result is the leanest possible real-time messaging system, perfect for financial and gaming applications, where every millisecond matters.

- is an extremely lightweight messaging protocol designed for broadcasting live notifications within a system. It’s ideal for propagating short-lived messages when low latency and huge throughput are critical.

**Message queue** = mutable state, push, deleted after successfully read/processed
**Event stream** = immutable state, pull, log of events, never deleted (eventually trimmed, cold stored or remove)

src: [Task Queue vs Event Stream](https://redis.com/solutions/use-cases/messaging/)

