# Table: simple_analytics_events

https://docs.simpleanalytics.com/api/export-data-points

The primary key for this table is **_cq_id**.
It supports incremental syncs.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|metadata|JSON|
|added_iso|Timestamp|
|added_unix|Int|
|browser_name|String|
|browser_version|String|
|country_code|String|
|device_type|String|
|document_referrer|String|
|hostname|String|
|hostname_original|String|
|is_robot|Bool|
|lang_language|String|
|lang_region|String|
|os_name|String|
|os_version|String|
|path|String|
|path_and_query|String|
|query|String|
|screen_height|Int|
|screen_width|Int|
|session_id|String|
|utm_campaign|String|
|utm_content|String|
|utm_medium|String|
|utm_source|String|
|utm_term|String|
|uuid|String|
|user_agent|String|
|viewport_height|Int|
|viewport_width|Int|