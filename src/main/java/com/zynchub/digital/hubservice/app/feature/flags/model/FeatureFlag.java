package com.zynchub.digital.hubservice.app.feature.flags.model;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class FeatureFlag {
    private String name;
    private String parent;
    private boolean enabled;
    private String description;
}
