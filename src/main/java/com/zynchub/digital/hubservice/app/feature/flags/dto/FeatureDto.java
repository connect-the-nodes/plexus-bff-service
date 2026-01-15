package com.zynchub.digital.hubservice.app.feature.flags.dto;

public record FeatureDto(String name, boolean enabled, String description, String parent) {}
