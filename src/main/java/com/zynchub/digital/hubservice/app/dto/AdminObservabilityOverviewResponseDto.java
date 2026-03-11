package com.zynchub.digital.hubservice.app.dto;

import java.time.Instant;

public class AdminObservabilityOverviewResponseDto {
    private Instant from;
    private Instant to;
    private long totalRequests;
    private double successRate;
    private double failureRate;
    private double p95LatencyMs;
    private long throttledCount;
    private long timeoutCount;
    private int activeTenants;
    private int activeConnectors;
    private String topFailingConnectorId;
    private String topFailingTenantId;

    public AdminObservabilityOverviewResponseDto() {
    }

    public Instant getFrom() {
        return from;
    }

    public void setFrom(Instant from) {
        this.from = from;
    }

    public Instant getTo() {
        return to;
    }

    public void setTo(Instant to) {
        this.to = to;
    }

    public long getTotalRequests() {
        return totalRequests;
    }

    public void setTotalRequests(long totalRequests) {
        this.totalRequests = totalRequests;
    }

    public double getSuccessRate() {
        return successRate;
    }

    public void setSuccessRate(double successRate) {
        this.successRate = successRate;
    }

    public double getFailureRate() {
        return failureRate;
    }

    public void setFailureRate(double failureRate) {
        this.failureRate = failureRate;
    }

    public double getP95LatencyMs() {
        return p95LatencyMs;
    }

    public void setP95LatencyMs(double p95LatencyMs) {
        this.p95LatencyMs = p95LatencyMs;
    }

    public long getThrottledCount() {
        return throttledCount;
    }

    public void setThrottledCount(long throttledCount) {
        this.throttledCount = throttledCount;
    }

    public long getTimeoutCount() {
        return timeoutCount;
    }

    public void setTimeoutCount(long timeoutCount) {
        this.timeoutCount = timeoutCount;
    }

    public int getActiveTenants() {
        return activeTenants;
    }

    public void setActiveTenants(int activeTenants) {
        this.activeTenants = activeTenants;
    }

    public int getActiveConnectors() {
        return activeConnectors;
    }

    public void setActiveConnectors(int activeConnectors) {
        this.activeConnectors = activeConnectors;
    }

    public String getTopFailingConnectorId() {
        return topFailingConnectorId;
    }

    public void setTopFailingConnectorId(String topFailingConnectorId) {
        this.topFailingConnectorId = topFailingConnectorId;
    }

    public String getTopFailingTenantId() {
        return topFailingTenantId;
    }

    public void setTopFailingTenantId(String topFailingTenantId) {
        this.topFailingTenantId = topFailingTenantId;
    }
}
