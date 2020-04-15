package filter

import (
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
)

type Handler func(string) Option

var optionMap = map[string]Handler{
	"domainProject": WithDomainProject,
	"app":           WithApp,
	"serviceName":   WithServiceName,
	"serviceId":     WithServiceId,
}

type Option func(c *condition)

func WithDomainProject(domainProject string) Option {
	return func(f *condition) { f.DomainProject = domainProject }
}

func WithApp(app string) Option {
	return func(f *condition) { f.App = app }
}

func WithServiceName(serviceName string) Option {
	return func(f *condition) { f.ServiceName = serviceName }
}

func WithServiceId(serviceId string) Option {
	return func(f *condition) { f.ServiceId = serviceId }
}

type condition struct {
	DomainProject string
	App           string
	ServiceName   string
	ServiceId     string
}

func (c *condition) Match(svc *pb.SyncService) bool {
	if c.DomainProject != "" && c.DomainProject != svc.DomainProject {
		return false
	}

	if c.App != "" && c.App != svc.App {
		return false
	}

	if c.ServiceName != "" && c.ServiceName != svc.Name {
		return false
	}

	if c.ServiceId != "" && c.ServiceId != svc.ServiceId {
		return false
	}
	return true
}

type Matcher interface {
	Match(svc *pb.SyncService) bool
}

func MapToMatcher(m map[string]string) Matcher {
	ops := make([]Option, 0, len(m))
	for key, val := range m {
		fn, ok := optionMap[key]
		if !ok {
			panic("Unsupported the condition key： " + key)
		}
		ops = append(ops, fn(val))
	}
	return NewMatcher(ops...)
}

func NewMatcher(ops ...Option) Matcher {
	f := &condition{}
	for _, op := range ops {
		op(f)
	}
	return f
}

type matchers []Matcher

func (m matchers) Match(svc *pb.SyncService) bool {
	for _, val := range m {
		if val.Match(svc) {
			return true
		}
	}
	return false
}

type Filter interface {
	Filter(data *pb.SyncData) *pb.SyncData
}

type filter struct {
	matchers matchers
	approve  bool
}

func (f *filter) Filter(data *pb.SyncData) *pb.SyncData {
	nd := &pb.SyncData{
		Services:  make([]*pb.SyncService, 0, len(data.Services)),
		Instances: make([]*pb.SyncInstance, 0, len(data.Instances)),
	}

	for _, svc := range data.Services {
		if f.approve == f.matchers.Match(svc) {
			nd.Services = append(nd.Services, svc)
		}
	}
	return nd
}

func NewWhiteList(filters ...Matcher) Filter {
	return &filter{matchers: filters, approve: true}
}

func NewBlackList(filters ...Matcher) Filter {
	return &filter{matchers: filters, approve: false}
}
