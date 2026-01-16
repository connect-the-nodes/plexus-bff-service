package com.zynchub.digital.hubservice.app.service.impl;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.zynchub.digital.hubservice.app.dto.FeatureDtoList;
import com.zynchub.digital.hubservice.app.exception.FeatureNotFoundException;
import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import java.io.IOException;
import java.net.URL;
import java.util.List;
import lombok.extern.slf4j.Slf4j;
import org.springframework.util.ResourceUtils;

@Slf4j
public class PropertyFeaturesRetrieverImpl implements FeaturesRetriever {

    private final List<FeatureFlag> features;

    public PropertyFeaturesRetrieverImpl(String fileName, FeatureMapper featureMapper) {
        var mapper =
                new ObjectMapper(new YAMLFactory())
                        .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        try {
            URL file = ResourceUtils.getURL("classpath:" + fileName);
            var featureDtoList = mapper.readValue(file, FeatureDtoList.class);
            features =
                    featureDtoList == null || featureDtoList.getFeatures() == null
                            ? List.of()
                            : featureMapper.mapFromDto(featureDtoList.getFeatures());
        } catch (IOException e) {
            log.error("Error occurred while reading features from {} file", fileName);
            throw new FeatureNotFoundException("Error while reading features.yml");
        }
    }

    @Override
    public List<FeatureFlag> retrieveFeatures() {
        return features;
    }
}
