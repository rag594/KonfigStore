# KonfigStore - A redis backed multi-tenant config store


A multi tenant configuration management. Configuration are like key/value pairs for a tenant.These configurations are specific to feature/workflow that a tenant uses in your system.

**Note**: Configurations in this case are not application specific configurations.

**Features**:
- Values are of JSON format
- Developers/users can define their custom values as structs.
- Grouping/Categorisation of config to ease the fetch of config of similar category/groups(do we need multiple groups per config?)
- default TTL and custom TTL for each config
- Register/De-Register a configuration using a hook based mechanism. Registration of a configuration will simply register a configuration of specific type. 
- Caching of configuration in redis
- Persistent storage backed by rdbms/nosql databases(mysql, postgres and mongo currently)
- Support of different cache write policies(Why? Because certain configurations are critical and we do not want stale values from cache). Each update to a configuration will o
- Current state of a group/category of configuration in cache/db.
- Lineage of a configuration with time(will include changes in configuration, timestamp, updatedBy etc)
- Web based UI for managing the configurations per tenant(RBAC for listing, editing, viewing configurations per tenant). It is optional.
