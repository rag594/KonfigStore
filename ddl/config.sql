CREATE TABLE Config
(
    id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    entityId  BIGINT      NOT NULL,
    configKey VARCHAR(50) NOT NULL,
    value     JSON        NOT NULL DEFAULT (JSON_OBJECT()),
    addedOn   DATETIME             DEFAULT CURRENT_TIMESTAMP,
    updatedOn DATETIME             DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `uc_entityId_configKey` UNIQUE (`entityId`, `configKey`)
);