# Sumo Logic Output Plugin

This plugin writes metrics to a [Sumo Logic HTTP Source][sumologic] using one
of the following data formats:

- `graphite` for Content-Type of `application/vnd.sumologic.graphite`
- `carbon2` for Content-Type of `application/vnd.sumologic.carbon2`
- `prometheus` for Content-Type of `application/vnd.sumologic.prometheus`

⭐ Telegraf v1.16.0
🏷️ logging
💻 all

[sumologic]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/HTTP-Source/Upload-Metrics-to-an-HTTP-Source

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# A plugin that can send metrics to Sumo Logic HTTP metric collector.
[[outputs.sumologic]]
  ## Unique URL generated for your HTTP Metrics Source.
  ## This is the address to send metrics to.
  # url = "https://events.sumologic.net/receiver/v1/http/<UniqueHTTPCollectorCode>"

  ## Data format to be used for sending metrics.
  ## This will set the "Content-Type" header accordingly.
  ## Currently supported formats:
  ## * graphite - for Content-Type of application/vnd.sumologic.graphite
  ## * carbon2 - for Content-Type of application/vnd.sumologic.carbon2
  ## * prometheus - for Content-Type of application/vnd.sumologic.prometheus
  ##
  ## More information can be found at:
  ## https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/HTTP-Source/Upload-Metrics-to-an-HTTP-Source#content-type-headers-for-metrics
  ##
  ## NOTE:
  ## When unset, telegraf will by default use the influx serializer which is currently unsupported
  ## in HTTP Source.
  data_format = "carbon2"

  ## Timeout used for HTTP request
  # timeout = "5s"

  ## Max HTTP request body size in bytes before compression (if applied).
  ## By default 1MB is recommended.
  ## NOTE:
  ## Bear in mind that in some serializer a metric even though serialized to multiple
  ## lines cannot be split any further so setting this very low might not work
  ## as expected.
  # max_request_body_size = 1000000

  ## Additional, Sumo specific options.
  ## Full list can be found here:
  ## https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/HTTP-Source/Upload-Metrics-to-an-HTTP-Source#supported-http-headers

  ## Desired source name.
  ## Useful if you want to override the source name configured for the source.
  # source_name = ""

  ## Desired host name.
  ## Useful if you want to override the source host configured for the source.
  # source_host = ""

  ## Desired source category.
  ## Useful if you want to override the source category configured for the source.
  # source_category = ""

  ## Comma-separated key=value list of dimensions to apply to every metric.
  ## Custom dimensions will allow you to query your metrics at a more granular level.
  # dimensions = ""
```
