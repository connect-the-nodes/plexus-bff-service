package com.zynchub.digital.hubservice.app.feature.flags.config;

import com.zynchub.digital.hubservice.app.feature.flags.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.feature.flags.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.feature.flags.service.impl.PropertyFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.feature.flags.service.impl.RemoteFeaturesRetrieverImpl;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestClient;

@Configuration
@RequiredArgsConstructor
public class FeaturesConfiguration {

    private final FeatureMapper featureMapper;

    @Bean
    @ConditionalOnProperty(
            value = "features.remote.enabled",
            havingValue = "true",
            matchIfMissing = false)
    public FeaturesRetriever remoteFeaturesRetriever(
            RestClient.Builder restClientBuilder,
            @Value("${features.remote.url}") String featuresUrl) {
        return new RemoteFeaturesRetrieverImpl(restClientBuilder.build(), featuresUrl, featureMapper);
    }

    @Bean
    @ConditionalOnProperty(
            value = "features.remote.enabled",
            havingValue = "false",
            matchIfMissing = true)
    public FeaturesRetriever propertyFeaturesRetriever(
            @Value("${features.file:features.yml}") String featureFileName) {
        return new PropertyFeaturesRetrieverImpl(featureFileName, featureMapper);
    }
}
