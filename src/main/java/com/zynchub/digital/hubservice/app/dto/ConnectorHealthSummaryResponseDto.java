package com.zynchub.digital.hubservice.app.dto;

public class ConnectorHealthSummaryResponseDto {
    private String connectorId;
    private long totalRequests;
    private double successRate;
    private double failureRate;
    private double p95LatencyMs;
    private int activeTenants;
    private String status;

    public ConnectorHealthSummaryResponseDto() {
    }

    public String getConnectorId() {
        return connectorId;
    }

    public void setConnectorId(String connectorId) {
        this.connectorId = connectorId;
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

    public int getActiveTenants() {
        return activeTenants;
    }

    public void setActiveTenants(int activeTenants) {
        this.activeTenants = activeTenants;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }
}
