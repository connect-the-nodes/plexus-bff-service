package com.zynchub.digital.hubservice.app.domain;


public class StatusInfo {
    private final String status;
    private final String message;

    public StatusInfo(String status, String message) {
        this.status = status;
        this.message = message;
    }

    public String getStatus() { return status; }
    public String getMessage() { return message; }
}
