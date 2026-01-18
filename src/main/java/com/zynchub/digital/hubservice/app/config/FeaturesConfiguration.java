package com.zynchub.digital.hubservice.app.config;

import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.service.impl.AppConfigFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.service.impl.PropertyFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.service.impl.SessionCachedFeaturesRetriever;
import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.context.annotation.Primary;
import org.springframework.core.env.Environment;

@Configuration
@RequiredArgsConstructor
public class FeaturesConfiguration {

  private final Environment environment;
  private final FeatureMapper featureMapper;

  @Bean(name = "baseFeaturesRetriever")
  @ConditionalOnProperty(
      value = "aws.app-config.features.enabled",
      havingValue = "true",
      matchIfMissing = true)
  public FeaturesRetriever baseAppConfigFeaturesRetriever() {
    return new AppConfigFeaturesRetrieverImpl(environment, featureMapper);
  }

  @Bean(name = "baseFeaturesRetriever")
  @ConditionalOnProperty(
      value = "aws.app-config.features.enabled",
      havingValue = "false",
      matchIfMissing = false)
  public FeaturesRetriever basePropertyFeaturesRetriever() {
    return new PropertyFeaturesRetrieverImpl("features.yml", featureMapper);
  }

  @Bean
  @ConditionalOnProperty(value = "features.session.enabled", havingValue = "true")
  @Primary
  public FeaturesRetriever sessionCachedFeaturesRetriever(
      @Qualifier("baseFeaturesRetriever") FeaturesRetriever delegate,
      ObjectProvider<HttpSession> sessionProvider) {
    return new SessionCachedFeaturesRetriever(delegate, sessionProvider);
  }
}
