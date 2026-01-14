package com.zynchub.digital.hubservice.app.dto;


public class ErrorResponseDto {
    private final String error;
    private final String message;

    public ErrorResponseDto(String error, String message) {
        this.error = error;
        this.message = message;
    }

    public String getError() { return error; }
    public String getMessage() { return message; }
}
