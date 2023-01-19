# Table: simple_analytics_data_points

https://docs.simpleanalytics.com/api/export-data-points

The composite primary key for this table is (**hostname**, **uuid**).
It supports incremental syncs.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|metadata|JSON|
|added_unix|Int|
|added_iso|Timestamp|
|hostname (PK)|String|
|hostname_original|String|
|path|String|
|query|String|
|is_unique|Bool|
|is_robot|Bool|
|document_referrer|String|
|utm_source|String|
|utm_medium|String|
|utm_campaign|String|
|utm_content|String|
|utm_term|String|
|scrolled_percentage|Float|
|duration_seconds|Float|
|viewport_width|Int|
|viewport_height|Int|
|screen_width|Int|
|screen_height|Int|
|user_agent|String|
|device_type|String|
|country_code|String|
|browser_name|String|
|browser_version|String|
|os_name|String|
|os_version|String|
|lang_region|String|
|lang_language|String|
|uuid (PK)|String|