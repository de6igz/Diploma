// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"aletheia-public-api/interfaces"
	v1 "aletheia-public-api/interfaces/types/v1"
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type metricsRules struct {
	next            interfaces.Rules
	requestCount    metrics.Counter
	requestCountAll metrics.Counter
	requestLatency  metrics.Histogram
}

func metricsMiddlewareRules(next interfaces.Rules) interfaces.Rules {
	return &metricsRules{
		next:            next,
		requestCount:    RequestCount.With("service", "Rules"),
		requestCountAll: RequestCountAll.With("service", "Rules"),
		requestLatency:  RequestLatency.With("service", "Rules"),
	}
}

func (m metricsRules) GetRules(ctx context.Context, userId int64) (items v1.RulesResponse, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "getRules", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "getRules", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "getRules").Add(1)

	return m.next.GetRules(ctx, userId)
}

func (m metricsRules) GetRuleByID(ctx context.Context, userId int64, ruleId string, ruleType string) (rule *v1.RuleDetailResponse, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "getRuleByID", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "getRuleByID", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "getRuleByID").Add(1)

	return m.next.GetRuleByID(ctx, userId, ruleId, ruleType)
}

func (m metricsRules) GetAvailableRules(ctx context.Context, userId int64) (items v1.RulesResponse, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "getAvailableRules", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "getAvailableRules", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "getAvailableRules").Add(1)

	return m.next.GetAvailableRules(ctx, userId)
}

func (m metricsRules) DeleteRuleByID(ctx context.Context, userId int64, req v1.DeleteRuleRequest) (status bool, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "deleteRuleByID", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "deleteRuleByID", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "deleteRuleByID").Add(1)

	return m.next.DeleteRuleByID(ctx, userId, req)
}

func (m metricsRules) CreateRule(ctx context.Context, userId int64, request v1.CreateRuleRequest) (status bool, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "createRule", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "createRule", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "createRule").Add(1)

	return m.next.CreateRule(ctx, userId, request)
}

func (m metricsRules) UpdateRuleById(ctx context.Context, userId int64, request v1.UpdateRuleRequest) (status bool, err error) {

	defer func(_begin time.Time) {
		m.requestCount.With("method", "updateRuleById", "success", fmt.Sprint(err == nil)).Add(1)
		m.requestLatency.With("method", "updateRuleById", "success", fmt.Sprint(err == nil)).Observe(time.Since(_begin).Seconds())
	}(time.Now())

	m.requestCountAll.With("method", "updateRuleById").Add(1)

	return m.next.UpdateRuleById(ctx, userId, request)
}
