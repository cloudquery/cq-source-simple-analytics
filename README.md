# CloudQuery Simple Analytics Source Plugin

This plugin is a Simple Analytics source plugin that can be used to sync data from Simple Analytics to any database, data warehouse, data lake supported by [CloudQuery](https://www.cloudquery.io/), such as PostgreSQL, BigQuery, Athena, and many more.

## Configuration

The following configuration file will sync all data points for `mywebsite.com` to a PostgreSQL database. See [the CloudQuery Quickstart](https://www.cloudquery.io/docs/quickstart) for more information on how to configure the source and destination.

```yaml
kind: source
spec:
  name: "simple-analytics"
  path: "simpleanalytics/simple-analytics"
  version: "${VERSION}"
  backend: "local"
  tables: 
    ["*"]
  destinations: 
    - "postgresql"
  spec:
    user_id: "${SA_USER_ID}"
    api_key: "${SA_API_KEY}"
    websites:
      - hostname: mywebsite.com
        metadata_fields: 
          - fieldname_text
          - fieldname_int
          # - ... 
```

## Example Queries

### List the top 10 pages by views for a given period, excluding robots

```sql
select 
  path, 
  count(*) 
from 
  simple_analytics_data_points 
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
  10;
```

```text
+----------------------------------+---------+
| path                             | count   |
|----------------------------------+---------|
| /                                | 100333  |
| /intro                           | 91234   |
| /how-we-use-cloudquery-for-elt   | 84567   |
| /blog/introduction               | 74342   |
| /google                          | 69333   |
| /another/page                    | 64935   |
| /deeply/nested/page              | 50404   |
| /yet/another                     | 42309   |
| /some/page                       | 34433   |
| /about-us                        | 20334   |
+----------------------------------+---------+
```
