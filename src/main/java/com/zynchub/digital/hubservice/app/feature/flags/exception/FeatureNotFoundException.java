package com.zynchub.digital.hubservice.app.feature.flags.exception;

public class FeatureNotFoundException extends RuntimeException {
    public FeatureNotFoundException(String message) {
        super(message);
    }
}
