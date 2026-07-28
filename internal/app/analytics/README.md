
```
Event -> EventSink -> *eventBackgroundService -> EventRepository
-> *eventProcessorWorker (notify via EventProcessorWorkerHandle) -> EventProcessor
-> Metric -> MetricSink -> *metricBackgroundService -> MetricRepository

```