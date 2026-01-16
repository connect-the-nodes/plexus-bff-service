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
import java.util.List;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.client.RestClient;

@Slf4j
@RequiredArgsConstructor
public class RemoteFeaturesRetrieverImpl implements FeaturesRetriever {

    private final RestClient restClient;
    private final String featuresUrl;
    private final FeatureMapper featureMapper;

    @Override
    public List<FeatureFlag> retrieveFeatures() {
        String content = restClient.get().uri(featuresUrl).retrieve().body(String.class);
        if (content == null || content.isBlank()) {
            return List.of();
        }
        try {
            var featureList =
                    new ObjectMapper(new YAMLFactory())
                            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
                            .readValue(content, FeatureDtoList.class);
            return featureList == null || featureList.getFeatures() == null
                    ? List.of()
                    : featureMapper.mapFromDto(featureList.getFeatures());
        } catch (IOException e) {
            log.error("Unable to parse features yaml retrieved from {}", featuresUrl, e);
            throw new FeatureNotFoundException("Unable to parse features yaml from remote source");
        }
    }
}
