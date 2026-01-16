package com.zynchub.digital.hubservice.app.service;

import com.zynchub.digital.hubservice.app.exception.FeatureNotFoundException;
import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import java.util.List;

public interface FeaturesRetriever {

    List<FeatureFlag> retrieveFeatures();

    default boolean isActive(String featureName) {
        List<FeatureFlag> features = retrieveFeatures();
        return verifyFeatureInList(featureName, features);
    }

    private static boolean verifyFeatureInList(String featureName, List<FeatureFlag> features) {
        FeatureFlag feature =
                features.stream()
                        .filter(f -> f.getName().equalsIgnoreCase(featureName))
                        .findFirst()
                        .orElseThrow(
                                () -> new FeatureNotFoundException("Feature " + featureName + " is not valid"));

        boolean isFeatureActive = feature.isEnabled();

        while (isFeatureActive && !feature.getName().equalsIgnoreCase(feature.getParent())) {
            final FeatureFlag temp = feature;
            feature =
                    features.stream()
                            .filter(f -> f.getName().equalsIgnoreCase(temp.getParent()))
                            .findFirst()
                            .orElseThrow(
                                    () ->
                                            new FeatureNotFoundException(
                                                    "Please provide valid parent feature to the given feature "
                                                            + temp.getName()));
            isFeatureActive = feature.isEnabled();
        }

        return isFeatureActive;
    }
}
