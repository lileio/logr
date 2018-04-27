# Logr

Logr is a simple helper for [Logrus](https://github.com/sirupsen/logrus) which helps wrap loggers to log [Opentracing](http://opentracing.io/) information (atm, just [Zipkin](https://zipkin.io/)).

### Example

If you have a context that might already be in a trace, then you can simple create a new Logrus logger with your context. This a gRPC handler for example.

``` go
func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
        logr.WithCtx(ctx).Info("Called GetFeature woo!")
}
```