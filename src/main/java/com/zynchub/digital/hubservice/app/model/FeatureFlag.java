package com.zynchub.digital.hubservice.app.model;

import java.io.Serializable;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class FeatureFlag implements Serializable {
    private static final long serialVersionUID = 1L;
    private String name;
    private String parent;
    private boolean enabled;
    private String description;
}
