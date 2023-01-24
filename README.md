# CloudQuery Simple Analytics Source Plugin

[![test](https://github.com/cloudquery/cq-source-simple-analytics/actions/workflows/test.yaml/badge.svg)](https://github.com/cloudquery/cq-source-simple-analytics/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudquery/cq-source-simple-analytics)](https://goreportcard.com/report/github.com/cloudquery/cq-source-simple-analytics)

A Simple Analytics source plugin for CloudQuery that loads data from Simple Analytics to any database, data warehouse or data lake supported by [CloudQuery](https://www.cloudquery.io/), such as PostgreSQL, BigQuery, Athena, and many more.

## Links

 - [CloudQuery Quickstart Guide](https://www.cloudquery.io/docs/quickstart)
 - [Supported Tables](docs/tables/README.md)

## Configuration

The following source configuration file will sync all data points for `mywebsite.com` to a PostgreSQL database. See [the CloudQuery Quickstart](https://www.cloudquery.io/docs/quickstart) for more information on how to configure the source and destination.

```yaml
kind: source
spec:
  name: "simple-analytics"
  path: "simpleanalytics/simple-analytics"
  version: "${VERSION}"
  backend: "local" # remove this to always sync all data
  tables: 
    ["*"]
  destinations: 
    - "postgresql"
  spec:
    # plugin spec section
    user_id: "${SA_USER_ID}"
    api_key: "${SA_API_KEY}"
    websites:
      - hostname: mywebsite.com
        metadata_fields: 
          - fieldname_text
          - fieldname_int
          # - ... 
```

### Plugin Spec

- `user_id` (string, required):

  A user ID from Simple Analytics, obtained from the [account settings](https://simpleanalytics.com/account) page. It should start with `sa_user_id...`

- `api_key` (string, required):

  An API Key from Simple Analytics, obtained from the [account settings](https://simpleanalytics.com/account) page. It should start with `sa_api_key...`

- `websites` (array, required):

  A list of websites to sync data for. Each website should have the following fields:

    - `hostname` (string, required):
    
      The hostname of the website to sync data for. This should be the same as the hostname in Simple Analytics.
  
    - `metadata_fields` (array[string], optional):

      A list of metadata fields to sync, e.g. `["path_text", "created_at_time"]`. If not specified, no metadata fields will be synced.

- `start_date` (string, optional):

  The date to start syncing data from. If not specified, the plugin will sync data from the beginning of time (or use a start time defined by `period`, if set).

- `end_date` (string, optional): 

  The date to stop syncing data at. If not specified, the plugin will sync data until the current date.

- `period` (string, optional):
  
  The duration of the time window to fetch historical data for, in days, months or years. It is used to calculate `start_date` if it is not specified. If `start_date` is specified, duration is ignored. Examples:
    - `7d`: last 7 days
    - `3m`: last 3 months
    - `1y`: last year


## Example Queries

### List the top 10 pages by views for a given period, excluding robots

```sql
select 
  path, 
  count(*) 
from 
  simple_analytics_page_views 
where 
  hostname = 'mywebsite.com'
  and is_robot is false 
  and added_iso between '2023-01-01' 
  and '2023-02-01'
group by 
  path 
order by
  count desc 
limit 
  10
```

```text
+----------------------------------+---------+
| path                             | count   |
|----------------------------------+---------|
| /                                | 100333  |
| /page                            | 91234   |
| /another-page                    | 84567   |
| /blog/introduction               | 74342   |
| /index                           | 69333   |
| /another/page                    | 64935   |
| /deeply/nested/page              | 50404   |
| /yet/another                     | 42309   |
| /some/page                       | 34433   |
| /about-us                        | 20334   |
+----------------------------------+---------+
```


### List events

```sql
select 
  added_iso, 
  datapoint, 
  path, 
  browser_name 
from 
  simple_analytics_events 
order by 
  added_iso desc 
limit 
  5
```

```text
+-------------------------+-----------+-----------------------------------------------+---------------+
| added_iso               | datapoint | path                                          | browser_name  |
|-------------------------+-----------+-----------------------------------------------+---------------|
| 2023-01-23 19:32:25.68  | 404       | /security                                     | Google Chrome |
| 2023-01-22 20:23:23.379 | 404       | /blog/running-cloudquery-in-gcp               | Google Chrome |
| 2023-01-19 12:04:57.095 | 404       | /docs/plugins/sources/vercel/configuration.md | Brave         |
| 2023-01-19 12:04:36.567 | 404       | /docsss                                       | Firefox       |
| 2023-01-19 01:50:19.259 | 404       | /imgs/gcp-cross-project-service-account       | Google Chrome |
+-------------------------+-----------+-----------------------------------------------+---------------+
```
