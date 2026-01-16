package com.zynchub.digital.hubservice.app.feature;

import com.zynchub.digital.hubservice.app.feature.FeatureAssociation;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.stereotype.Component;

@Slf4j
@Aspect
@Component
@RequiredArgsConstructor
public class FeaturesAspect {

    private final FeaturesRetriever featuresRetriever;

    @Around("@within(featureAssociation) || @annotation(featureAssociation)")
    public Object checkAspect(ProceedingJoinPoint joinPoint, FeatureAssociation featureAssociation)
            throws Throwable {
        boolean active = featuresRetriever.isActive(featureAssociation.name());
        if (featureAssociation.invert() ? !active : active) {
            return joinPoint.proceed();
        }
        log.info("Feature {} is not enabled", featureAssociation.name());
        throw new AccessDeniedException("Feature " + featureAssociation.name() + " is not enabled!");
    }
}
