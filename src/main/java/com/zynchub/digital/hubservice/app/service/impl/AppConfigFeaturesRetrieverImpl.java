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
import java.nio.charset.StandardCharsets;
import java.util.List;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.env.Environment;
import software.amazon.awssdk.services.appconfigdata.AppConfigDataClient;
import software.amazon.awssdk.services.appconfigdata.model.GetLatestConfigurationRequest;
import software.amazon.awssdk.services.appconfigdata.model.StartConfigurationSessionRequest;

@Slf4j
@RequiredArgsConstructor
public class AppConfigFeaturesRetrieverImpl implements FeaturesRetriever {

  private static final String APP_ID = "aws.app-config.features.application-id";
  private static final String ENV_ID = "aws.app-config.features.environment-id";
  private static final String CONFIG_ID = "aws.app-config.features.configuration-id";

  private final Environment environment;
  private final FeatureMapper featureMapper;

  @Override
  public List<FeatureFlag> retrieveFeatures() {
    String content = fetchFeaturesContent();
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
      log.error("Unable to parse features yaml retrieved from AWS AppConfig", e);
      throw new FeatureNotFoundException(
          "Unable to parse features yaml retrieved from AWS AppConfig", e);
    }
  }

  private String fetchFeaturesContent() {
    try (AppConfigDataClient client = AppConfigDataClient.create()) {
      String token = getAwsConfigurationToken(client);
      var request = GetLatestConfigurationRequest.builder().configurationToken(token).build();
      var latest = client.getLatestConfiguration(request);
      return new String(latest.configuration().asByteArray(), StandardCharsets.UTF_8);
    }
  }

  private String getAwsConfigurationToken(final AppConfigDataClient client) {
    String applicationId = environment.getProperty(APP_ID);
    String environmentId = environment.getProperty(ENV_ID);
    String configurationId = environment.getProperty(CONFIG_ID);
    if (isBlank(applicationId) || isBlank(environmentId) || isBlank(configurationId)) {
      throw new FeatureNotFoundException(
          "AWS AppConfig identifiers are not configured for feature flags");
    }
    var request =
        StartConfigurationSessionRequest.builder()
            .applicationIdentifier(applicationId)
            .environmentIdentifier(environmentId)
            .configurationProfileIdentifier(configurationId)
            .build();
    var session = client.startConfigurationSession(request);
    return session.initialConfigurationToken();
  }

  private static boolean isBlank(String value) {
    return value == null || value.isBlank();
  }
}
