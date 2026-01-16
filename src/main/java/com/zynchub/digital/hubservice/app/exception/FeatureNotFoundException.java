package com.zynchub.digital.hubservice.app.exception;

public class FeatureNotFoundException extends RuntimeException {
    public FeatureNotFoundException(String message) {
        super(message);
    }

    public FeatureNotFoundException(String message, Throwable cause) {
        super(message, cause);
    }
}
