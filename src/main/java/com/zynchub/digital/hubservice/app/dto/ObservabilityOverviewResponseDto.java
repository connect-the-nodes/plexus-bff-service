package com.zynchub.digital.hubservice.app.dto;

import java.time.Instant;

public class ObservabilityOverviewResponseDto {
    private Instant from;
    private Instant to;
    private long totalRequests;
    private long successCount;
    private long failureCount;
    private double successRate;
    private long throttledCount;
    private long timeoutCount;
    private int activeApiCount;
    private double peakRps;
    private int activeFlowCount;
    private int idleFlowCount;
    private double goldenSuccessRate;

    public ObservabilityOverviewResponseDto() {
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

    public long getSuccessCount() {
        return successCount;
    }

    public void setSuccessCount(long successCount) {
        this.successCount = successCount;
    }

    public long getFailureCount() {
        return failureCount;
    }

    public void setFailureCount(long failureCount) {
        this.failureCount = failureCount;
    }

    public double getSuccessRate() {
        return successRate;
    }

    public void setSuccessRate(double successRate) {
        this.successRate = successRate;
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

    public int getActiveApiCount() {
        return activeApiCount;
    }

    public void setActiveApiCount(int activeApiCount) {
        this.activeApiCount = activeApiCount;
    }

    public double getPeakRps() {
        return peakRps;
    }

    public void setPeakRps(double peakRps) {
        this.peakRps = peakRps;
    }

    public int getActiveFlowCount() {
        return activeFlowCount;
    }

    public void setActiveFlowCount(int activeFlowCount) {
        this.activeFlowCount = activeFlowCount;
    }

    public int getIdleFlowCount() {
        return idleFlowCount;
    }

    public void setIdleFlowCount(int idleFlowCount) {
        this.idleFlowCount = idleFlowCount;
    }

    public double getGoldenSuccessRate() {
        return goldenSuccessRate;
    }

    public void setGoldenSuccessRate(double goldenSuccessRate) {
        this.goldenSuccessRate = goldenSuccessRate;
    }
}
