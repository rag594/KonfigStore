# KonfigStore - A cache backed multi-tenant config store


A multi tenant configuration management. Configuration are like key/value pairs for a tenant.These configurations are specific to feature/workflow that a tenant uses in your system.

**Note**: Configurations in this case are not application specific configurations, these are configurations specific to a tenant(tenant can be any entity in the ecosystem)

**Features**:
- [x] Values are of JSON format
- [x] Developers/users can define their custom values as structs.
- [x] default TTL and custom TTL for each config
- [x] Register/De-Register a configuration using a hook based mechanism. Registration of a configuration will simply register a configuration of specific type. 
- [x] Caching of configuration in redis
- [x] Persistent storage backed by rdbms/nosql databases(mysql)
- [x] Support of different cache write policies(write-through, write-around, write-back)
- [x] Cache stampede protection
- [ ] Settings/Options at each configuration
  - [x] ttl of the config in cache
  - [ ] distributed cache mode
  - [ ] db persistence mode
  - [ ] db timeout
  - [ ] cache timeout
  - [ ] eager refresh
  - [x] write policy
  - [x] custom configKey
  - [ ] custom cacheKey
- [ ] Monitoring
- [ ] Logging
- [ ] Current state of a group/category of configuration in cache/db.
- [ ] Grouping/Categorisation of config to ease the fetch of config of similar category/groups(do we need multiple groups per config?)
- [ ] Lineage of a configuration with time(will include changes in configuration, timestamp, updatedBy etc)
- [ ] Web based UI for managing the configurations per tenant(RBAC for listing, editing, viewing configurations per tenant). It is optional.
