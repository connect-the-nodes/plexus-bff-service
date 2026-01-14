package com.zynchub.digital.hubservice.app.dto;


public class StatusResponseDto {
    private final String status;
    private final String message;
    private final String timestamp;

    public StatusResponseDto(String status, String message, String timestamp) {
        this.status = status;
        this.message = message;
        this.timestamp = timestamp;
    }

    public String getStatus() { return status; }
    public String getMessage() { return message; }
    public String getTimestamp() { return timestamp; }
}
