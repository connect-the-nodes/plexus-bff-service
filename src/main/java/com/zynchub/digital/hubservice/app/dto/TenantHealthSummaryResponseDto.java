package com.zynchub.digital.hubservice.app.dto;

public class TenantHealthSummaryResponseDto {
    private String tenantId;
    private long totalRequests;
    private double successRate;
    private double failureRate;
    private double quotaUsedPercent;
    private long throttledCount;
    private long timeoutCount;
    private String status;

    public TenantHealthSummaryResponseDto() {
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
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

    public double getQuotaUsedPercent() {
        return quotaUsedPercent;
    }

    public void setQuotaUsedPercent(double quotaUsedPercent) {
        this.quotaUsedPercent = quotaUsedPercent;
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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }
}
