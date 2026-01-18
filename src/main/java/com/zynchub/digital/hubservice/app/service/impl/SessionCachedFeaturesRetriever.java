package com.zynchub.digital.hubservice.app.service.impl;

import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import jakarta.servlet.http.HttpSession;
import java.util.List;
import org.springframework.beans.factory.ObjectProvider;

public class SessionCachedFeaturesRetriever implements FeaturesRetriever {

  private static final String SESSION_KEY = "feature_flags";

  private final FeaturesRetriever delegate;
  private final ObjectProvider<HttpSession> sessionProvider;

  public SessionCachedFeaturesRetriever(
      FeaturesRetriever delegate, ObjectProvider<HttpSession> sessionProvider) {
    this.delegate = delegate;
    this.sessionProvider = sessionProvider;
  }

  @Override
  public List<FeatureFlag> retrieveFeatures() {
    HttpSession session = sessionProvider.getIfAvailable();
    if (session != null) {
      Object cached = session.getAttribute(SESSION_KEY);
      if (cached instanceof List<?> cachedList && cacheIsValid(cachedList)) {
        @SuppressWarnings("unchecked")
        List<FeatureFlag> features = (List<FeatureFlag>) cachedList;
        return features;
      }
    }

    List<FeatureFlag> features = delegate.retrieveFeatures();
    if (session != null) {
      session.setAttribute(SESSION_KEY, features);
    }
    return features;
  }

  private boolean cacheIsValid(List<?> cachedList) {
    return cachedList.stream().allMatch(FeatureFlag.class::isInstance);
  }
}
